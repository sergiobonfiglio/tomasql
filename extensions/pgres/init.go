package pgres

import "github.com/sergiobonfiglio/tomasql/dialects/pgres"
import "github.com/sergiobonfiglio/tomasql"


func init() {
	tomasql.SetDialect(&pgres.PostgresDialect{})
}