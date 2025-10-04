# GoQL - Type-safe SQL Query Builder for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/sergiobonfiglio/tomasql.svg)](https://pkg.go.dev/github.com/sergiobonfiglio/tomasql)
[![Go Report Card](https://goreportcard.com/badge/github.com/sergiobonfiglio/tomasql)](https://goreportcard.com/report/github.com/sergiobonfiglio/tomasql)

GoQL is a type-safe SQL query builder for Go that provides a fluent API for constructing SQL queries with compile-time type checking and excellent IDE support.

## Features

- 🔒 **Type-safe**: Compile-time type checking for SQL queries
- ⚡ **Performance**: Zero reflection
- 🎯 **Fluent API**: Intuitive, chainable query building
- 📊 **Rich SQL Support**: JOINs, subqueries, aggregations, and more
- 🗃️ **Database Schema Integration**: Generate type-safe table definitions from your database schema

## Installation

```bash
go get github.com/sergiobonfiglio/tomasql
```

## Quick Start

There are two ways to use GoQL:

1. **With Code Generation (Recommended)** - Generate type-safe table definitions from your database schema
2. **Manual Definition** - Create table and column definitions manually

### Option 1: With Code Generation (Recommended)

1. Install the library and code generation tool:
   ```bash
   go get github.com/sergiobonfiglio/tomasql
   ```

2. Generate table definitions from your schema directly through the command line:
   ```bash
   # First create your schema.sql file with your database schema
   go run github.com/sergiobonfiglio/tomasql/cmd/table-def-gen --schema ./schema.sql --package-dir . --package-name main

   ```
   or a go:generate comment in your main.go file:
   ```go
   //go:generate go run github.com/sergiobonfiglio/tomasql/cmd/table-def-gen --schema ./schema.sql --package-dir . --package-name main
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
       query := goql.NewBuilder().
           Select(Users.Id, Users.Name, Users.Email).
           From(Users).
           Where(Users.IsActive.EqParam(true)).
           OrderBy(Users.Name.Asc()).
           Limit(10)
       
       sql, params := query.SQL()
       fmt.Println("SQL:", sql)
       fmt.Println("Params:", params)
       
       // Output:
       // SQL: SELECT users.id, users.name, users.email FROM users WHERE users.is_active = $1 ORDER BY users.name ASC LIMIT 10
       // Params: [true]
   }
   ```

### Option 2: Manual Definition

```go
package main

import (
    "fmt"
    "github.com/sergiobonfiglio/tomasql"
)

func main() {
    // Create columns manually
    userID := goql.NewCol[int]("id", nil)
    userName := goql.NewCol[string]("name", nil) 
    userEmail := goql.NewCol[string]("email", nil)
    
    // Build a query
    query := goql.NewBuilder().
        Select(userID, userName, userEmail).
        Where(userID.GtParam(100).And(userName.LikeParam("%john%"))).
        OrderBy(userName.Asc()).
        Limit(10)
    
    sql, params := query.SQL()
    fmt.Println("SQL:", sql)
    fmt.Println("Params:", params)
}
```

### Working with JOINs

```go
// Define table relationships
usersTable := goql.NewCol[int]("users", nil)
postsTable := goql.NewCol[int]("posts", nil)

userID := goql.NewCol[int]("id", usersTable)
userName := goql.NewCol[string]("name", usersTable)
postID := goql.NewCol[int]("id", postsTable)
postTitle := goql.NewCol[string]("title", postsTable)
postUserID := goql.NewCol[int]("user_id", postsTable)

query := goql.NewBuilder().
    SelectCols(userName, postTitle).
    From(usersTable).
    Join(postsTable).On(userID.Eq(postUserID)).
    Where(userName.IsNotNull()).
    OrderBy(userName.Asc())

sql, params := query.SQL()
// SQL: SELECT users.name, posts.title FROM users JOIN posts ON users.id = posts.user_id WHERE users.name IS NOT NULL ORDER BY users.name ASC
```

### Using Functions and Aggregations

```go
import "github.com/sergiobonfiglio/tomasql"

// Count users by status
query := goql.NewBuilder().
    Select(
        userStatus,
        goql.Count().As("total_users"),
        goql.Avg[float64](userAge).As("avg_age"),
    ).
    From(usersTable).
    Where(userStatus.InArray([]string{"active", "pending"})).
    GroupBy(userStatus).
    Having(goql.Count().GtParam(10)).
    OrderBy(goql.Count().Desc())

sql, params := query.SQL()
```

### Subqueries

```go
// Subquery example
subQuery := goql.NewBuilder().
    SelectCols(userID).
    From(usersTable).
    Where(userAge.GeParam(18))

mainQuery := goql.NewBuilder().
    SelectAll().
    From(postsTable).
    Where(postUserID.In(subQuery))

sql, params := mainQuery.SQL()
// SQL: SELECT * FROM posts WHERE user_id IN (SELECT id FROM users WHERE age >= $1)
```

### Working with Complex Conditions

```go
// Complex WHERE conditions
condition1 := userName.LikeParam("%admin%")
condition2 := userAge.GtParam(25).And(userStatus.EqParam("active"))
condition3 := userEmail.IsNotNull()

query := goql.NewBuilder().
    SelectAll().
    From(usersTable).
    Where(condition1.Or(condition2).And(condition3))

sql, params := query.SQL()
```

## Advanced Features

### Custom Column Types

```go
// Define custom column types
type UserStatus string
const (
    UserStatusActive   UserStatus = "active"
    UserStatusInactive UserStatus = "inactive"
)

userStatus := goql.NewCol[UserStatus]("status", usersTable)
query := goql.NewBuilder().
    SelectAll().
    From(usersTable).
    Where(userStatus.EqParam(UserStatusActive))
```

### Working with Arrays and ANY/ALL

```go
// Array operations
userIDs := []int{1, 2, 3, 4, 5}
query := goql.NewBuilder().
    SelectAll().
    From(usersTable).
    Where(userID.InArray(userIDs))

// ANY/ALL operations with subqueries  
subQuery := goql.NewBuilder().SelectCols(postUserID).From(postsTable)
query2 := goql.NewBuilder().
    SelectAll().
    From(usersTable).
    Where(userID.EqAny(subQuery))
```

### Table Definition Generation

GoQL includes a code generation tool to create type-safe table definitions from your database schema:

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
       *goql.SqlableTable
       alias     *string
       Id        *goql.Col[int]
       Name      *goql.Col[string]
       Email     *goql.Col[string]
       IsActive  *goql.Col[bool]
   }
   
   var Users = newUsersTable()
   
   // Usage with generated tables
   query := goql.NewBuilder().
       Select(Users.Name, Users.Email).
       From(Users).
       Where(Users.IsActive.EqParam(true))
   ```

> **Note**: See `example-app/` directory for a complete working example with generated tables.

## API Reference

### Core Interfaces

- `Builder1` - Entry point for query building
- `Column` - Represents a database column with type safety
- `Table` - Represents a database table
- `Condition` - Represents WHERE/HAVING conditions  
- `SQLable` - Can be converted to SQL with parameters

### Key Functions

- `NewBuilder()` - Create a new query builder
- `NewCol[T](name, table)` - Create a typed column
- `NewTableFromSubQuery()` - Create table from subquery

### SQL Functions

- `Count()`, `Sum[T]()`, `Avg[T]()`, `Min[T]()`, `Max[T]()`
- `Upper()`, `Lower()`, `Length()`, `Trim()`
- `Coalesce[T]()`, `Round()`, `Abs[T]()`
- `Exists()`

## Example Application

The repository includes a complete example application demonstrating GoQL usage:

```bash
# Run the example application
go run ./example-app
```

The `example-app/` directory contains:
- **`schema.sql`** - Complete PostgreSQL schema (users, products, orders, etc.)
- **`tables-definitions_test.gen.go`** - Generated table definitions with full type safety
- **`tables-graph.gen.go`** - Generated graph of table relationships
- **`main.go`** - Comprehensive examples showing all GoQL features
- **`README.md`** - Detailed documentation

### Example App Features Demonstrated

- ✅ Basic SELECT queries with type safety
- ✅ Complex JOINs across multiple tables
- ✅ WHERE clauses with parameterized conditions
- ✅ Table aliases and column aliasing
- ✅ Aggregation functions (COUNT, SUM, etc.)
- ✅ ORDER BY with multiple columns
- ✅ Left/Right/Inner JOIN operations

## Examples

Check out the **example-app** directory for comprehensive usage examples including:

- Complex JOINs and subqueries  
- Aggregation functions
- Database schema integration