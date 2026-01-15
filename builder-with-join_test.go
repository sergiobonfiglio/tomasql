package tomasql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuilderWithJoin_On(t *testing.T) {
	t.Run("set join condition", func(t *testing.T) {
		sql, _ := Select(Account.Id).From(Account).
			Join(ShoppingCart).On(Account.Id.Eq(ShoppingCart.OwnerId)).
			SQL()

		expected := "SELECT account.id FROM account JOIN shopping_cart ON account.id = shopping_cart.owner_id"
		require.Equal(t, expected, sql)
	})

	t.Run("multiple joins with conditions", func(t *testing.T) {
		acc := Account.As("a")
		sc := ShoppingCart.As("sc")
		sc2 := ShoppingCart.As("sc2")

		sql, _ := Select(acc.Id).From(acc).
			Join(sc).On(acc.Id.Eq(sc.OwnerId)).
			Join(sc2).On(acc.Id.Eq(sc2.OwnerId)).
			SQL()

		expected := "SELECT a.id FROM account AS a " +
			"JOIN shopping_cart AS sc ON a.id = sc.owner_id " +
			"JOIN shopping_cart AS sc2 ON a.id = sc2.owner_id"
		require.Equal(t, expected, sql)
	})
}

func TestBuilderWithJoin_JoinTypes(t *testing.T) {
	t.Run("inner join", func(t *testing.T) {
		sql, _ := Select(Account.Id).From(Account).
			Join(ShoppingCart).On(Account.Id.Eq(ShoppingCart.OwnerId)).
			SQL()

		expected := "SELECT account.id FROM account JOIN shopping_cart ON account.id = shopping_cart.owner_id"
		require.Equal(t, expected, sql)
	})

	t.Run("left join", func(t *testing.T) {
		sql, _ := Select(Account.Id).From(Account).
			LeftJoin(ShoppingCart).On(Account.Id.Eq(ShoppingCart.OwnerId)).
			SQL()

		expected := "SELECT account.id FROM account LEFT JOIN shopping_cart ON account.id = shopping_cart.owner_id"
		require.Equal(t, expected, sql)
	})

	t.Run("right join", func(t *testing.T) {
		sql, _ := Select(Account.Id).From(Account).
			RightJoin(ShoppingCart).On(Account.Id.Eq(ShoppingCart.OwnerId)).
			SQL()

		expected := "SELECT account.id FROM account RIGHT JOIN shopping_cart ON account.id = shopping_cart.owner_id"
		require.Equal(t, expected, sql)
	})

	t.Run("mixed join types", func(t *testing.T) {
		acc := Account.As("a")
		sc := ShoppingCart.As("sc")

		sql, _ := Select(acc.Id).From(acc).
			LeftJoin(sc).On(acc.Id.Eq(sc.OwnerId)).
			SQL()

		expected := "SELECT a.id FROM account AS a LEFT JOIN shopping_cart AS sc ON a.id = sc.owner_id"
		require.Equal(t, expected, sql)
	})
}

func TestBuilderWithJoin_LeftJoin(t *testing.T) {
	t.Run("single left join", func(t *testing.T) {
		sql, _ := Select(Account.Id).From(Account).
			LeftJoin(ShoppingCart).On(Account.Id.Eq(ShoppingCart.OwnerId)).
			SQL()

		expected := "SELECT account.id FROM account LEFT JOIN shopping_cart ON account.id = shopping_cart.owner_id"
		require.Equal(t, expected, sql)
	})

	t.Run("multiple left joins", func(t *testing.T) {
		sc1 := ShoppingCart.As("sc1")
		sc2 := ShoppingCart.As("sc2")

		sql, _ := Select(Account.Id).From(Account).
			LeftJoin(sc1).On(Account.Id.Eq(sc1.OwnerId)).
			LeftJoin(sc2).On(Account.Id.Eq(sc2.OwnerId)).
			SQL()

		expected := "SELECT account.id FROM account LEFT JOIN shopping_cart AS sc1 ON account.id = sc1.owner_id LEFT JOIN shopping_cart AS sc2 ON account.id = sc2.owner_id"
		require.Equal(t, expected, sql)
	})
}

func TestBuilderWithJoin_RightJoin(t *testing.T) {
	t.Run("single right join", func(t *testing.T) {
		sql, _ := Select(Account.Id).From(Account).
			RightJoin(ShoppingCart).On(Account.Id.Eq(ShoppingCart.OwnerId)).
			SQL()

		expected := "SELECT account.id FROM account RIGHT JOIN shopping_cart ON account.id = shopping_cart.owner_id"
		require.Equal(t, expected, sql)
	})

	t.Run("multiple right joins", func(t *testing.T) {
		sc1 := ShoppingCart.As("sc1")
		sc2 := ShoppingCart.As("sc2")

		sql, _ := Select(Account.Id).From(Account).
			RightJoin(sc1).On(Account.Id.Eq(sc1.OwnerId)).
			RightJoin(sc2).On(Account.Id.Eq(sc2.OwnerId)).
			SQL()

		expected := "SELECT account.id FROM account RIGHT JOIN shopping_cart AS sc1 ON account.id = sc1.owner_id RIGHT JOIN shopping_cart AS sc2 ON account.id = sc2.owner_id"
		require.Equal(t, expected, sql)
	})
}

func TestBuilderWithJoin_Where(t *testing.T) {
	t.Run("join with where clause", func(t *testing.T) {
		sql, _ := Select(Account.Id).From(Account).
			Join(ShoppingCart).On(Account.Id.Eq(ShoppingCart.OwnerId)).
			Where(Account.Id.EqParam(1)).
			SQL()

		expected := "SELECT account.id FROM account JOIN shopping_cart ON account.id = shopping_cart.owner_id WHERE account.id = ?"
		require.Equal(t, expected, sql)
	})

	t.Run("join with complex where condition", func(t *testing.T) {
		sql, _ := Select(Account.Id).From(Account).
			Join(ShoppingCart).On(Account.Id.Eq(ShoppingCart.OwnerId)).
			Where(Account.Id.EqParam(1).And(ShoppingCart.Id.EqParam(2))).
			SQL()

		expected := "SELECT account.id FROM account JOIN shopping_cart ON account.id = shopping_cart.owner_id WHERE account.id = ? AND shopping_cart.id = ?"
		require.Equal(t, expected, sql)
	})
}

func TestBuilderWithJoin_GroupBy(t *testing.T) {
	t.Run("join with group by", func(t *testing.T) {
		sql, _ := Select(Account.Id, Count().As("cnt")).From(Account).
			Join(ShoppingCart).On(Account.Id.Eq(ShoppingCart.OwnerId)).
			GroupBy(Account.Id).
			SQL()

		expected := "SELECT account.id, COUNT(1) AS cnt FROM account JOIN shopping_cart ON account.id = shopping_cart.owner_id GROUP BY account.id"
		require.Equal(t, expected, sql)
	})

	t.Run("join with multiple group by columns", func(t *testing.T) {
		sql, _ := Select(Account.Id, Account.Type, Count().As("cnt")).From(Account).
			Join(ShoppingCart).On(Account.Id.Eq(ShoppingCart.OwnerId)).
			GroupBy(Account.Id, Account.Type).
			SQL()

		expected := "SELECT account.id, account.type, COUNT(1) AS cnt FROM account JOIN shopping_cart ON account.id = shopping_cart.owner_id GROUP BY account.id, account.type"
		require.Equal(t, expected, sql)
	})
}

func TestBuilderWithJoin_OrderBy(t *testing.T) {
	t.Run("join with order by", func(t *testing.T) {
		sql, _ := Select(Account.Id).From(Account).
			Join(ShoppingCart).On(Account.Id.Eq(ShoppingCart.OwnerId)).
			OrderBy(Account.Id.Asc()).
			SQL()

		expected := "SELECT account.id FROM account JOIN shopping_cart ON account.id = shopping_cart.owner_id ORDER BY account.id ASC"
		require.Equal(t, expected, sql)
	})

	t.Run("join with multiple order by", func(t *testing.T) {
		sql, _ := Select(Account.Id).From(Account).
			Join(ShoppingCart).On(Account.Id.Eq(ShoppingCart.OwnerId)).
			OrderBy(Account.Id.Asc(), ShoppingCart.Id.Desc()).
			SQL()

		expected := "SELECT account.id FROM account JOIN shopping_cart ON account.id = shopping_cart.owner_id ORDER BY account.id ASC, shopping_cart.id DESC"
		require.Equal(t, expected, sql)
	})
}

func TestBuilderWithJoin_AsSubQuery(t *testing.T) {
	t.Run("join as named subquery", func(t *testing.T) {
		subquery := Select(Account.Id).From(Account).
			Join(ShoppingCart).On(Account.Id.Eq(ShoppingCart.OwnerId)).
			AsNamedSubQuery("joined_data")

		sql, _ := SelectAll().From(subquery).SQL()
		expected := "SELECT * FROM (SELECT account.id FROM account JOIN shopping_cart ON account.id = shopping_cart.owner_id) AS joined_data"
		require.Equal(t, expected, sql)
	})

	t.Run("join as unnamed subquery", func(t *testing.T) {
		subquery := Select(Account.Id).From(Account).
			Join(ShoppingCart).On(Account.Id.Eq(ShoppingCart.OwnerId)).
			AsSubQuery()

		sql, _ := SelectAll().From(subquery.(Table)).SQL()
		expected := "SELECT * FROM (SELECT account.id FROM account JOIN shopping_cart ON account.id = shopping_cart.owner_id)"
		require.Equal(t, expected, sql)
	})
}

func TestBuilderWithJoin_ComplexQuery(t *testing.T) {
	t.Run("multiple joins with where and group by", func(t *testing.T) {
		acc := Account.As("a")
		sc := ShoppingCart.As("sc")

		sql, _ := Select(acc.Id, Count(sc.Id).As("cart_count")).From(acc).
			Join(sc).On(acc.Id.Eq(sc.OwnerId)).
			Where(acc.Id.EqParam(1)).
			GroupBy(acc.Id).
			OrderBy(Count(sc.Id).As("cart_count").Desc()).
			SQL()

		expected := "SELECT a.id, COUNT(sc.id) AS cart_count FROM account AS a " +
			"JOIN shopping_cart AS sc ON a.id = sc.owner_id " +
			"WHERE a.id = ? " +
			"GROUP BY a.id " +
			"ORDER BY cart_count DESC"
		require.Equal(t, expected, sql)
	})

	t.Run("three table join", func(t *testing.T) {
		acc := Account.As("a")
		sc := ShoppingCart.As("sc")

		sql, _ := Select(acc.Id, sc.Id, acc.Uuid).From(acc).
			Join(sc).On(acc.Id.Eq(sc.OwnerId)).
			SQL()

		expected := "SELECT a.id, sc.id, a.uuid FROM account AS a JOIN shopping_cart AS sc ON a.id = sc.owner_id"
		require.Equal(t, expected, sql)
	})
}

func TestBuilderWithJoin_SqlWithParams(t *testing.T) {
	t.Run("sql with params", func(t *testing.T) {
		sql, paramsMap := Select(Account.Id).From(Account).
			Join(ShoppingCart).On(Account.Id.Eq(ShoppingCart.OwnerId)).
			Where(Account.Id.EqParam(1)).
			SqlWithParams(ParamsMap{})

		expected := "SELECT account.id FROM account JOIN shopping_cart ON account.id = shopping_cart.owner_id WHERE account.id = ?"
		require.Equal(t, expected, sql)
		require.Equal(t, 1, len(paramsMap))
	})

	t.Run("sql with multiple params", func(t *testing.T) {
		sql, paramsMap := Select(Account.Id).From(Account).
			Join(ShoppingCart).On(Account.Id.Eq(ShoppingCart.OwnerId)).
			Where(Account.Id.EqParam(1).And(ShoppingCart.Id.EqParam(2))).
			SqlWithParams(ParamsMap{})

		expected := "SELECT account.id FROM account JOIN shopping_cart ON account.id = shopping_cart.owner_id WHERE account.id = ? AND shopping_cart.id = ?"
		require.Equal(t, expected, sql)
		require.Equal(t, 2, len(paramsMap))
	})
}

func TestBuilderWithJoin_Joins(t *testing.T) {
	t.Run("joins with single join item (inner join)", func(t *testing.T) {
		joinItem1 := &JoinItem{
			Target:      ShoppingCart,
			OnCondition: Account.Id.Eq(ShoppingCart.OwnerId),
		}

		sql, _ := Select(Account.Id).From(Account).
			Joins(joinItem1).
			SQL()

		expected := "SELECT account.id FROM account JOIN shopping_cart ON account.id = shopping_cart.owner_id"
		require.Equal(t, expected, sql)
	})

	t.Run("joins with multiple join items", func(t *testing.T) {
		acc := Account.As("a")
		sc := ShoppingCart.As("sc")
		sc2 := ShoppingCart.As("sc2")

		joinItem1 := &JoinItem{
			Target:      sc,
			OnCondition: acc.Id.Eq(sc.OwnerId),
		}
		joinItem2 := &JoinItem{
			Target:      sc2,
			OnCondition: acc.Id.Eq(sc2.OwnerId),
		}

		sql, _ := Select(acc.Id).From(acc).
			Joins(joinItem1, joinItem2).
			SQL()

		expected := "SELECT a.id FROM account AS a " +
			"JOIN shopping_cart AS sc ON a.id = sc.owner_id " +
			"JOIN shopping_cart AS sc2 ON a.id = sc2.owner_id"
		require.Equal(t, expected, sql)
	})

	t.Run("joins with single item then chain left join", func(t *testing.T) {
		joinItem1 := &JoinItem{
			Target:      ShoppingCart,
			OnCondition: Account.Id.Eq(ShoppingCart.OwnerId),
		}

		sql, _ := Select(Account.Id).From(Account).
			Joins(joinItem1).
			LeftJoin(ShoppingCart.As("sc2")).On(Account.Id.Eq(ShoppingCart.As("sc2").OwnerId)).
			SQL()

		expected := "SELECT account.id FROM account JOIN shopping_cart ON account.id = shopping_cart.owner_id LEFT JOIN shopping_cart AS sc2 ON account.id = sc2.owner_id"
		require.Equal(t, expected, sql)
	})

	t.Run("joins with multiple items and additional joins", func(t *testing.T) {
		sc1 := ShoppingCart.As("sc1")
		sc2 := ShoppingCart.As("sc2")
		sc3 := ShoppingCart.As("sc3")

		joinItem1 := &JoinItem{
			Target:      sc1,
			OnCondition: Account.Id.Eq(sc1.OwnerId),
		}
		joinItem2 := &JoinItem{
			Target:      sc2,
			OnCondition: Account.Id.Eq(sc2.OwnerId),
		}

		sql, _ := Select(Account.Id).From(Account).
			Joins(joinItem1, joinItem2).
			RightJoin(sc3).On(Account.Id.Eq(sc3.OwnerId)).
			SQL()

		expected := "SELECT account.id FROM account " +
			"JOIN shopping_cart AS sc1 ON account.id = sc1.owner_id " +
			"JOIN shopping_cart AS sc2 ON account.id = sc2.owner_id " +
			"RIGHT JOIN shopping_cart AS sc3 ON account.id = sc3.owner_id"
		require.Equal(t, expected, sql)
	})
}
