package goql

import (
	"github.com/stretchr/testify/require"
	"testing"
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
			name: "binary conditions",
			tests: []test{
				{
					want: "account.id = config.account_id",
					got:  Account.Id.Eq(Config.AccountId).SQL(ParamsMap{}),
				},
				{
					want: "a.id = c.account_id",
					got:  Account.As("a").Id.Eq(Config.As("c").AccountId).SQL(ParamsMap{}),
				},
				{
					want: "account.id > config.account_id",
					got:  Account.Id.Gt(Config.AccountId).SQL(ParamsMap{}),
				},
				{
					want: "account.id >= config.account_id",
					got:  Account.Id.Ge(Config.AccountId).SQL(ParamsMap{}),
				},
				{
					want: "account.id < config.account_id",
					got:  Account.Id.Lt(Config.AccountId).SQL(ParamsMap{}),
				},
				{
					want: "account.id <= config.account_id",
					got:  Account.Id.Le(Config.AccountId).SQL(ParamsMap{}),
				},
			},
		},
		{
			name: "binary conditions param",
			tests: []test{
				{
					want: "account.id = $1",
					got:  Account.Id.EqParam(1).SQL(ParamsMap{}),
				},
				{
					want: "account.id > $1",
					got:  Account.Id.GtParam(1).SQL(ParamsMap{}),
				},
				{
					want: "account.id >= $1",
					got:  Account.Id.GeParam(1).SQL(ParamsMap{}),
				},
				{
					want: "account.id < $1",
					got:  Account.Id.LtParam(1).SQL(ParamsMap{}),
				},
				{
					want: "account.id <= $1",
					got:  Account.Id.LeParam(1).SQL(ParamsMap{}),
				},
			},
		},
		{
			name: "in",
			tests: []test{

				{
					want: "account.id IN ($1)",
					got:  Account.Id.InArray([]int64{1}).SQL(ParamsMap{}),
				},
				{
					want: "account.id IN ($1, $1, $2, $3, $1)",
					got:  Account.Id.InArray([]int64{1, 1, 2, 3, 1}).SQL(ParamsMap{}),
				},
				{
					want: "account.id IN (SELECT account.id FROM account)",
					got:  Account.Id.In(newBuilder().Select(Account.Id).From(Account).AsSubQuery()).SQL(ParamsMap{}),
				},
			},
		},
		{
			name: "any",
			tests: []test{

				{
					want: "account.id = ANY($1)",
					got:  Account.Id.EqAny(SQLArray([]int64{1})).SQL(ParamsMap{}),
				},
				{
					want: "account.id = ANY(SELECT account.id FROM account)",
					got:  Account.Id.EqAny(newBuilder().Select(Account.Id).From(Account).AsSubQuery()).SQL(ParamsMap{}),
				},
				{
					want: "account.id > ANY($1)",
					got:  Account.Id.GtAny(SQLArray([]int64{1, 1, 2, 3, 1})).SQL(ParamsMap{}),
				},
			},
		},
		{
			name: "all",
			tests: []test{

				{
					want: "account.id = ALL($1)",
					got:  Account.Id.EqAll(SQLArray([]int64{1})).SQL(ParamsMap{}),
				},
				{
					want: "account.id = ALL(SELECT account.id FROM account)",
					got:  Account.Id.EqAll(newBuilder().Select(Account.Id).From(Account).AsSubQuery()).SQL(ParamsMap{}),
				},
				{
					want: "account.id > ALL($1)",
					got:  Account.Id.GtAll(SQLArray([]int64{1, 1, 2, 3, 1})).SQL(ParamsMap{}),
				},
			},
		},
		{
			name: "multiple params",
			tests: []test{

				{
					want: "account.id = $1 AND account.id = $1",
					got:  Account.Id.EqParam(1).And(Account.Id.EqParam(1)).SQL(ParamsMap{}),
				},
				{
					want: "account.id = $1 AND account.id = $2",
					got:  Account.Id.EqParam(7).And(Account.Id.EqParam(1)).SQL(ParamsMap{}),
				},
			},
		},
		{
			name: "AND/OR concatenation",
			tests: []test{
				{
					want: "account.id = $1 AND account.id = $1",
					got:  Account.Id.EqParam(1).And(Account.Id.EqParam(1)).SQL(ParamsMap{}),
				},
				{
					want: "account.id = $1 OR account.id = $2",
					got:  Account.Id.EqParam(7).Or(Account.Id.EqParam(1)).SQL(ParamsMap{}),
				},
				{
					want: "account.id = $1 AND account.uuid = $2 OR account.created_ts = $3",
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
					want: "(account.id = $1 AND account.id = $2)",
					got:  Grouped(Account.Id.EqParam(1).And(Account.Id.EqParam(2))).SQL(ParamsMap{}),
				},
				{
					want: "account.id = $1 AND (account.id = $2)",
					got:  Account.Id.EqParam(1).And(Grouped(Account.Id.EqParam(2))).SQL(ParamsMap{}),
				},
				{
					want: "(account.id = $1) AND account.id = $2",
					got:  Grouped(Account.Id.EqParam(1)).And(Account.Id.EqParam(2)).SQL(ParamsMap{}),
				},
				{
					want: "account.id = $1 AND (account.uuid = $2 OR account.created_ts = $3)",
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
