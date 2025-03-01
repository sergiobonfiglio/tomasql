package goql

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

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
						sql, _ := newBuilder().Select(Account.Id, Account.Uuid).SQL()
						return sql
					}(),
				},
				{
					want: "SELECT account.id AS a1, account.uuid",
					got: func() string {
						sql, _ := newBuilder().Select(Account.Id.As("a1"), Account.Uuid).SQL()
						return sql
					}(),
				},

				{
					want: "SELECT a.id, a.uuid FROM account AS a",
					got: func() string {
						acc := Account.As("a")
						sql, _ := newBuilder().Select(acc.Id, Account.As("a").Uuid).
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
						sql, _ := newBuilder().SelectAll().
							From(Account).
							SQL()
						return sql
					}(),
				},
				{
					want: "SELECT config.* FROM config",
					got: func() string {
						sql, _ := newBuilder().Select(Config.Star()).
							From(Config).
							SQL()
						return sql
					}(),
				},
				{
					want: "SELECT a.* FROM account AS a",
					got: func() string {
						accA := Account.As("a")
						sql, _ := newBuilder().Select(accA.Star()).
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
						sql, _ := newBuilder().SelectDistinctAll().
							From(Account).
							SQL()
						return sql
					}(),
				},
				{
					want: "SELECT DISTINCT account.id, account.uuid FROM account",
					got: func() string {
						sql, _ := newBuilder().SelectDistinct(Account.Id, Account.Uuid).
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
						sql, _ := newBuilder().Select(Account.Id, Account.Uuid).
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
						sql, _ := newBuilder().Select(acc.Id, acc.Uuid).
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
						sql, _ := newBuilder().Select(acc.Id, acc.Uuid, sc.Id).
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
						sql, _ := newBuilder().Select(acc.Id, acc.Uuid, sc.Id).
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
						sql, _ := newBuilder().Select(Account.Id, Account.Uuid).
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
						sql, _ := newBuilder().Select(Account.Id, Account.Uuid).
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
						sql, _ := newBuilder().Select(acc.Id).
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
						sql, _ := newBuilder().Select(acc.Id).
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
						sql, _ := newBuilder().Select(acc.Id).
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
						sql, _ := newBuilder().Select(acc.Id).
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
					want: "SELECT (SELECT c.uuid FROM config AS c WHERE c.account_id = $1) AS conf_uuid, a.id FROM account AS a WHERE a.id = $1 OR a.id = $2",
					got: func() string {
						acc := Account.As("a")
						conf := Config.As("c")
						sql, _ := newBuilder().Select(
							newBuilder().Select(conf.Uuid).From(conf).Where(conf.AccountId.EqParam(1)).AsNamedSubQuery("conf_uuid"),
							acc.Id,
						).From(acc).
							Where(acc.Id.EqParam(1).Or(acc.Id.EqParam(2))).
							SQL()
						return sql
					}(),
				},
				{
					want: "SELECT * FROM (SELECT config.uuid FROM config) AS uuids",
					got: func() string {
						sql, _ := newBuilder().SelectAll().
							From(newBuilder().Select(Config.Uuid).From(Config).AsNamedSubQuery("uuids")).
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
						sql, _ := newBuilder().
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
						sql, _ := newBuilder().
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

func TestBuilder_Query(t *testing.T) {

	xtest := xtestSetup.New(t, []any{})
	defer xtest.Cleanup()

	//setup
	type accountRow struct {
		Id        int64  `db:"id"`
		Uuid      string `db:"uuid"`
		CreatedTs int64  `db:"created_ts"`
	}
	var inserted []accountRow
	for i := 0; i < 10; i++ {

		accUuid := uuid.NewString()
		createdTs := time.Now().Unix()
		var accId *int64 = new(int64)
		err := xtest.SQL.Get(accId,
			`INSERT INTO account (uuid, created_ts) VALUES ($1, $2) RETURNING id`,
			accUuid,
			createdTs,
		)
		require.Nil(t, err)
		require.NotNil(t, accId)

		inserted = append(inserted, accountRow{
			Id:        *accId,
			Uuid:      accUuid,
			CreatedTs: createdTs,
		})
	}

	t.Run("query with params", func(t *testing.T) {

		builder := newBuilder().
			Select(Account.Id, Account.Uuid, Account.CreatedTs).
			From(Account).
			Where(Account.Id.EqParam(inserted[0].Id).
				And(Account.Id.EqParam(inserted[0].Id)).
				And(Account.Uuid.EqParam(inserted[0].Uuid)))

		res := accountRow{}
		sql, params := builder.SQL()
		require.Len(t, params, 2) // only 2 distinct params

		err := xtest.SQL.Get(&res, sql, params...)
		require.Nil(t, err)

		require.Equal(t, inserted[0].Id, res.Id)
		require.Equal(t, inserted[0].Uuid, res.Uuid)
		require.Equal(t, inserted[0].CreatedTs, res.CreatedTs)

	})

}
