package setup

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

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
