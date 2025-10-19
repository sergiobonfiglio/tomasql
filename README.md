# TomaSQL - Type-safe SQL Query Builder for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/sergiobonfiglio/tomasql.svg)](https://pkg.go.dev/github.com/sergiobonfiglio/tomasql)
[![Go Report Card](https://goreportcard.com/badge/github.com/sergiobonfiglio/tomasql)](https://goreportcard.com/report/github.com/sergiobonfiglio/tomasql)

TomaSQL is a type-safe SQL query builder for Go that provides a fluent API for constructing SQL queries with compile-time type checking.

## Features

- **Type-safe**: Compile-time type checking for SQL queries
- **Performance**: No reflection
- **Fluent API**: Intuitive, chainable query building
- **Rich SQL Support**: JOINs, subqueries, aggregations, and more
- **Database Schema Integration**: Generate type-safe table definitions from your database schema

## Quick Start

1. Get the library and code generation tool:

   ```bash
   go get github.com/sergiobonfiglio/tomasql
   ```

2. Generate table definitions from your schema directly through the command line:

   ```bash
   # First create your schema.sql file with your database schema
   go run github.com/sergiobonfiglio/tomasql/cmd/table-def-gen@latest --schema ./schema.sql --package-dir . --package-name main
   ```

   or a `go:generate` comment:

   ```go
   //go:generate go run github.com/sergiobonfiglio/tomasql/cmd/table-def-gen@latest --schema ./schema.sql --package-dir . --package-name main
   ```

3. Use the generated tables:

   ```go
   package main

   import (
       "fmt"
       "github.com/sergiobonfiglio/tomasql"
   )

   func main() {
       // Use generated table definitions (Users, Products, etc.)
       query := tomasql.Select(Users.Id, Users.Name, Users.Email).
           From(Users).
           Where(Users.IsActive.EqParam(true)).
           OrderBy(Users.Name.Asc()).
           Limit(10)

       sql, params := query.SQL()
       fmt.Println("SQL:", sql)
       fmt.Println("Params:", params)

       // Output:
       // SQL:
       // SELECT users.id, users.name, users.email
       // FROM users
       // WHERE users.is_active = ?
       // ORDER BY users.name ASC
       // LIMIT 10
       // Params: [true]
   }
   ```

### Working with JOINs

```go
query := tomasql.SelectCols(Users.Name, Posts.Title).
    From(Users).
    Join(Posts).On(Users.Id.Eq(Posts.UserId)).
    Where(Users.Name.IsNotNull()).
    OrderBy(Users.Name.Asc())

sql, params := query.SQL()
// SQL:
// SELECT users.name, posts.title
// FROM users
// JOIN posts ON users.id = posts.user_id
// WHERE users.name IS NOT NULL
// ORDER BY users.name ASC
```

### Using Functions and Aggregations

```go
import "github.com/sergiobonfiglio/tomasql"

// Count users by status
totalUsers := tomasql.Count().As("total_users")
query := tomasql.Select(
        Users.Status,
        totalUsers,
        tomasql.Avg[float64](userAge).As("avg_age"),
    ).
    From(Users).
    GroupBy(Users.Status).
    Having(tomasql.Count().GtParam(10)).
    OrderBy(tomasql.Count().Desc())

sql, params := query.SQL()
// SQL:
// SELECT users.status, COUNT(*) AS total_users, AVG(users.age) AS avg_age
// FROM users
// GROUP BY users.status
// HAVING COUNT(*) > ?
// ORDER BY COUNT(*) DESC
```

### Subqueries

```go
// Subquery example
subQuery := tomasql.SelectCols(Users.Id).
    From(Users).
    Where(Users.Age.GeParam(18))

mainQuery := tomasql.SelectAll().
    From(Posts).
    Where(Posts.UserId.In(subQuery))

sql, params := mainQuery.SQL()
// SQL:
// SELECT *
// FROM posts
// WHERE posts.user_id IN (
//      SELECT users.id FROM users WHERE users.age >= ?)
```

### Working with Complex Conditions

```go
query := tomasql.SelectAll().
    From(Users).
    Where(User.Name.LikeParam("%admin%").
        Or(User.Age.GtParam(25).And(User.Status.EqParam("active")))

sql, params := query.SQL()
// SQL:
// SELECT *
// FROM users
// WHERE users.name LIKE ? OR (users.age > ? AND users.status = ?)
```

## Table Definition Generation

TomaSQL includes a code generation tool to create type-safe table definitions from your database schema:

1. **Create your database schema** (`schema.sql`)
2. **Generate table definitions**:

   ```bash
   # Install the code generation tool
   go install github.com/sergiobonfiglio/tomasql/cmd/table-def-gen@latest

   # Generate table definitions (customize for your database)
   table-def-gen -schema schema.sql -output tables.go
   ```

3. **Use generated type-safe tables**:

   ```go
   // Generated table definitions provide full type safety
   type UsersTableDef struct {
       *tomasql.SqlableTable
       alias     *string
       Id        *tomasql.Col[int]
       Name      *tomasql.Col[string]
       Email     *tomasql.Col[string]
       IsActive  *tomasql.Col[bool]
   }

   var Users = newUsersTable()

   // Usage with generated tables
   query := tomasql.Select(Users.Name, Users.Email).
       From(Users).
       Where(Users.IsActive.EqParam(true))
   ```

> **Note**: See `example-app/` directory for a complete working example with generated tables.

## API Reference

### Entry Points

- `Select(cols ...ParametricSql)` - Start a SELECT query
- `SelectAll()` - Start a SELECT \* query (equivalent to `Select(<GenTable>.Star())`
- `SelectDistinct(cols ...ParametricSql)` - Start a SELECT DISTINCT query

Every entry point also has an alternative version which takes `Column` as parameters instead of `ParametricSql` to avoid manual casting if you have an array of columns you want to select.

### Comparisons

The standard Column implementation also provides type-safe comparison methods:

- `Eq(value T)` -> =
- `Neq(value T)` -> <>
- `Gt(value T)` -> >
- `Gte(value T)` -> >=
- `Lt(value T)` -> <
- `Lte(value T)` -> <=
- `Like(value ParametricSql)` -> LIKE
- `IsNull()` -> IS NULL
- `IsNotNull()` -> IS NOT NULL

Each comparison method also has a `*Param` variant that takes a value and generates a parameter placeholder for it, e.g. `EqParam(value T)`.

### SQL Functions

- `Count()`, `Sum[T]()`, `Avg[T]()`, `Min[T]()`, `Max[T]()`
- `Upper()`, `Lower()`, `Length()`, `Trim()`
- `Coalesce[T]()`, `Round()`, `Abs[T]()`
- `Exists(ParametricSql)`, `Any(ParametricSql)`, `All(ParametricSql)`, `In(ParametricSql)`

## Dialects

TomaSQL is designed to support multiple SQL dialects, which you can set globally.
Dialects are implemented in the `dialects` package and they can be set using:

```go
import (
    "github.com/sergiobonfiglio/tomasql"
    "github.com/sergiobonfiglio/tomasql/dialects/pgres"
)

func init() {
    // example of setting Postgres dialect
    pgres.SetDialect()
    // or, equivalently:
    tomasql.SetDialect(pgres.GetDialect())
}
```

### Extensions

Also, some dialects provide extensions with additional column types and methods. Since these extensions can introduce additional dependencies, they are defined in different modules. You can import them explicitly:

```go
import _ "github.com/sergiobonfiglio/tomasql/extensions/pgres"
```

To enable PostgreSQL-specific features (like array support or ILIKE) you need to either: wrap table columns with the extension column types manually, e.g. `pgres.Wrap(...)`, or generate the table definitions with the `--with-pgres-extensions` flag enabled in the `table-def-gen` tool. See the example-app for a complete example.


## Example Application

The repository includes a complete example application demonstrating TomaSQL usage:

```bash
# Run the example application
go run ./example-app
```