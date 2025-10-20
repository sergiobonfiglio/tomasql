package pgres

import "github.com/sergiobonfiglio/tomasql"

var pgresDialect *PostgresDialect = &PostgresDialect{}

func GetDialect() *PostgresDialect {
	return pgresDialect
}

// SetDialect sets the Postgres dialect as the current tomasql dialect.
// Equivalent to calling tomasql.SetDialect(&PostgresDialect{}).
func SetDialect() {
	tomasql.SetDialect(&PostgresDialect{})
}
