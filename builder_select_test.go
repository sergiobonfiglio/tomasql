package tomasql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSelectCols(t *testing.T) {
	t.Run("single column", func(t *testing.T) {
		sql, _ := SelectCols(Account.Id).SQL()
		expected := "SELECT account.id"
		require.Equal(t, expected, sql)
	})

	t.Run("multiple columns", func(t *testing.T) {
		sql, _ := SelectCols(Account.Id, Account.Uuid, Account.Type).SQL()
		expected := "SELECT account.id, account.uuid, account.type"
		require.Equal(t, expected, sql)
	})

	t.Run("with FROM", func(t *testing.T) {
		sql, _ := SelectCols(Account.Id, Account.Uuid).
			From(Account).
			SQL()
		expected := "SELECT account.id, account.uuid FROM account"
		require.Equal(t, expected, sql)
	})

	t.Run("with alias", func(t *testing.T) {
		acc := Account.As("a")
		sql, _ := SelectCols(acc.Id, acc.Uuid).
			From(acc).
			SQL()
		expected := "SELECT a.id, a.uuid FROM account AS a"
		require.Equal(t, expected, sql)
	})

	t.Run("with column alias", func(t *testing.T) {
		sql, _ := SelectCols(Account.Id.As("account_id"), Account.Uuid).SQL()
		expected := "SELECT account.id AS account_id, account.uuid"
		require.Equal(t, expected, sql)
	})
}

func TestSelectDistinctCols(t *testing.T) {
	t.Run("single column", func(t *testing.T) {
		sql, _ := SelectDistinctCols(Account.Type).SQL()
		expected := "SELECT DISTINCT account.type"
		require.Equal(t, expected, sql)
	})

	t.Run("multiple columns", func(t *testing.T) {
		sql, _ := SelectDistinctCols(Account.Type, Account.Uuid).SQL()
		expected := "SELECT DISTINCT account.type, account.uuid"
		require.Equal(t, expected, sql)
	})

	t.Run("with FROM", func(t *testing.T) {
		sql, _ := SelectDistinctCols(Account.Type, Account.Uuid).
			From(Account).
			SQL()
		expected := "SELECT DISTINCT account.type, account.uuid FROM account"
		require.Equal(t, expected, sql)
	})

	t.Run("with alias", func(t *testing.T) {
		acc := Account.As("a")
		sql, _ := SelectDistinctCols(acc.Type, acc.Uuid).
			From(acc).
			SQL()
		expected := "SELECT DISTINCT a.type, a.uuid FROM account AS a"
		require.Equal(t, expected, sql)
	})

	t.Run("with column alias", func(t *testing.T) {
		sql, _ := SelectDistinctCols(Account.Type.As("acc_type"), Account.Uuid).SQL()
		expected := "SELECT DISTINCT account.type AS acc_type, account.uuid"
		require.Equal(t, expected, sql)
	})

	t.Run("with WHERE clause", func(t *testing.T) {
		sql, _ := SelectDistinctCols(Account.Type).
			From(Account).
			Where(Account.Id.GtParam(1)).
			SQL()
		expected := "SELECT DISTINCT account.type FROM account WHERE account.id > " + GetDialect().Placeholder(1)
		require.Equal(t, expected, sql)
	})
}

func TestSelect(t *testing.T) {
	t.Run("with function", func(t *testing.T) {
		sql, _ := Select(Count(), Account.Id).
			From(Account).
			GroupBy(Account.Id).
			SQL()
		expected := "SELECT COUNT(1), account.id FROM account GROUP BY account.id"
		require.Equal(t, expected, sql)
	})

	t.Run("with mixed ParametricSql types", func(t *testing.T) {
		sql, _ := Select(Account.Id, Count().As("cnt")).
			From(Account).
			GroupBy(Account.Id).
			SQL()
		expected := "SELECT account.id, COUNT(1) AS cnt FROM account GROUP BY account.id"
		require.Equal(t, expected, sql)
	})
}

func TestSelectDistinct(t *testing.T) {
	t.Run("with function", func(t *testing.T) {
		sql, _ := SelectDistinct(Account.Type, Count()).
			From(Account).
			GroupBy(Account.Type).
			SQL()
		expected := "SELECT DISTINCT account.type, COUNT(1) FROM account GROUP BY account.type"
		require.Equal(t, expected, sql)
	})
}

func TestSelectAll(t *testing.T) {
	t.Run("basic select all", func(t *testing.T) {
		sql, _ := SelectAll().From(Account).SQL()
		expected := "SELECT * FROM account"
		require.Equal(t, expected, sql)
	})

	t.Run("with WHERE", func(t *testing.T) {
		sql, _ := SelectAll().
			From(Account).
			Where(Account.Id.EqParam(1)).
			SQL()
		expected := "SELECT * FROM account WHERE account.id = " + GetDialect().Placeholder(1)
		require.Equal(t, expected, sql)
	})

	t.Run("with JOIN", func(t *testing.T) {
		sql, _ := SelectAll().
			From(Account).
			Join(ShoppingCart).On(Account.Id.Eq(ShoppingCart.OwnerId)).
			SQL()
		expected := "SELECT * FROM account JOIN shopping_cart ON account.id = shopping_cart.owner_id"
		require.Equal(t, expected, sql)
	})

	t.Run("with ORDER BY and LIMIT", func(t *testing.T) {
		sql, _ := SelectAll().
			From(Account).
			OrderBy(Account.Id.Asc()).
			Limit(5).
			SQL()
		expected := "SELECT * FROM account ORDER BY account.id ASC LIMIT 5"
		require.Equal(t, expected, sql)
	})
}

func TestSelectDistinctAll(t *testing.T) {
	t.Run("basic select distinct all", func(t *testing.T) {
		sql, _ := SelectDistinctAll().From(Account).SQL()
		expected := "SELECT DISTINCT * FROM account"
		require.Equal(t, expected, sql)
	})

	t.Run("with WHERE", func(t *testing.T) {
		sql, _ := SelectDistinctAll().
			From(Account).
			Where(Account.Type.EqParam("premium")).
			SQL()
		expected := "SELECT DISTINCT * FROM account WHERE account.type = " + GetDialect().Placeholder(1)
		require.Equal(t, expected, sql)
	})

	t.Run("with JOIN", func(t *testing.T) {
		sql, _ := SelectDistinctAll().
			From(Account).
			LeftJoin(ShoppingCart).On(Account.Id.Eq(ShoppingCart.OwnerId)).
			SQL()
		expected := "SELECT DISTINCT * FROM account LEFT JOIN shopping_cart ON account.id = shopping_cart.owner_id"
		require.Equal(t, expected, sql)
	})

	t.Run("with GROUP BY", func(t *testing.T) {
		sql, _ := SelectDistinctAll().
			From(Account).
			GroupBy(Account.Type).
			SQL()
		expected := "SELECT DISTINCT * FROM account GROUP BY account.type"
		require.Equal(t, expected, sql)
	})
}
