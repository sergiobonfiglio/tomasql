package tomasql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
				// {
				// 	want: "account.id IN ($1)",
				// 	got:  Account.Id.InArray([]int64{1}).SQL(ParamsMap{}),
				// },
				// {
				// 	want: "account.id IN ($1, $1, $2, $3, $1)",
				// 	got:  Account.Id.InArray([]int64{1, 1, 2, 3, 1}).SQL(ParamsMap{}),
				// },
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
				// {
				// 	want: "account.id = ANY($1)",
				// 	got:  Account.Id.EqAny(SQLArray([]int64{1})).SQL(ParamsMap{}),
				// },
				// {
				// 	want: "account.id > ANY($1)",
				// 	got:  Account.Id.GtAny(SQLArray([]int64{1, 1, 2, 3, 1})).SQL(ParamsMap{}),
				// },
			},
		},
		{
			name: "all",
			tests: []test{
				// {
				// 	want: "account.id = ALL($1)",
				// 	got:  Account.Id.EqAll(SQLArray([]int64{1})).SQL(ParamsMap{}),
				// },
				// {
				// 	want: "account.id > ALL($1)",
				// 	got:  Account.Id.GtAll(SQLArray([]int64{1, 1, 2, 3, 1})).SQL(ParamsMap{}),
				// },
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

	for _, group := range testGroups {
		t.Run(group.name, func(tg *testing.T) {
			for _, testItem := range group.tests {
				tg.Run(testItem.want, func(tt *testing.T) {
					require.Equal(tt, testItem.want, testItem.got)
				})
			}
		})
	}
}

func TestCondition_Columns(t *testing.T) {
	type test struct {
		name string
		impl Condition
		want []Column
	}

	tests := []test{
		{
			name: "identity condition columns",
			impl: IdentityCond,
			want: []Column{NewCol[string]("1", nil), NewCol[string]("1", nil)},
		},
		{
			name: "binary condition columns",
			impl: newBinaryCondition(NewCol[int64]("col1", nil), NewCol[int64]("col2", nil), comparerEq),
			want: []Column{NewCol[int64]("col1", nil), NewCol[int64]("col2", nil)},
		},
		{
			name: "binary condition with param columns",
			impl: newBinaryParamCondition(NewCol[int64]("col1", nil), int64(42), comparerEq),
			want: []Column{NewCol[int64]("col1", nil)},
		},
		// {
		// 	name: "in array condition columns",
		// 	impl: newInArrayCondition(NewCol[int64]("col1", nil), []int64{1, 2, 3}),
		// 	want: []Column{NewCol[int64]("col1", nil)},
		// },
		{
			name: "in subquery condition columns",
			impl: newInCondition(NewCol[int64]("col1", nil), Select(NewCol[int64]("col2", nil))),
			want: []Column{NewCol[int64]("col1", nil)},
		},
		{
			name: "is condition columns",
			impl: newIsCondition(NewCol[int64]("col1", nil), comparerNull),
			want: []Column{NewCol[int64]("col1", nil)},
		},
		{
			name: "grouped condition columns",
			impl: Grouped(newBinaryCondition(NewCol[int64]("col1", nil), NewCol[int64]("col2", nil), comparerEq)),
			want: []Column{NewCol[int64]("col1", nil), NewCol[int64]("col2", nil)},
		},
		{
			name: "exists condition columns",
			impl: NewExistsCondition(Select(NewCol[int64]("col1", nil))),
			want: []Column{},
		},
		{
			name: "and condition columns",
			impl: newBinaryCondition(NewCol[int64]("col1", nil), NewCol[int64]("col2", nil), comparerEq).
				And(newBinaryCondition(NewCol[int64]("col3", nil), NewCol[int64]("col4", nil), comparerEq)),
			want: []Column{
				NewCol[int64]("col1", nil),
				NewCol[int64]("col2", nil),
				NewCol[int64]("col3", nil),
				NewCol[int64]("col4", nil),
			},
		},
		{
			name: "or condition columns",
			impl: newBinaryCondition(NewCol[int64]("col1", nil), NewCol[int64]("col2", nil), comparerEq).
				Or(newBinaryCondition(NewCol[int64]("col3", nil), NewCol[int64]("col4", nil), comparerEq)),
			want: []Column{
				NewCol[int64]("col1", nil),
				NewCol[int64]("col2", nil),
				NewCol[int64]("col3", nil),
				NewCol[int64]("col4", nil),
			},
		},
	}

	for _, testItem := range tests {
		got := testItem.impl.Columns()
		name := testItem.impl.SQL(ParamsMap{})
		t.Run(testItem.name+"_"+name, func(tt *testing.T) {
			require.ElementsMatch(tt, testItem.want, got)
		})
	}
}
