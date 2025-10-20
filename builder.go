// package tomasql provides a type-safe SQL query builder for Go.
//
// TomaSQL allows you to build complex SQL queries using a fluent API with compile-time type checking.
// It supports SELECT, JOINs, WHERE conditions, subqueries, aggregations, and more.
//
// Basic usage:
//
//	// Build query
//	query := tomasql.Select(Users.Id, Users.Name, Users.Email).
//		From(Users).
//		Where(Users.IsActive.EqParam(true)).
//	 	OrderBy(Users.Name.Asc()).
//		Limit(10)
//
//	// Generate SQL
//	sql, params := query.SQL()
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
