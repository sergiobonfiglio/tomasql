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
//	query := tomasql.NewBuilder().
//		SelectCols(userID, userName).
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

type builder1 struct{}

func (b *builder1) SelectCols(first Column, columns ...Column) BuilderWithSelect {
	convertedColumns := make([]ParametricSql, len(columns))
	for i, col := range columns {
		convertedColumns[i] = col
	}

	return newBuilderWithSelect(false, first, convertedColumns...)
}

func (b *builder1) Select(first ParametricSql, columns ...ParametricSql) BuilderWithSelect {
	return newBuilderWithSelect(false, first, columns...)
}

func (b *builder1) SelectDistinctCols(first Column, columns ...Column) BuilderWithSelect {
	convertedColumns := make([]ParametricSql, len(columns))
	for i, col := range columns {
		convertedColumns[i] = col
	}
	return newBuilderWithSelect(true, first, convertedColumns...)
}

func (b *builder1) SelectDistinct(first ParametricSql, column ...ParametricSql) BuilderWithSelect {
	return newBuilderWithSelect(true, first, column...)
}

func (b *builder1) SelectAll() BuilderWithSelect {
	return newBuilderWithSelectAll(false)
}

func (b *builder1) SelectDistinctAll() BuilderWithSelect {
	return newBuilderWithSelectAll(true)
}

func NewBuilder() Builder1 {
	return &builder1{}
}
