package tomasql

type Dialect interface {

	// // QuoteIdentifier quotes table/column names (e.g., users -> "users" or `users`)
	// QuoteIdentifier(identifier string) string

	// Placeholder returns the parameter placeholder for position n (e.g., $1, ?, :1)
	Placeholder(position int) string

	// ArrayToSQL converts Go arrays to SQL array representation
	ArrayToSQL(array any) ParametricSql
}

// DefaultDialect is used when no dialect is specified
var DefaultDialect Dialect = &standardDialect{}

var dialect Dialect = DefaultDialect

// SetDialect allows changing the dialect
func SetDialect(d Dialect) {
	dialect = d
}

// GetDialect returns the current dialect
func GetDialect() Dialect {
	return dialect
}

type standardDialect struct {
}

var _ Dialect = (*standardDialect)(nil)

func (d *standardDialect) Placeholder(_ int) string {
	return "?"
}

func (d *standardDialect) ArrayToSQL(array any) ParametricSql {
	panic("arrays are not supported in standard SQL")
}
