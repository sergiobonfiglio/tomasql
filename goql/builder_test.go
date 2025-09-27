package goql

import (
	"testing"

	"github.com/stretchr/testify/require"
)


//go:generate go run ../cmd/table-def-gen --schema ../cmd/table-def-gen/example_schema.sql --package-dir ../goql --table-def-file tables-definitions_test.gen.go --table-graph-file= --goql-import-mode none


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
						sql, _ := NewBuilder().Select(Account.Id, Account.Uuid).SQL()
						return sql
					}(),
				},
				{
					want: "SELECT account.id AS a1, account.uuid",
					got: func() string {
						sql, _ := NewBuilder().Select(Account.Id.As("a1"), Account.Uuid).SQL()
						return sql
					}(),
				},

				{
					want: "SELECT a.id, a.uuid FROM account AS a",
					got: func() string {
						acc := Account.As("a")
						sql, _ := NewBuilder().Select(acc.Id, Account.As("a").Uuid).
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
						sql, _ := NewBuilder().SelectAll().
							From(Account).
							SQL()
						return sql
					}(),
				},
				{
					want: "SELECT account.* FROM account",
					got: func() string {
						sql, _ := NewBuilder().Select(Account.Star()).
							From(Account).
							SQL()
						return sql
					}(),
				},
				{
					want: "SELECT a.* FROM account AS a",
					got: func() string {
						accA := Account.As("a")
						sql, _ := NewBuilder().Select(accA.Star()).
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
						sql, _ := NewBuilder().SelectDistinctAll().
							From(Account).
							SQL()
						return sql
					}(),
				},
				{
					want: "SELECT DISTINCT account.id, account.uuid FROM account",
					got: func() string {
						sql, _ := NewBuilder().SelectDistinct(Account.Id, Account.Uuid).
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
						sql, _ := NewBuilder().Select(Account.Id, Account.Uuid).
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
						sql, _ := NewBuilder().Select(acc.Id, acc.Uuid).
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
						sql, _ := NewBuilder().Select(acc.Id, acc.Uuid, sc.Id).
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
						"WHERE a.id = $1 AND sc.id = $2",
					got: func() string {
						acc := Account.As("a")
						sc := ShoppingCart.As("sc")
						sc2 := ShoppingCart.As("sc2")
						sql, _ := NewBuilder().Select(acc.Id, acc.Uuid, sc.Id).
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
						sql, _ := NewBuilder().Select(Account.Id, Account.Uuid).
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
						sql, _ := NewBuilder().Select(Account.Id, Account.Uuid).
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
					want: "SELECT a.id FROM account AS a WHERE a.id = $1",
					got: func() string {
						acc := Account.As("a")
						sql, _ := NewBuilder().Select(acc.Id).
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
					want: "SELECT a.id FROM account AS a WHERE a.id = $1 ORDER BY a.id ASC, a.uuid DESC",
					got: func() string {
						acc := Account.As("a")
						sql, _ := NewBuilder().Select(acc.Id).
							From(acc).
							Where(acc.Id.EqParam(1)).
							OrderBy(acc.Id.Asc(), acc.Uuid.Desc()).
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
						sql, _ := NewBuilder().Select(acc.Id).
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
						sql, _ := NewBuilder().Select(acc.Id).
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
					want: "SELECT (SELECT c.uuid FROM shopping_cart AS c WHERE c.owner_id = $1) AS cart_uuid, a.id FROM account AS a WHERE a.id = $1 OR a.id = $2",
					got: func() string {
						acc := Account.As("a")
						cart := ShoppingCart.As("c")
						sql, _ := NewBuilder().Select(
							NewBuilder().Select(cart.Uuid).From(cart).Where(cart.OwnerId.EqParam(1)).AsNamedSubQuery("cart_uuid"),
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
						sql, _ := NewBuilder().SelectAll().
							From(NewBuilder().Select(Account.Uuid).From(Account).AsNamedSubQuery("uuids")).
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
						"GROUP BY account.id HAVING COUNT(1) > $1",
					got: func() string {
						sql, _ := NewBuilder().
							Select(Account.Id, Count().As("cnt")).
							From(Account).
							GroupBy(Account.Id).
							Having(Count().GtParam(1)).
							SQL()
						return sql
					}(),
				},
				{
					want: "SELECT account.id, COUNT(1) AS cnt FROM account " +
						"GROUP BY account.id HAVING COUNT(1) > $1 " +
						"ORDER BY account.id ASC",
					got: func() string {
						sql, _ := NewBuilder().
							Select(Account.Id, Count().As("cnt")).
							From(Account).
							GroupBy(Account.Id).
							Having(Count().GtParam(1)).
							OrderBy(Account.Id.Asc()).
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

// func TestBuilder_Subquery(t *testing.T) {
// 	xtest := xtestSetup.New(t, []any{})
// 	defer xtest.Cleanup()

// 	type test struct {
// 		got  string
// 		want string
// 	}

// 	t.Run("subqueries with params", func(t *testing.T) {
// 		testItem := test{
// 			want: "SELECT * FROM (SELECT config.uuid FROM config WHERE config.uuid = $1) AS t1 " +
// 				"JOIN (SELECT config.uuid FROM config WHERE config.uuid = $2) AS t2 ON t1.uuid = t2.uuid " +
// 				"JOIN (SELECT config.uuid FROM config WHERE config.uuid = $2) AS t3 ON t2.uuid = t3.uuid",
// 			got: func() string {
// 				t1 := NewBuilder().Select(Config.Uuid).From(Config).Where(Config.Uuid.EqParam("1")).AsNamedSubQuery("t1")
// 				t2 := NewBuilder().Select(Config.Uuid).From(Config).Where(Config.Uuid.EqParam("2")).AsNamedSubQuery("t2")
// 				t3 := NewBuilder().Select(Config.Uuid).From(Config).Where(Config.Uuid.EqParam("2")).AsNamedSubQuery("t3")

// 				t1Col := NewCol[string]("uuid", t1)
// 				t2Col := NewCol[string]("uuid", t2)
// 				t3Col := NewCol[string]("uuid", t3)

// 				sql, params := NewBuilder().SelectAll().
// 					From(t1).
// 					Join(t2).On(t1Col.Eq(t2Col)).
// 					Join(t3).On(t2Col.Eq(t3Col)).
// 					SQL()

// 				require.Len(t, params, 2) // only 2 distinct params
// 				return sql
// 			}(),
// 		}

// 		require.Equal(t, testItem.want, testItem.got)
// 	})
// }

// func TestBuilder_Query(t *testing.T) {
// 	xtest := xtestSetup.New(t, []any{})
// 	defer xtest.Cleanup()

// 	// setup
// 	type accountRow struct {
// 		Id        int64  `db:"id"`
// 		Uuid      string `db:"uuid"`
// 		CreatedTs int64  `db:"created_ts"`
// 	}
// 	var inserted []accountRow
// 	for i := 0; i < 10; i++ {

// 		accUuid := uuid.NewString()
// 		createdTs := time.Now().Unix()
// 		accId := new(int64)
// 		err := xtest.SQL.Get(accId,
// 			`INSERT INTO account (uuid, created_ts) VALUES ($1, $2) RETURNING id`,
// 			accUuid,
// 			createdTs,
// 		)
// 		require.NoError(t, err)
// 		require.NotNil(t, accId)

// 		inserted = append(inserted, accountRow{
// 			Id:        *accId,
// 			Uuid:      accUuid,
// 			CreatedTs: createdTs,
// 		})
// 	}

// 	t.Run("query with params", func(t *testing.T) {
// 		builder := NewBuilder().
// 			Select(Account.Id, Account.Uuid, Account.CreatedTs).
// 			From(Account).
// 			Where(Account.Id.EqParam(inserted[0].Id).
// 				And(Account.Id.EqParam(inserted[0].Id)).
// 				And(Account.Uuid.EqParam(inserted[0].Uuid)))

// 		res := accountRow{}
// 		sql, params := builder.SQL()
// 		require.Len(t, params, 2) // only 2 distinct params

// 		err := xtest.SQL.Get(&res, sql, params...)
// 		require.NoError(t, err)

// 		require.Equal(t, inserted[0].Id, res.Id)
// 		require.Equal(t, inserted[0].Uuid, res.Uuid)
// 		require.Equal(t, inserted[0].CreatedTs, res.CreatedTs)
// 	})
// }
