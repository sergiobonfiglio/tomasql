package main

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"html/template"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

//go:generate go run .

func main() {

	// Get the directory of the current Go source file
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	tmpl, err := template.ParseFiles(filepath.Join(dir, "table-def.tmpl"))
	if err != nil {
		panic(err)
	}

	outPath := filepath.Join(dir, "..", "..", "goql", "table-definitions.gen.go")
	outputFile, err := os.Create(outPath)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	tableDefData, err := getTableDefinitionFromTestDB()
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(outputFile, tableDefData)
	if err != nil {
		panic(err)
	}
	log.Println("Generated table definitions")
}

func getTableDefinitionFromTestDB() (*TemplateData, error) {
	container, err := SetupTestContainer(nil)
	if err != nil {
		return nil, err
	}

	type row struct {
		TableName     string `db:"table_name"`
		ColumnName    string `db:"column_name"`
		UdtName       string `db:"udt_name"` //type
		IsNullable    bool   `db:"is_nullable"`
		IsUserDefined bool   `db:"is_user_defined"`
		IsEnum        bool   `db:"is_enum"`
		BaseType      string `db:"base_type"`
	}
	result := []row{}
	err = container.Select(&result, `
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

	data := &TemplateData{Package: "goql"}
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
		//TODO: constants for enums?
		column := Column{
			Name:    snakeToCamel(item.ColumnName),
			SqlName: item.ColumnName,
			Type:    psqlTypeToGo(item.BaseType),
		}

		currTable.Columns = append(currTable.Columns, column)
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

func SetupTestContainer(tt *testing.T) (*sqlx.DB, error) {

	// we also use this function outside of tests
	logFn := log.Printf
	if tt != nil {
		logFn = tt.Logf
	}

	// in case creating testcontainers is stuck, purge all testcontainers instances running `make purge-testcontainers`
	ctx := context.Background()
	containerName := "postgis-testcontainer"
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "public.ecr.aws/bigprofiles/postgres-postgis:15.1.0",
			ExposedPorts: []string{"5432/tcp"},
			WaitingFor:   wait.ForLog("database system is ready to accept connections"),
			// necessary to reuse container
			Name: containerName,
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
	schema := readDbSchema()

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

	return rootDB, nil
}

func readDbSchema() string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	sqlSchema := filepath.Join(basepath, "example_schema.sql")
	content, err := os.ReadFile(sqlSchema)
	if err != nil {
		panic(fmt.Sprintf("failed to read sql schema: %s", err))
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
