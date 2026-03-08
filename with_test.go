package tomasql

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type testPostgresDialect struct{}

func (d *testPostgresDialect) Name() string {
	return "postgres"
}

func (d *testPostgresDialect) Placeholder(position int) string {
	return fmt.Sprintf("$%d", position)
}

func TestWith_SingleCTE(t *testing.T) {
	cte := CTE("recent", Select(Account.Id).From(Account))

	sql, params := With(cte).SelectAll().From(cte.Table()).SQL()

	require.Equal(t, "WITH recent AS (SELECT account.id FROM account) SELECT * FROM recent", sql)
	require.Empty(t, params)
}

func TestWith_MultipleCTEs_PreservesOrder(t *testing.T) {
	a := CTE("a", Select(Account.Id).From(Account))
	b := CTE("b", Select(ShoppingCart.Id).From(ShoppingCart))

	sql, _ := With(a, b).SelectAll().From(b.Table()).SQL()

	require.Equal(t, "WITH a AS (SELECT account.id FROM account), b AS (SELECT shopping_cart.id FROM shopping_cart) SELECT * FROM b", sql)
}

func TestWith_ParameterOrderingAcrossCTEAndMainQuery(t *testing.T) {
	cte := CTE("x", Select(Account.Id).From(Account).Where(Account.Id.EqParam(int64(1))))

	sql, params := With(cte).SelectAll().From(cte.Table()).Where(NewCol[int64]("id", cte.Table()).EqParam(int64(2))).SQL()

	require.Equal(t, "WITH x AS (SELECT account.id FROM account WHERE account.id = ?) SELECT * FROM x WHERE x.id = ?", sql)
	require.Equal(t, []any{int64(1), int64(2)}, params)
}

func TestWith_PostgresParameterOrdering(t *testing.T) {
	originalDialect := GetDialect()
	defer SetDialect(originalDialect)
	SetDialect(&testPostgresDialect{})

	cte := CTE("x", Select(Account.Id).From(Account).Where(Account.Id.EqParam(int64(1))))

	sql, params := With(cte).SelectAll().From(cte.Table()).Where(NewCol[int64]("id", cte.Table()).EqParam(int64(2))).SQL()

	require.Equal(t, "WITH x AS (SELECT account.id FROM account WHERE account.id = $1) SELECT * FROM x WHERE x.id = $2", sql)
	require.Equal(t, []any{int64(1), int64(2)}, params)
}

func TestWith_CTECanReferencePreviousCTE(t *testing.T) {
	a := CTE("a", Select(Account.Id).From(Account))
	b := CTE("b", SelectAll().From(a.Table()))

	sql, _ := With(a, b).SelectAll().From(b.Table()).SQL()

	require.Equal(t, "WITH a AS (SELECT account.id FROM account), b AS (SELECT * FROM a) SELECT * FROM b", sql)
}

func TestWith_AsNamedSubQueryPreservesWithClause(t *testing.T) {
	cte := CTE("recent", Select(Account.Id).From(Account))
	inner := With(cte).SelectAll().From(cte.Table()).AsNamedSubQuery("q")

	sql, _ := SelectAll().From(inner).SQL()

	require.Equal(t, "SELECT * FROM (WITH recent AS (SELECT account.id FROM account) SELECT * FROM recent) AS q", sql)
}

func TestCTE_InvalidInputs(t *testing.T) {
	require.Panics(t, func() {
		CTE("", SelectAll().From(Account))
	})

	require.Panics(t, func() {
		CTE("x", nil)
	})

	require.Panics(t, func() {
		With(nil)
	})
}
