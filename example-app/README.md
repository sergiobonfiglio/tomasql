# GoQL Test Application

This is a demonstration application showing how to use the GoQL library with custom table definitions.

## Structure

- **`schema.sql`** - PostgreSQL schema definition
- **`main.go`** - Example queries demonstrating various GoQL features
- **`table-definitions.gen.go`** - Generated type-safe table definitions (see go:generate comment in main.go)
- **`tables-graph.gen.go`** - Generated graph of table relationships (see go:generate comment in main.go)

## Running the Examples

From the test-app directory, run:

```bash
# Build and run the test application
go run .
```

If you want to regenerate the table definitions after modifying the schema, run:

```bash
go generate ./...
```

## Generated Table Definitions

The `table-definitions.gen.go` file contains type-safe table definitions that correspond to the database schema. They are generated using the table-def-gen tool:

```bash
go run github.com/sergiobonfiglio/goql/cmd/table-def-gen --schema ./schema.sql --package-dir . --package-name main
```

## Example Queries

The application demonstrates:

1. **Basic SELECT** - Simple column selection
2. **SELECT with FROM** - Explicit table specification  
3. **WHERE clauses** - Parameterized conditions
4. **JOIN operations** - Table relationships
5. **Complex queries** - Aliases, multiple conditions, ordering
6. **Aggregation** - COUNT functions

## Key Features Demonstrated

- **Type Safety** - Column types are enforced at compile time
- **Parameterized Queries** - SQL injection prevention
- **Fluent API** - Readable query building
- **Table Aliases** - Clean complex queries
- **Join Operations** - Relationship handling
- **Function Support** - Built-in SQL functions

## Integration with GoQL Library

This test application shows the recommended pattern for using GoQL:

1. Define your database schema in SQL
2. Generate type-safe table definitions 
3. Use the fluent query builder for type-safe operations
4. Execute generated SQL with parameters

The library handles SQL generation while providing compile-time guarantees about query correctness.