package tomasql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuilderWithSelect_AsNamedSubQuery(t *testing.T) {
	t.Run("single column as named subquery", func(t *testing.T) {
		builder := Select(Account.Id).From(Account).AsNamedSubQuery("acc")

		sql, _ := SelectAll().From(builder).SQL()
		expected := "SELECT * FROM (SELECT account.id FROM account) AS acc"
		require.Equal(t, expected, sql)
	})

	t.Run("multiple columns as named subquery", func(t *testing.T) {
		builder := Select(Account.Id, Account.Uuid).From(Account).AsNamedSubQuery("user_data")

		sql, _ := SelectAll().From(builder).SQL()
		expected := "SELECT * FROM (SELECT account.id, account.uuid FROM account) AS user_data"
		require.Equal(t, expected, sql)
	})

	t.Run("distinct as named subquery", func(t *testing.T) {
		builder := SelectDistinct(Account.Id).From(Account).AsNamedSubQuery("distinct_ids")

		sql, _ := SelectAll().From(builder).SQL()
		expected := "SELECT * FROM (SELECT DISTINCT account.id FROM account) AS distinct_ids"
		require.Equal(t, expected, sql)
	})
}

func TestBuilderWithSelect_AsSubQuery(t *testing.T) {
	t.Run("single column as unnamed subquery", func(t *testing.T) {
		builder := Select(Account.Id).From(Account).AsSubQuery()

		sql, _ := SelectAll().From(builder.(Table)).SQL()
		expected := "SELECT * FROM (SELECT account.id FROM account)"
		require.Equal(t, expected, sql)
	})

	t.Run("multiple columns as unnamed subquery", func(t *testing.T) {
		builder := Select(Account.Id, Account.Uuid).From(Account).AsSubQuery()

		sql, _ := SelectAll().From(builder.(Table)).SQL()
		expected := "SELECT * FROM (SELECT account.id, account.uuid FROM account)"
		require.Equal(t, expected, sql)
	})

	t.Run("distinct as unnamed subquery", func(t *testing.T) {
		builder := SelectDistinct(Account.Uuid).From(Account).AsSubQuery()

		sql, _ := SelectAll().From(builder.(Table)).SQL()
		expected := "SELECT * FROM (SELECT DISTINCT account.uuid FROM account)"
		require.Equal(t, expected, sql)
	})
}

func TestBuilderWithSelect_From(t *testing.T) {
	t.Run("select with single column and FROM", func(t *testing.T) {
		sql, _ := Select(Account.Id).From(Account).SQL()
		expected := "SELECT account.id FROM account"
		require.Equal(t, expected, sql)
	})

	t.Run("select with multiple columns and FROM", func(t *testing.T) {
		sql, _ := Select(Account.Id, Account.Uuid, Account.Type).From(Account).SQL()
		expected := "SELECT account.id, account.uuid, account.type FROM account"
		require.Equal(t, expected, sql)
	})

	t.Run("select distinct with FROM", func(t *testing.T) {
		sql, _ := SelectDistinct(Account.Type).From(Account).SQL()
		expected := "SELECT DISTINCT account.type FROM account"
		require.Equal(t, expected, sql)
	})

	t.Run("select with aliased table", func(t *testing.T) {
		a := Account.As("a")
		sql, _ := Select(a.Id, a.Uuid).From(a).SQL()
		expected := "SELECT a.id, a.uuid FROM account AS a"
		require.Equal(t, expected, sql)
	})
}

func TestBuilderWithSelect_SqlWithParams(t *testing.T) {
	t.Run("single column without FROM", func(t *testing.T) {
		builder := Select(Account.Id)
		sql, paramsMap := builder.SqlWithParams(ParamsMap{})

		expected := "SELECT account.id"
		require.Equal(t, expected, sql)
		require.Equal(t, 0, len(paramsMap))
	})

	t.Run("single column with FROM", func(t *testing.T) {
		builder := Select(Account.Id).From(Account)
		sql, paramsMap := builder.SqlWithParams(ParamsMap{})

		expected := "SELECT account.id FROM account"
		require.Equal(t, expected, sql)
		require.Equal(t, 0, len(paramsMap))
	})

	t.Run("multiple columns without FROM", func(t *testing.T) {
		builder := Select(Account.Id, Account.Uuid)
		sql, paramsMap := builder.SqlWithParams(ParamsMap{})

		expected := "SELECT account.id, account.uuid"
		require.Equal(t, expected, sql)
		require.Equal(t, 0, len(paramsMap))
	})

	t.Run("multiple columns with FROM", func(t *testing.T) {
		builder := Select(Account.Id, Account.Uuid).From(Account)
		sql, paramsMap := builder.SqlWithParams(ParamsMap{})

		expected := "SELECT account.id, account.uuid FROM account"
		require.Equal(t, expected, sql)
		require.Equal(t, 0, len(paramsMap))
	})

	t.Run("with functions and parameters", func(t *testing.T) {
		builder := Select(Count().As("cnt")).From(Account)
		sql, _ := builder.SqlWithParams(ParamsMap{})

		expected := "SELECT COUNT(1) AS cnt FROM account"
		require.Equal(t, expected, sql)
	})

	t.Run("distinct flag without FROM", func(t *testing.T) {
		builder := SelectDistinct(Account.Type)
		sql, paramsMap := builder.SqlWithParams(ParamsMap{})

		expected := "SELECT DISTINCT account.type"
		require.Equal(t, expected, sql)
		require.Equal(t, 0, len(paramsMap))
	})

	t.Run("distinct flag with FROM", func(t *testing.T) {
		builder := SelectDistinct(Account.Type).From(Account)
		sql, paramsMap := builder.SqlWithParams(ParamsMap{})

		expected := "SELECT DISTINCT account.type FROM account"
		require.Equal(t, expected, sql)
		require.Equal(t, 0, len(paramsMap))
	})
}

func TestBuilderWithSelectAll_Basic(t *testing.T) {
	t.Run("select all", func(t *testing.T) {
		sql, _ := SelectAll().SQL()
		expected := "SELECT *"
		require.Equal(t, expected, sql)
	})

	t.Run("select all with from", func(t *testing.T) {
		sql, _ := SelectAll().From(Account).SQL()
		expected := "SELECT * FROM account"
		require.Equal(t, expected, sql)
	})

	t.Run("select distinct all", func(t *testing.T) {
		sql, _ := SelectDistinctAll().SQL()
		expected := "SELECT DISTINCT *"
		require.Equal(t, expected, sql)
	})

	t.Run("select distinct all with from", func(t *testing.T) {
		sql, _ := SelectDistinctAll().From(Account).SQL()
		expected := "SELECT DISTINCT * FROM account"
		require.Equal(t, expected, sql)
	})
}

func TestWithOptionalAlias_TableName(t *testing.T) {
	t.Run("with alias", func(t *testing.T) {
		alias := "my_alias"
		subquery := Select(Account.Id).From(Account).AsNamedSubQuery(alias)

		require.Equal(t, alias, subquery.TableName())
	})

	t.Run("without alias", func(t *testing.T) {
		subquery := Select(Account.Id).From(Account).AsSubQuery().(Table)

		require.Equal(t, "", subquery.TableName())
	})
}

func TestWithOptionalAlias_Alias(t *testing.T) {
	t.Run("returns alias pointer when set", func(t *testing.T) {
		alias := "test_alias"
		subquery := Select(Account.Id).From(Account).AsNamedSubQuery(alias)

		require.NotNil(t, subquery.Alias())
		require.Equal(t, alias, *subquery.Alias())
	})

	t.Run("returns nil alias when not set", func(t *testing.T) {
		subquery := Select(Account.Id).From(Account).AsSubQuery().(Table)

		require.Nil(t, subquery.Alias())
	})
}

func TestWithOptionalAlias_As(t *testing.T) {
	t.Run("change alias on existing subquery", func(t *testing.T) {
		subquery := Select(Account.Id).From(Account).AsSubQuery().(*withOptionalAlias)

		// Initially has no alias
		require.Nil(t, subquery.Alias())

		// Add alias via As method
		newAlias := "updated_alias"
		subquery.As(&newAlias)

		require.NotNil(t, subquery.Alias())
		require.Equal(t, newAlias, *subquery.Alias())
	})

	t.Run("As method returns self for chaining", func(t *testing.T) {
		subquery := Select(Account.Id).From(Account).AsSubQuery().(*withOptionalAlias)
		newAlias := "alias"

		result := subquery.As(&newAlias)

		require.NotNil(t, result)
	})
}

func TestWithOptionalAlias_SqlWithParams(t *testing.T) {
	t.Run("subquery with alias", func(t *testing.T) {
		subquery := Select(Account.Id, Account.Uuid).From(Account).AsNamedSubQuery("acc")

		sql, _ := SelectAll().From(subquery).SQL()
		expected := "SELECT * FROM (SELECT account.id, account.uuid FROM account) AS acc"
		require.Equal(t, expected, sql)
	})

	t.Run("subquery without alias", func(t *testing.T) {
		subquery := Select(Account.Id).From(Account).AsSubQuery()

		sql, _ := SelectAll().From(subquery.(Table)).SQL()
		expected := "SELECT * FROM (SELECT account.id FROM account)"
		require.Equal(t, expected, sql)
	})

	t.Run("nested subqueries with aliases", func(t *testing.T) {
		inner := Select(Account.Id).From(Account).AsNamedSubQuery("inner")
		outer := Select(Count().As("cnt")).From(inner).AsNamedSubQuery("outer")

		sql, _ := SelectAll().From(outer).SQL()
		expected := "SELECT * FROM (SELECT COUNT(1) AS cnt FROM (SELECT account.id FROM account) AS inner) AS outer"
		require.Equal(t, expected, sql)
	})

	t.Run("subquery with parameters", func(t *testing.T) {
		subquery := Select(Account.Id).From(Account).Where(Account.Id.EqParam(1)).AsNamedSubQuery("filtered")

		sql, params := SelectAll().From(subquery).SQL()
		expected := "SELECT * FROM (SELECT account.id FROM account WHERE account.id = ?) AS filtered"
		require.Equal(t, expected, sql)
		require.Equal(t, 1, len(params))
	})
}

func TestBuilderWithSelect_EdgeCases(t *testing.T) {
	t.Run("select with function columns", func(t *testing.T) {
		sql, _ := Select(Count().As("total"), CountDistinct(Account.Id).As("sum_id")).From(Account).SQL()
		expected := "SELECT COUNT(1) AS total, COUNT(DISTINCT account.id) AS sum_id FROM account"
		require.Equal(t, expected, sql)
	})

	t.Run("select with mixed column types", func(t *testing.T) {
		sql, _ := Select(Account.Id, Count().As("cnt"), Account.Uuid).From(Account).SQL()
		expected := "SELECT account.id, COUNT(1) AS cnt, account.uuid FROM account"
		require.Equal(t, expected, sql)
	})

	t.Run("select all from subquery with alias", func(t *testing.T) {
		subquery := Select(Account.Id, Account.Uuid).From(Account).AsNamedSubQuery("data")
		sql, _ := SelectAll().From(subquery).SQL()
		expected := "SELECT * FROM (SELECT account.id, account.uuid FROM account) AS data"
		require.Equal(t, expected, sql)
	})
}
