package main

import (
	"bytes"
	"context"
	_ "embed"
	"flag"
	"fmt"
	"go/format"
	"html/template"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func main() {
	// Flags
	packageDirFlag := flag.String("package-dir", "", "Target package directory where generated files will be written (required)")
	schemaPathFlag := flag.String("schema", "", "Path to SQL schema file (required)")
	// Optional flags
	packageNameFlag := flag.String("package-name", "", "Override the package name in generated files (default: use directory name)")
	tableDefFileFlag := flag.String("table-def-file", "table-definitions.gen.go", "Name of the generated table definitions file (default: table-definitions.gen.go)")
	tableGraphFileFlag := flag.String("table-graph-file", "tables-graph.gen.go", "Name of the generated tables graph file (default: tables-graph.gen.go). If empty, the graph file will not be generated.")
	tomasqlImportModeFlag := flag.String("tomasql-import-mode", "full", "How to import tomasql package: 'full' (tomasql.Type), 'dot' (. import), 'none' (no import)")
	postgresImageFlag := flag.String("postgres-image", "postgres:latest", "Postgres image to use for tables generation (default: postgres:latest)")

	flag.Parse()

	if *packageDirFlag == "" || *schemaPathFlag == "" {
		flag.Usage()
		log.Fatal("--package-dir and --schema are required")
	}

	packageDir := *packageDirFlag
	schemaPath := *schemaPathFlag
	packageName := *packageNameFlag
	tableDefFile := *tableDefFileFlag
	tableGraphFile := *tableGraphFileFlag
	tomasqlImportMode := *tomasqlImportModeFlag
	dockerImage := *postgresImageFlag

	// Use provided package name or default to directory name
	pkgName := packageName
	if pkgName == "" {
		pkgName = filepath.Base(packageDir)
	}

	// Get the directory of the current Go source file
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	tmpl, err := template.New("table-def.tmpl").Funcs(template.FuncMap{
		"TomasqlImportMode": func() string { return tomasqlImportMode },
		"TomasqlPrefix": func() string {
			switch tomasqlImportMode {
			case "full":
				return "tomasql."
			case "dot", "none":
				return ""
			default:
				return "tomasql."
			}
		},
	}).ParseFiles(filepath.Join(dir, "table-def.tmpl"))
	if err != nil {
		panic(err)
	}

	// Helper function to execute template and format the output
	executeAndFormat := func(tmpl *template.Template, data interface{}, outputPath string) error {
		// Execute template to a buffer first
		var buf bytes.Buffer
		err := tmpl.Execute(&buf, data)
		if err != nil {
			return err
		}

		// Format the generated code
		formatted, err := format.Source(buf.Bytes())
		if err != nil {
			// If formatting fails, use the unformatted version
			log.Printf("Warning: failed to format %s: %v", outputPath, err)
			formatted = buf.Bytes()
		}

		// Write to file
		return os.WriteFile(outputPath, formatted, 0644)
	}

	container, err := SetupTestContainer(nil, schemaPath, dockerImage)
	if err != nil {
		panic(err)
	}

	tableDefData, err := getTableDefinitionFromTestDB(container, pkgName)
	if err != nil {
		panic(err)
	}

	outPath := filepath.Join(packageDir, tableDefFile)
	err = executeAndFormat(tmpl, tableDefData, outPath)
	if err != nil {
		panic(err)
	}
	log.Println("Generated and formatted table definitions")

	// Generate graph definitions if tableGraphFile is not empty
	if tableGraphFile != "" {
		tmplGraph, err := template.New("tables-graph.tmpl").Funcs(template.FuncMap{
			"TomasqlImportMode": func() string { return tomasqlImportMode },
			"TomasqlPrefix": func() string {
				switch tomasqlImportMode {
				case "full":
					return "tomasql."
				case "dot", "none":
					return ""
				default:
					return "tomasql."
				}
			},
		}).ParseFiles(filepath.Join(dir, "tables-graph.tmpl"))
		if err != nil {
			panic(err)
		}

		dbGraphData, err := getDbGraphFromTestDB(container, pkgName)
		if err != nil {
			panic(err)
		}

		outGraphPath := filepath.Clean(filepath.Join(packageDir, tableGraphFile))
		err = executeAndFormat(tmplGraph, dbGraphData, outGraphPath)
		if err != nil {
			panic(err)
		}
		log.Println("Generated and formatted table graph")
	} else {
		log.Println("Skipping table graph generation as --table-graph-file is empty")
	}
}

func getTableDefinitionFromTestDB(container *sqlx.DB, pkgName string) (*TemplateData, error) {
	type row struct {
		TableName     string `db:"table_name"`
		ColumnName    string `db:"column_name"`
		UdtName       string `db:"udt_name"` // type
		IsNullable    bool   `db:"is_nullable"`
		IsUserDefined bool   `db:"is_user_defined"`
		IsEnum        bool   `db:"is_enum"`
		BaseType      string `db:"base_type"`
	}
	result := []row{}
	err := container.Select(&result, `
SELECT c.table_name,
                c.column_name,
                c.udt_name,
                CASE WHEN c.is_nullable = 'NO' THEN false ELSE true END as is_nullable,
                CASE
                    WHEN c.data_type = 'USER-DEFINED' THEN true
                    ELSE false
                    END as is_user_defined,
                CASE
                    WHEN t.typcategory = 'E' THEN true
                    ELSE false
                    END AS is_enum,
                CASE
                    WHEN t.typcategory = 'E' THEN 'string' -- Enums are stored as strings
                    WHEN t.typcategory = 'C' THEN 'composite' -- Composite types map to Go structs
                    WHEN t.typcategory = 'D' AND t.typbasetype <> 0 THEN bt.typname
                    ELSE t.typname
                    END AS base_type
FROM information_schema.columns c
         JOIN pg_type t ON c.udt_name = t.typname
         JOIN pg_namespace n ON t.typnamespace = n.oid
         LEFT JOIN pg_type bt ON t.typbasetype = bt.oid -- To get the base type of a domain
WHERE c.table_schema = 'public'
ORDER BY 1, 2
`)
	if err != nil {
		return nil, err
	}

	data := &TemplateData{Package: pkgName}
	importsSet := map[string]struct{}{}
	var currTable *Table
	for _, item := range result {
		if currTable == nil || currTable.SqlName != item.TableName {
			currTable = &Table{
				SqlName:     item.TableName,
				TypeDefName: snakeToCamel(item.TableName),
				Columns:     []Column{},
			}
			data.Tables = append(data.Tables, currTable)
		}
		// TODO: constants for enums?
		mappedType := psqlTypeToGo(item.BaseType)
		if idx := strings.Index(mappedType, "."); idx > 0 { // pkg.Type pattern
			pkgPart := mappedType[:idx]
			importsSet[pkgPart] = struct{}{}
		}
		column := Column{
			Name:    snakeToCamel(item.ColumnName),
			SqlName: item.ColumnName,
			Type:    mappedType,
		}

		currTable.Columns = append(currTable.Columns, column)
	}

	if len(importsSet) > 0 {
		for k := range importsSet {
			data.Imports = append(data.Imports, k)
		}
		slices.Sort(data.Imports)
	}

	return data, nil
}

func psqlTypeToGo(psqlType string) string {
	switch psqlType {
	case "bool":
		return "bool"
	case "float4":
		return "float32"
	case "float8":
		return "float64"
	case "int2":
		return "int16"
	case "int4":
		return "int"
	case "int8":
		return "int64"
	case "uuid":
		return "string"
	case "bpchar":
		return "string"
	case "string":
		return "string"
	case "text":
		return "string"
	case "varchar":
		return "string"
	case "numeric":
		// numeric/decimal: defaulting to float64; consider using a fixed-point/decimal type if precision is critical
		return "float64"
	case "timestamp":
		return "time.Time"
	case "timestamptz":
		return "time.Time"
	case "date":
		return "time.Time"

	default:
		panic("Unknown type: " + psqlType)
	}
}

func snakeToCamel(name string) string {
	parts := strings.Split(name, "_")
	for i := range parts {
		parts[i] = cases.Title(language.English).String(parts[i])
	}
	return strings.Join(parts, "")
}

type TemplateData struct {
	Package string
	Tables  []*Table
	Imports []string
}

type Table struct {
	TypeDefName string
	SqlName     string
	Columns     []Column
}

type Column struct {
	Name    string
	SqlName string
	Type    string
}

// getDbGraphFromTestDB retrieves the database graph from the test database. Note that at the moment it only retunrs
// 'forward' links, i.e. from the table that has a foreign key to the table that is referenced by the foreign key.
func getDbGraphFromTestDB(container *sqlx.DB, pkgName string) (*DbGraphTemplateData, error) {
	type row struct {
		FromTable  string `db:"from_table"`
		ToTable    string `db:"to_table"`
		FromColumn string `db:"from_column"`
		ToColumn   string `db:"to_column"`
	}
	result := []row{}
	err := container.Select(&result, `SELECT
    tc.table_name AS from_table,
    ccu.table_name AS to_table,
    kcu.column_name AS from_column,
    ccu.column_name AS to_column
FROM
    information_schema.table_constraints AS tc
        JOIN information_schema.key_column_usage AS kcu
             ON tc.constraint_name = kcu.constraint_name
                 AND tc.constraint_schema = kcu.constraint_schema
        JOIN information_schema.constraint_column_usage AS ccu
             ON tc.constraint_name = ccu.constraint_name
                 AND tc.constraint_schema = ccu.constraint_schema
WHERE
    tc.constraint_type = 'FOREIGN KEY'`)
	if err != nil {
		return nil, err
	}

	data := &DbGraphTemplateData{
		Package: pkgName,
		Links:   map[string]map[string]map[string]*Link{},
	}

	addLink := func(fromType, toType, fromField, toField string) {
		if _, exists := data.Links[fromType]; !exists {
			data.Links[fromType] = map[string]map[string]*Link{}
		}
		if _, exists := data.Links[fromType][toType]; !exists {
			data.Links[fromType][toType] = map[string]*Link{}
		}
		data.Links[fromType][toType][fromField] = &Link{
			FromTable:  fromType,
			FromColumn: fromField,
			ToTable:    toType,
			ToColumn:   toField,
		}
	}

	for _, item := range result {
		fromType := snakeToCamel(item.FromTable)
		toType := snakeToCamel(item.ToTable)
		fromField := snakeToCamel(item.FromColumn)
		toField := snakeToCamel(item.ToColumn)
		addLink(fromType, toType, fromField, toField)
		// inverse link
		addLink(toType, fromType, toField, fromField)
	}

	return data, nil
}

type DbGraphTemplateData struct {
	Package string
	// fromTable -> toTable -> fromColumn -> Link
	Links map[string]map[string]map[string]*Link
}

type Link struct {
	FromTable  string
	FromColumn string
	ToTable    string
	ToColumn   string
}

func SetupTestContainer(tt *testing.T, schemaPath string, dockerImage string) (*sqlx.DB, error) {

	// we also use this function outside of tests
	logFn := log.Printf
	if tt != nil {
		logFn = tt.Logf
	}

	ctx := context.Background()
	containerName := "tomasql-table-def-gen-testcontainer"
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        dockerImage,
			ExposedPorts: []string{"5432/tcp"},
			WaitingFor:   wait.ForLog("database system is ready to accept connections"),
			Name:         containerName,
			Env: map[string]string{
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "postgres",
				"POSTGRES_DB":       "postgres_user",
			},
		},
		Started: true,
		Reuse:   true,
	})
	if err != nil {
		return nil, err
	}

	port, err := nat.NewPort("tcp", "5432")
	if err != nil {
		return nil, err
	}
	port, err = container.MappedPort(ctx, port)
	if err != nil {
		return nil, err
	}

	newPostgresUrl := func(dbName string) string {
		return fmt.Sprintf("postgresql://postgres:postgres@%s:%d%s?sslmode=disable&connect_timeout=3", "0.0.0.0", port.Int(), dbName)
	}

	purl := newPostgresUrl("")
	log.Printf("connecting to postgres through: %s", purl)
	var rootDB *sqlx.DB
	for i := 0; i < 10; i++ {
		rootDB, err = sqlx.Connect("postgres", purl)
		if err == nil {
			break
		}
		logFn("failed to connect sql: %d/10", i+1)

		time.Sleep(time.Second * 3)
	}
	if err != nil {
		panic(fmt.Sprintf("failed to connect to postgres: %s", err))
	}

	// create new database

	nw := time.Now()
	tm := nw.Format("2006_01_02__15_04_05")
	randDbName := "testcontainer_" + strconv.Itoa(int(nw.UnixMilli())) + "_" + strings.ToLower(tm) + "_" + createRandString(16)
	schema := readDbSchema(schemaPath)

	_, err = rootDB.Exec("create database " + randDbName)
	if err != nil {
		return nil, err
	}

	purl = newPostgresUrl("/" + randDbName)
	log.Printf("trying to connect to newly created database: %s", purl)
	db, err := sqlx.Connect("postgres", purl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to newly created test database instance: %s", err)
	}

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	logFn("set up test database: %s", randDbName)

	return db, nil
}

func readDbSchema(schemaPath string) string {
	if schemaPath == "" {
		log.Fatal("schema path must not be empty")
	}
	content, err := os.ReadFile(schemaPath)
	if err != nil {
		panic(fmt.Sprintf("failed to read sql schema '%s': %s", schemaPath, err))
	}

	return string(content)
}

func createRandString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
