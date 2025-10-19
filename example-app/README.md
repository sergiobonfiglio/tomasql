# TomaSQL Test Application

This is a demonstration application showing how to use the TomaSQL library with custom table definitions.

## Structure

- **`schema.sql`** - SQL schema definition
- **`main.go`** - Example queries demonstrating various TomaSQL features
- **`<basic | postgres>/table-definitions.gen.go`** - Generated type-safe table definitions (see go:generate comment in <basic | postgres>/main.go)
- **`basic/tables-graph.gen.go`** - Generated graph of table relationships (see go:generate comment in basic/main.go)

## Running the Examples

From the example-app directory, run:

```bash
go run .
```

If you want to regenerate the table definitions after modifying the schema, run:

```bash
go generate ./...
```

## Generated Table Definitions

The `table-definitions.gen.go` file contains type-safe table definitions that correspond to the database schema. They are generated using the table-def-gen tool:

```bash
go run github.com/sergiobonfiglio/tomasql/cmd/table-def-gen --schema ./schema.sql --package-dir . --package-name main
```