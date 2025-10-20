package tomasql_test

import (
	"testing"

	. "github.com/sergiobonfiglio/tomasql"
	"github.com/sergiobonfiglio/tomasql/dialects/pgres"
	"github.com/stretchr/testify/require"
)

//go:generate go run ../../cmd/table-def-gen --schema ../../cmd/table-def-gen/example_schema.sql --package-dir ../tomasql_test --table-def-file tables-definitions_gen_test.go --table-graph-file=

func TestConditions(t *testing.T) {
	type test struct {
		got  string
		want string
	}
	type testGroup struct {
		name  string
		tests []test
	}

	testGroups := []testGroup{
		{
			name: "identity",
			tests: []test{
				{
					want: "1 = 1",
					got:  IdentityCond.SQL(ParamsMap{}),
				},
			},
		},
		{
			name: "binary conditions",
			tests: []test{
				{
					want: "account.id = shopping_cart.owner_id",
					got:  Account.Id.Eq(ShoppingCart.OwnerId).SQL(ParamsMap{}),
				},
				{
					want: "a.id = s.owner_id",
					got:  Account.As("a").Id.Eq(ShoppingCart.As("s").OwnerId).SQL(ParamsMap{}),
				},
				{
					want: "account.id > shopping_cart.owner_id",
					got:  Account.Id.Gt(ShoppingCart.OwnerId).SQL(ParamsMap{}),
				},
				{
					want: "account.id >= shopping_cart.owner_id",
					got:  Account.Id.Ge(ShoppingCart.OwnerId).SQL(ParamsMap{}),
				},
				{
					want: "account.id < shopping_cart.owner_id",
					got:  Account.Id.Lt(ShoppingCart.OwnerId).SQL(ParamsMap{}),
				},
				{
					want: "account.id <= shopping_cart.owner_id",
					got:  Account.Id.Le(ShoppingCart.OwnerId).SQL(ParamsMap{}),
				},
			},
		},
		{
			name: "binary conditions param",
			tests: []test{
				{
					want: "account.id = " + GetDialect().Placeholder(1),
					got:  Account.Id.EqParam(1).SQL(ParamsMap{}),
				},
				{
					want: "account.id > " + GetDialect().Placeholder(1),
					got:  Account.Id.GtParam(1).SQL(ParamsMap{}),
				},
				{
					want: "account.id >= " + GetDialect().Placeholder(1),
					got:  Account.Id.GeParam(1).SQL(ParamsMap{}),
				},
				{
					want: "account.id < " + GetDialect().Placeholder(1),
					got:  Account.Id.LtParam(1).SQL(ParamsMap{}),
				},
				{
					want: "account.id <= " + GetDialect().Placeholder(1),
					got:  Account.Id.LeParam(1).SQL(ParamsMap{}),
				},
			},
		},
		{
			name: "in",
			tests: []test{
				{
					want: "account.id IN (SELECT shopping_cart.owner_id FROM shopping_cart)",
					got:  Account.Id.In(Select(ShoppingCart.OwnerId).From(ShoppingCart).AsSubQuery()).SQL(ParamsMap{}),
				},
			},
		},
		{
			name: "exists",
			tests: []test{
				{
					want: "EXISTS(SELECT 1)",
					got:  NewExistsCondition(Select(NewFixedCol(1, nil))).SQL(ParamsMap{}),
				},
			},
		},
		{
			name: "any",
			tests: []test{
				{
					want: "account.id = ANY(SELECT shopping_cart.owner_id FROM shopping_cart)",
					got:  Account.Id.EqAny(Select(ShoppingCart.OwnerId).From(ShoppingCart).AsSubQuery()).SQL(ParamsMap{}),
				},
			},
		},
		{
			name: "all",
			tests: []test{
				{
					want: "account.id = ALL(SELECT shopping_cart.owner_id FROM shopping_cart)",
					got:  Account.Id.EqAll(Select(ShoppingCart.OwnerId).From(ShoppingCart).AsSubQuery()).SQL(ParamsMap{}),
				},
			},
		},
		{
			name: "multiple params",
			tests: []test{
				{
					want: "account.id = " + GetDialect().Placeholder(1) + " AND account.id = " + GetDialect().Placeholder(1),
					got:  Account.Id.EqParam(1).And(Account.Id.EqParam(1)).SQL(ParamsMap{}),
				},
				{
					want: "account.id = " + GetDialect().Placeholder(1) + " AND account.id = " + GetDialect().Placeholder(2),
					got:  Account.Id.EqParam(7).And(Account.Id.EqParam(1)).SQL(ParamsMap{}),
				},
			},
		},
		{
			name: "AND/OR concatenation",
			tests: []test{
				{
					want: "account.id = " + GetDialect().Placeholder(1) + " AND account.id = " + GetDialect().Placeholder(1),
					got:  Account.Id.EqParam(1).And(Account.Id.EqParam(1)).SQL(ParamsMap{}),
				},
				{
					want: "account.id = " + GetDialect().Placeholder(1) + " OR account.id = " + GetDialect().Placeholder(2),
					got:  Account.Id.EqParam(7).Or(Account.Id.EqParam(1)).SQL(ParamsMap{}),
				},
				{
					want: "account.id = " + GetDialect().Placeholder(1) + " AND account.uuid = " + GetDialect().Placeholder(2) + " OR account.created_ts = " + GetDialect().Placeholder(3),
					got: Account.Id.EqParam(1).
						And(Account.Uuid.EqParam("abc")).
						Or(Account.CreatedTs.EqParam(3)).
						SQL(ParamsMap{}),
				},
			},
		},
		{
			name: "grouping",
			tests: []test{
				{
					want: "(account.id = " + GetDialect().Placeholder(1) + " AND account.id = " + GetDialect().Placeholder(2) + ")",
					got:  Grouped(Account.Id.EqParam(1).And(Account.Id.EqParam(2))).SQL(ParamsMap{}),
				},
				{
					want: "account.id = " + GetDialect().Placeholder(1) + " AND (account.id = " + GetDialect().Placeholder(2) + ")",
					got:  Account.Id.EqParam(1).And(Grouped(Account.Id.EqParam(2))).SQL(ParamsMap{}),
				},
				{
					want: "(account.id = " + GetDialect().Placeholder(1) + ") AND account.id = " + GetDialect().Placeholder(2),
					got:  Grouped(Account.Id.EqParam(1)).And(Account.Id.EqParam(2)).SQL(ParamsMap{}),
				},
				{
					want: "account.id = " + GetDialect().Placeholder(1) + " AND (account.uuid = " + GetDialect().Placeholder(2) + " OR account.created_ts = " + GetDialect().Placeholder(3) + ")",
					got: Account.Id.EqParam(1).And(
						Grouped(Account.Uuid.EqParam("abc").Or(Account.CreatedTs.EqParam(3))),
					).SQL(ParamsMap{}),
				},
			},
		},
	}

	dialects := []Dialect{
		DefaultDialect,
		pgres.GetDialect(),
	}

	for _, dialect := range dialects {
		t.Run("dialect_"+dialect.Name(), func(td *testing.T) {
			SetDialect(dialect)

			for _, group := range testGroups {
				td.Run(group.name, func(tg *testing.T) {
					for _, testItem := range group.tests {
						tg.Run(testItem.want, func(tt *testing.T) {
							require.Equal(tt, testItem.want, testItem.got)
						})
					}
				})
			}
		})
	}
}
