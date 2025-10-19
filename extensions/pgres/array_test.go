package pgres

import (
	"testing"

	"github.com/sergiobonfiglio/tomasql"
	"github.com/sergiobonfiglio/tomasql/dialects/pgres"
	"github.com/stretchr/testify/require"
)

//go:generate go run ../../cmd/table-def-gen --schema ../../cmd/table-def-gen/example_schema.sql --package-dir ../pgres --table-def-file tables-definitions_gen_test.go --table-graph-file= --with-pgres-extensions

func TestArray(t *testing.T) {

	pgres.SetDialect()

	type test struct {
		got  string
		want string
	}

	tests := []test{
		{
			want: "account.id IN ($1)",
			got:  Account.Id.InArray([]int64{1}).SQL(tomasql.ParamsMap{}),
		},
		{
			want: "account.id IN ($1, $1, $2, $3, $1)",
			got:  Account.Id.InArray([]int64{1, 1, 2, 3, 1}).SQL(tomasql.ParamsMap{}),
		},
		{
			want: "account.id = ANY($1)",
			got:  Account.Id.EqAny(Array([]int64{1})).SQL(tomasql.ParamsMap{}),
		},
		{
			want: "account.id > ANY($1)",
			got:  Account.Id.GtAny(Array([]int64{1, 1, 2, 3, 1})).SQL(tomasql.ParamsMap{}),
		},
		{
			want: "account.id = ALL($1)",
			got:  Account.Id.EqAll(Array([]int64{1})).SQL(tomasql.ParamsMap{}),
		},
		{
			want: "account.id > ALL($1)",
			got:  Account.Id.GtAll(Array([]int64{1, 1, 2, 3, 1})).SQL(tomasql.ParamsMap{}),
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.want, func(tt *testing.T) {
			require.Equal(tt, testItem.want, testItem.got)
		})
	}

}
