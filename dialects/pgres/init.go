package pgres

import "github.com/sergiobonfiglio/tomasql"


func init() {
	tomasql.SetDialect(&PostgresDialect{})
}