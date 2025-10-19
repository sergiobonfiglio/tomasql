# Table Definition Generator

This tool generates type-safe Go table definitions from a database schema definition SQL file.

## Usage

The generator creates Go code that provides compile-time type safety for your database tables and columns.

### Quick Start

1. Run the generator:

```bash
go run github.com/sergiobonfiglio/tomasql/cmd/table-def-gen --schema ./path/to/your/schema.sql --package-dir ./path/to/your/package --package-name yourpkg
```

### Generated Code Example

From a SQL table definition like:

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    age INTEGER,
    created_at TIMESTAMP DEFAULT NOW()
);
```

The generator creates:

```go
type UsersTableDef struct {
    *sqlableTable
    alias *string
    ID        *Col[int]
    Name      *Col[string]
    Email     *Col[string]
    Age       *Col[int]
    CreatedAt *Col[time.Time]
}

var Users = newUsersTable()

func (u *UsersTableDef) TableName() string {
    return "users"
}

func (u *UsersTableDef) As(alias string) *UsersTableDef {
    // ... aliasing logic
}
```

### Usage with TomaSQL

```go
import "github.com/sergiobonfiglio/tomasql"

// Use generated table definitions
query := tomasql.SelectCols(Users.Name, Users.Email).
    From(Users).
    Where(Users.Age.GtParam(18).And(Users.Email.IsNotNull())).
    OrderBy(Users.Name.Asc())

sql, params := query.SQL()
```

## Configuration

The generator currently supports PostgreSQL databases and uses the following parameters:

| Parameter            | Type   | Required | Default                    | Description                                                                           |
| -------------------- | ------ | -------- | -------------------------- | ------------------------------------------------------------------------------------- |
| `--schema`           | string | Yes      | -                          | Path to the PostgreSQL schema SQL file.                                               |
| `--package-dir`      | string | Yes      | -                          | Directory where the generated Go code will be written.                                |
| `--package-name`     | string | No       | (directory name)           | Name of the Go package for the generated code.                                        |
| `--table-def-file`   | string | No       | `table-definitions.gen.go` | Name of the generated table definitions file.                                         |
| `--table-graph-file` | string | No       | `tables-graph.gen.go`      | Name of the generated tables graph file. If empty, graph file won't be generated.     |
| `--tomasql-import-mode` | string | No       | `full`                     | How to import tomasql package: 'full' (tomasql.Type), 'dot' (. import), 'none' (no import). |
| `--postgres-image`   | string | No       | `postgres:latest`          | Postgres Docker image to use for tables generation.                                   |
| `--with-pgres-extensions`   | bool | No       | `false`          | Generates tables with Postgres columns so that Postgres specific methods can be used                                    |
| `--help`             | bool   | No       | `false`                    | Show help message and exit.                                                           |

#### Example

```bash
go run github.com/sergiobonfiglio/tomasql/cmd/table-def-gen \
    --schema ./schema.sql \
    --package-dir ./example-app \
    --package-name exampleapp \
    --table-def-file my-tables.gen.go \
    --table-graph-file my-graph.gen.go \
    --tomasql-import-mode dot
```

## Requirements

- Go 1.23+
- Docker (for test database setup)

## Database Type Mappings

| SQL Type | Go Type     |
| --------------- | ----------- |
| `bool`          | `bool`      |
| `int2`          | `int16`     |
| `int4`          | `int`       |
| `int8`          | `int64`     |
| `float4`        | `float32`   |
| `float8`        | `float64`   |
| `numeric`       | `float64`   |
| `text`          | `string`    |
| `varchar`       | `string`    |
| `string`        | `string`    |
| `uuid`          | `string`    |
| `bpchar`        | `string`    |
| `timestamp`     | `time.Time` |
| `timestamptz`   | `time.Time` |
| `date`          | `time.Time` |

## Generated Features

The generated table definitions provide:

- **Type Safety**: Column types match your database schema
- **Table Aliasing**: Support for table aliases in queries
- **Column References**: Easy access to table columns