// package tomasql provides a type-safe SQL query builder for Go.
//
// TomaSQL allows you to build complex SQL queries using a fluent API with compile-time type checking.
// It supports SELECT, JOINs, WHERE conditions, subqueries, aggregations, and more.
//
// Basic usage:
//
//	// Create columns
//	userID := tomasql.NewCol[int]("id", nil)
//	userName := tomasql.NewCol[string]("name", nil)
//
//	// Build query
//	query := tomasql.SelectCols(userID, userName).
//		From(usersTable).
//		Where(userID.GtParam(100))
//
//	// Generate SQL
//	sql, params := query.SQL()
//
// The library provides type safety by associating Go types with database columns,
// ensuring that comparisons and operations are performed with compatible types.
//
// Key features:
//   - Type-safe column definitions and operations
//   - Fluent query building API
//   - Support for complex JOINs and subqueries
//   - Built-in SQL functions (COUNT, SUM, AVG, etc.)
//   - Parameterized queries to prevent SQL injection
//   - Code generation for database schema integration
package tomasql

func SelectCols(first Column, columns ...Column) BuilderWithSelect {
	convertedColumns := make([]ParametricSql, len(columns))
	for i, col := range columns {
		convertedColumns[i] = col
	}

	return newBuilderWithSelect(false, first, convertedColumns...)
}

func Select(first ParametricSql, columns ...ParametricSql) BuilderWithSelect {
	return newBuilderWithSelect(false, first, columns...)
}

func SelectDistinctCols(first Column, columns ...Column) BuilderWithSelect {
	convertedColumns := make([]ParametricSql, len(columns))
	for i, col := range columns {
		convertedColumns[i] = col
	}
	return newBuilderWithSelect(true, first, convertedColumns...)
}

func SelectDistinct(first ParametricSql, column ...ParametricSql) BuilderWithSelect {
	return newBuilderWithSelect(true, first, column...)
}

func SelectAll() BuilderWithSelect {
	return newBuilderWithSelectAll(false)
}

func SelectDistinctAll() BuilderWithSelect {
	return newBuilderWithSelectAll(true)
}
