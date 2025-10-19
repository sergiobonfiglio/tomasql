package pgres

import (
	"fmt"

	"github.com/sergiobonfiglio/tomasql"
)

type PostgresDialect struct{}

var _ tomasql.Dialect = (*PostgresDialect)(nil)

func (p *PostgresDialect) Name() string {
	return "postgres"
}

// ArrayToSQL implements tomasql.Dialect.
func (p *PostgresDialect) ArrayToSQL(array any) tomasql.ParametricSql {
	panic("unimplemented")
}

// Placeholder implements tomasql.Dialect.
func (p *PostgresDialect) Placeholder(position int) string {
	return fmt.Sprintf("$%d", position)
}