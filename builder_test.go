package tomasql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

//go:generate go run ./cmd/table-def-gen --schema ./cmd/table-def-gen/example_schema.sql --package-dir ../tomasql --table-def-file tables-definitions_gen_test.go --table-graph-file= --tomasql-import-mode none

func TestBuilder(t *testing.T) {
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
			name: "select",
			tests: []test{
				{
					want: "SELECT account.id, account.uuid",
					got: func() string {
						sql, _ := Select(Account.Id, Account.Uuid).SQL()
						return sql
					}(),
				},
				{
					want: "SELECT account.id AS a1, account.uuid",
					got: func() string {
						sql, _ := Select(Account.Id.As("a1"), Account.Uuid).SQL()
						return sql
					}(),
				},

				{
					want: "SELECT a.id, a.uuid FROM account AS a",
					got: func() string {
						acc := Account.As("a")
						sql, _ := Select(acc.Id, Account.As("a").Uuid).
							From(acc).
							SQL()
						return sql
					}(),
				},
			},
		},
		{
			name: "select *",
			tests: []test{
				{
					want: "SELECT * FROM account",
					got: func() string {
						sql, _ := SelectAll().
							From(Account).
							SQL()
						return sql
					}(),
				},
				{
					want: "SELECT account.* FROM account",
					got: func() string {
						sql, _ := Select(Account.Star()).
							From(Account).
							SQL()
						return sql
					}(),
				},
				{
					want: "SELECT a.* FROM account AS a",
					got: func() string {
						accA := Account.As("a")
						sql, _ := Select(accA.Star()).
							From(accA).
							SQL()
						return sql
					}(),
				},
			},
		},
		{
			name: "select distinct",
			tests: []test{
				{
					want: "SELECT DISTINCT * FROM account",
					got: func() string {
						sql, _ := SelectDistinctAll().
							From(Account).
							SQL()
						return sql
					}(),
				},
				{
					want: "SELECT DISTINCT account.id, account.uuid FROM account",
					got: func() string {
						sql, _ := SelectDistinct(Account.Id, Account.Uuid).
							From(Account).
							SQL()
						return sql
					}(),
				},
			},
		},
		{
			name: "from",
			tests: []test{
				{
					want: "SELECT account.id, account.uuid FROM account",
					got: func() string {
						sql, _ := Select(Account.Id, Account.Uuid).
							From(Account).
							SQL()
						return sql
					}(),
				},
			},
		},
		{
			name: "join",
			tests: []test{
				{
					want: "SELECT a.id, a.uuid FROM account AS a " +
						"JOIN shopping_cart ON a.id = shopping_cart.owner_id",
					got: func() string {
						acc := Account.As("a")
						sql, _ := Select(acc.Id, acc.Uuid).
							From(acc).
							Join(ShoppingCart).On(acc.Id.Eq(ShoppingCart.OwnerId)).
							SQL()
						return sql
					}(),
				},
				{
					want: "SELECT a.id, a.uuid, sc.id FROM account AS a " +
						"JOIN shopping_cart AS sc ON a.id = sc.owner_id " +
						"JOIN shopping_cart AS sc2 ON a.id = sc2.owner_id",
					got: func() string {
						acc := Account.As("a")
						sc := ShoppingCart.As("sc")
						sc2 := ShoppingCart.As("sc2")
						sql, _ := Select(acc.Id, acc.Uuid, sc.Id).
							From(acc).
							Join(sc).On(acc.Id.Eq(sc.OwnerId)).
							Join(sc2).On(acc.Id.Eq(sc2.OwnerId)).
							SQL()
						return sql
					}(),
				},
				{
					want: "SELECT a.id, a.uuid, sc.id FROM account AS a " +
						"JOIN shopping_cart AS sc ON a.id = sc.owner_id " +
						"JOIN shopping_cart AS sc2 ON a.id = sc2.owner_id " +
						"WHERE a.id = " + GetDialect().Placeholder(1) + " AND sc.id = " + GetDialect().Placeholder(2),
					got: func() string {
						acc := Account.As("a")
						sc := ShoppingCart.As("sc")
						sc2 := ShoppingCart.As("sc2")
						sql, _ := Select(acc.Id, acc.Uuid, sc.Id).
							From(acc).
							Join(sc).On(acc.Id.Eq(sc.OwnerId)).
							Join(sc2).On(acc.Id.Eq(sc2.OwnerId)).
							Where(acc.Id.EqParam(1).And(sc.Id.EqParam(2))).
							SQL()
						return sql
					}(),
				},
				{
					want: "SELECT account.id, account.uuid FROM account " +
						"LEFT JOIN shopping_cart ON account.id = shopping_cart.owner_id",
					got: func() string {
						sql, _ := Select(Account.Id, Account.Uuid).
							From(Account).
							LeftJoin(ShoppingCart).On(Account.Id.Eq(ShoppingCart.OwnerId)).
							SQL()
						return sql
					}(),
				},
				{
					want: "SELECT account.id, account.uuid FROM account " +
						"RIGHT JOIN shopping_cart ON account.id = shopping_cart.owner_id",
					got: func() string {
						sql, _ := Select(Account.Id, Account.Uuid).
							From(Account).
							RightJoin(ShoppingCart).On(Account.Id.Eq(ShoppingCart.OwnerId)).
							SQL()
						return sql
					}(),
				},
			},
		},
		{
			name: "where",
			tests: []test{
				{
					want: "SELECT a.id FROM account AS a WHERE a.id = " + GetDialect().Placeholder(1),
					got: func() string {
						acc := Account.As("a")
						sql, _ := Select(acc.Id).
							From(acc).
							Where(acc.Id.EqParam(1)).
							SQL()
						return sql
					}(),
				},
			},
		},
		{
			name: "order by",
			tests: []test{
				{
					want: "SELECT account.id AS accId, account.uuid FROM account WHERE account.id = account.id ORDER BY accId ASC, account.uuid DESC",
					got: func() string {

						idCol := Account.Id.As("accId")
						sql, _ := Select(idCol, Account.Uuid).
							From(Account).
							Where(idCol.Eq(Account.Id)).
							OrderBy(idCol.Asc(), Account.Uuid.Desc()).
							SQL()
						return sql
					}(),
				},
				{
					want: "SELECT account.id AS accId, account.uuid FROM account WHERE account.id = " + GetDialect().Placeholder(1) + " ORDER BY accId ASC, account.uuid DESC",
					got: func() string {

						idCol := Account.Id.As("accId")
						sql, _ := Select(idCol, Account.Uuid).
							From(Account).
							Where(Account.Id.EqParam(1)).
							OrderBy(idCol.Asc(), Account.Uuid.Desc()).
							SQL()
						return sql
					}(),
				},
				{
					want: "SELECT a.id FROM account AS a WHERE a.id = " + GetDialect().Placeholder(1) + " ORDER BY a.id ASC, a.uuid DESC",
					got: func() string {
						acc := Account.As("a")
						sql, _ := Select(acc.Id).
							From(acc).
							Where(acc.Id.EqParam(1)).
							OrderBy(acc.Id.Asc(), acc.Uuid.Desc()).
							SQL()
						return sql
					}(),
				},
				{
					want: "SELECT account.id, COUNT(1) AS cnt FROM account GROUP BY account.id ORDER BY cnt ASC",
					got: func() string {
						sql, _ := Select(Account.Id, Count().As("cnt")).
							From(Account).
							GroupBy(Account.Id).
							OrderBy(Count().As("cnt").Asc()).
							SQL()
						return sql
					}(),
				},
				{
					want: "SELECT account.id, COUNT(1) FROM account GROUP BY account.id ORDER BY COUNT(1) ASC",
					got: func() string {
						sql, _ := Select(Account.Id, Count()).
							From(Account).
							GroupBy(Account.Id).
							OrderBy(Count().Asc()).
							SQL()
						return sql
					}(),
				},
			},
		},
		{
			name: "limit/offset",
			tests: []test{
				{
					want: "SELECT a.id FROM account AS a ORDER BY a.id ASC LIMIT 10",
					got: func() string {
						acc := Account.As("a")
						sql, _ := Select(acc.Id).
							From(acc).
							OrderBy(acc.Id.Asc()).
							Limit(10).
							SQL()
						return sql
					}(),
				},
				{
					want: "SELECT a.id FROM account AS a ORDER BY a.id ASC LIMIT 10 OFFSET 0",
					got: func() string {
						acc := Account.As("a")
						sql, _ := Select(acc.Id).
							From(acc).
							OrderBy(acc.Id.Asc()).
							Limit(10).Offset(0).
							SQL()
						return sql
					}(),
				},
			},
		},
		{
			name: "subqueries",
			tests: []test{
				{
					want: "SELECT (SELECT c.uuid FROM shopping_cart AS c WHERE c.owner_id = " + GetDialect().Placeholder(1) + ") AS cart_uuid, a.id FROM account AS a WHERE a.id = " + GetDialect().Placeholder(1) + " OR a.id = " + GetDialect().Placeholder(2),
					got: func() string {
						acc := Account.As("a")
						cart := ShoppingCart.As("c")
						sql, _ := Select(
							Select(cart.Uuid).From(cart).Where(cart.OwnerId.EqParam(1)).AsNamedSubQuery("cart_uuid"),
							acc.Id,
						).From(acc).
							Where(acc.Id.EqParam(1).Or(acc.Id.EqParam(2))).
							SQL()
						return sql
					}(),
				},
				{
					want: "SELECT * FROM (SELECT account.uuid FROM account) AS uuids",
					got: func() string {
						sql, _ := SelectAll().
							From(Select(Account.Uuid).From(Account).AsNamedSubQuery("uuids")).
							SQL()
						return sql
					}(),
				},
			},
		},
		{
			name: "group by/having",
			tests: []test{
				{
					want: "SELECT account.id, COUNT(1) AS cnt FROM account " +
						"WHERE account.id > " + GetDialect().Placeholder(1) + " " +
						"GROUP BY account.id HAVING COUNT(1) > " + GetDialect().Placeholder(1),
					got: func() string {
						sql, _ := Select(Account.Id, Count().As("cnt")).
							From(Account).
							Where(Account.Id.GtParam(0)).
							GroupBy(Account.Id).
							Having(Count().GtParam(1)).
							SQL()
						return sql
					}(),
				},
				{
					want: "SELECT account.id, COUNT(1) AS cnt FROM account " +
						"GROUP BY account.id HAVING COUNT(1) > " + GetDialect().Placeholder(1) + " " +
						"ORDER BY account.id ASC",
					got: func() string {
						sql, _ := Select(Account.Id, Count().As("cnt")).
							From(Account).
							GroupBy(Account.Id).
							Having(Count().GtParam(1)).
							OrderBy(Account.Id.Asc()).
							SQL()
						return sql
					}(),
				},
				{
					want: "SELECT account.id, COUNT(account.uuid) AS cnt FROM account " +
						"GROUP BY account.id",
					got: func() string {
						sql, _ := Select(Account.Id, Count(Account.Uuid).As("cnt")).
							From(Account).
							GroupBy(Account.Id).
							SQL()
						return sql
					}(),
				},
				{
					want: "SELECT account.id, COUNT(DISTINCT account.uuid) AS cnt FROM account " +
						"GROUP BY account.id",
					got: func() string {
						sql, _ := Select(Account.Id, CountDistinct(Account.Uuid).As("cnt")).
							From(Account).
							GroupBy(Account.Id).
							SQL()
						return sql
					}(),
				},
				{
					want: "SELECT account.id, COUNT(DISTINCT account.uuid, account.type) AS cnt FROM account " +
						"GROUP BY account.id",
					got: func() string {
						sql, _ := Select(Account.Id, CountDistinct(Account.Uuid, Account.Type).As("cnt")).
							From(Account).
							GroupBy(Account.Id).
							SQL()
						return sql
					}(),
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
