package tomasql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// simpleTable is a minimal Table implementation for testing
type simpleTable struct {
	name  string
	alias *string
}

func (s *simpleTable) TableName() string {
	return s.name
}

func (s *simpleTable) Alias() *string {
	return s.alias
}

func (s *simpleTable) SqlWithParams(params ParamsMap, _ RenderContext) (string, ParamsMap) {
	if s.alias != nil {
		return s.name + " AS " + *s.alias, params
	}
	return s.name, params
}

// TestCol_BinaryComparisons tests all binary comparison operators
func TestCol_BinaryComparisons(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		testFn   func(col1, col2 *Col[int]) Condition
	}{
		{"Eq", "=", func(col1, col2 *Col[int]) Condition { return col1.Eq(col2) }},
		{"Neq", "<>", func(col1, col2 *Col[int]) Condition { return col1.Neq(col2) }},
		{"Gt", ">", func(col1, col2 *Col[int]) Condition { return col1.Gt(col2) }},
		{"Ge", ">=", func(col1, col2 *Col[int]) Condition { return col1.Ge(col2) }},
		{"Lt", "<", func(col1, col2 *Col[int]) Condition { return col1.Lt(col2) }},
		{"Le", "<=", func(col1, col2 *Col[int]) Condition { return col1.Le(col2) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			col1 := NewCol[int]("col1", nil)
			col2 := NewCol[int]("col2", nil)

			cond := tt.testFn(col1, col2)

			sql := cond.SQL(ParamsMap{})
			require.Equal(t, "col1 "+tt.operator+" col2", sql)

			columns := cond.Columns()
			require.Len(t, columns, 2)
			require.Equal(t, "col1", columns[0].Name())
			require.Equal(t, "col2", columns[1].Name())
		})
	}
}

// TestCol_Like tests the LIKE operator
func TestCol_Like(t *testing.T) {
	col1 := NewCol[string]("col1", nil)
	col2 := NewCol[string]("col2", nil)

	cond := col1.Like(col2)

	sql := cond.SQL(ParamsMap{})
	require.Equal(t, "col1 LIKE col2", sql)
}

// TestCol_BinaryComparisonsParam tests all binary comparison operators with parameters
func TestCol_BinaryComparisonsParam(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		testFn   func(col *Col[int]) Condition
	}{
		{"EqParam", "=", func(col *Col[int]) Condition { return col.EqParam(42) }},
		{"NeqParam", "<>", func(col *Col[int]) Condition { return col.NeqParam(42) }},
		{"GtParam", ">", func(col *Col[int]) Condition { return col.GtParam(42) }},
		{"GeParam", ">=", func(col *Col[int]) Condition { return col.GeParam(42) }},
		{"LtParam", "<", func(col *Col[int]) Condition { return col.LtParam(42) }},
		{"LeParam", "<=", func(col *Col[int]) Condition { return col.LeParam(42) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			col1 := NewCol[int]("col1", nil)

			cond := tt.testFn(col1)

			params := ParamsMap{}
			sql := cond.SQL(params)
			require.Equal(t, "col1 "+tt.operator+" "+GetDialect().Placeholder(1), sql)
			require.Equal(t, ParamsMap{42: 1}, params)

			columns := cond.Columns()
			require.Len(t, columns, 1)
			require.Equal(t, "col1", columns[0].Name())
		})
	}
}

// TestCol_LikeParam tests the LIKE operator with parameter
func TestCol_LikeParam(t *testing.T) {
	col1 := NewCol[string]("col1", nil)

	cond := col1.LikeParam("%test%")

	params := ParamsMap{}
	sql := cond.SQL(params)
	require.Equal(t, "col1 LIKE "+GetDialect().Placeholder(1), sql)
	require.Equal(t, ParamsMap{"%test%": 1}, params)
}

func TestCol_In(t *testing.T) {
	table := &simpleTable{name: "table1"}
	col1 := NewCol[int]("col1", table)
	col2 := NewCol[int]("col2", table)
	subquery := Select(col2).From(table).AsSubQuery()

	cond := col1.In(subquery)

	params := ParamsMap{}
	sql := cond.SQL(params)
	require.Equal(t, "table1.col1 IN (SELECT table1.col2 FROM table1)", sql)
}

// TestCol_QuantifiedComparisons tests all quantified comparison operators (ANY/ALL)
func TestCol_QuantifiedComparisons(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		testFn   func(col1 *Col[int], subquery SQLable) Condition
	}{
		{"EqAny", "= ANY", func(col1 *Col[int], sq SQLable) Condition { return col1.EqAny(sq) }},
		{"EqAll", "= ALL", func(col1 *Col[int], sq SQLable) Condition { return col1.EqAll(sq) }},
		{"GtAny", "> ANY", func(col1 *Col[int], sq SQLable) Condition { return col1.GtAny(sq) }},
		{"GtAll", "> ALL", func(col1 *Col[int], sq SQLable) Condition { return col1.GtAll(sq) }},
		{"GeAny", ">= ANY", func(col1 *Col[int], sq SQLable) Condition { return col1.GeAny(sq) }},
		{"GeAll", ">= ALL", func(col1 *Col[int], sq SQLable) Condition { return col1.GeAll(sq) }},
		{"LtAny", "< ANY", func(col1 *Col[int], sq SQLable) Condition { return col1.LtAny(sq) }},
		{"LtAll", "< ALL", func(col1 *Col[int], sq SQLable) Condition { return col1.LtAll(sq) }},
		{"LeAny", "<= ANY", func(col1 *Col[int], sq SQLable) Condition { return col1.LeAny(sq) }},
		{"LeAll", "<= ALL", func(col1 *Col[int], sq SQLable) Condition { return col1.LeAll(sq) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := &simpleTable{name: "table1"}
			col1 := NewCol[int]("col1", table)
			col2 := NewCol[int]("col2", table)
			subquery := Select(col2).From(table).AsSubQuery()

			cond := tt.testFn(col1, subquery)

			params := ParamsMap{}
			sql := cond.SQL(params)
			require.Equal(t, "table1.col1 "+tt.operator+"(SELECT table1.col2 FROM table1)", sql)

			columns := cond.Columns()
			require.Len(t, columns, 1)
			require.Equal(t, "col1", columns[0].Name())
		})
	}
}

func TestCol_IsNull(t *testing.T) {
	col1 := NewCol[int]("col1", nil)

	cond := col1.IsNull()

	sql := cond.SQL(ParamsMap{})
	require.Equal(t, "col1 IS NULL", sql)
}

func TestCol_IsNotNull(t *testing.T) {
	col1 := NewCol[int]("col1", nil)

	cond := col1.IsNotNull()

	sql := cond.SQL(ParamsMap{})
	require.Equal(t, "col1 IS NOT NULL", sql)
}

// TestCol_Comparisons_WithDifferentTypes tests comparisons with different column types
func TestCol_Comparisons_WithDifferentTypes(t *testing.T) {
	tests := []struct {
		name        string
		expectedSQL string
		expectedVal interface{}
		testFn      func() Condition
	}{
		{
			name:        "string column with EqParam",
			expectedSQL: "name = " + GetDialect().Placeholder(1),
			expectedVal: "John",
			testFn: func() Condition {
				col := NewCol[string]("name", nil)
				return col.EqParam("John")
			},
		},
		{
			name:        "float64 column with GtParam",
			expectedSQL: "price > " + GetDialect().Placeholder(1),
			expectedVal: 99.99,
			testFn: func() Condition {
				col := NewCol[float64]("price", nil)
				return col.GtParam(99.99)
			},
		},
		{
			name:        "bool column with EqParam",
			expectedSQL: "active = " + GetDialect().Placeholder(1),
			expectedVal: true,
			testFn: func() Condition {
				col := NewCol[bool]("active", nil)
				return col.EqParam(true)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cond := tt.testFn()

			params := ParamsMap{}
			sql := cond.SQL(params)

			require.Equal(t, tt.expectedSQL, sql)
			require.Equal(t, ParamsMap{tt.expectedVal: 1}, params)
		})
	}
}

// Test parameter handling with multiple conditions
func TestCol_Comparisons_MultipleParams(t *testing.T) {
	col1 := NewCol[int]("col1", nil)
	col2 := NewCol[int]("col2", nil)

	cond1 := col1.EqParam(10)
	cond2 := col2.GtParam(20)
	combined := cond1.And(cond2)

	params := ParamsMap{}
	sql := combined.SQL(params)

	require.Equal(t, "col1 = "+GetDialect().Placeholder(1)+" AND col2 > "+GetDialect().Placeholder(2), sql)
	require.Equal(t, ParamsMap{10: 1, 20: 2}, params)
}

// Test that conditions can be combined with And/Or
func TestCol_Comparisons_CombinedConditions(t *testing.T) {
	col1 := NewCol[int]("col1", nil)
	col2 := NewCol[int]("col2", nil)

	t.Run("AND combination", func(t *testing.T) {
		cond := col1.Eq(col2).And(col1.Gt(col2))
		sql := cond.SQL(ParamsMap{})
		require.Equal(t, "col1 = col2 AND col1 > col2", sql)
	})

	t.Run("OR combination", func(t *testing.T) {
		cond := col1.Eq(col2).Or(col1.Gt(col2))
		sql := cond.SQL(ParamsMap{})
		require.Equal(t, "col1 = col2 OR col1 > col2", sql)
	})
}

// Test ANY/ALL Columns() method
func TestCol_AnyAll_Columns(t *testing.T) {
	table := &simpleTable{name: "table1"}
	col1 := NewCol[int]("col1", table)
	col2 := NewCol[int]("col2", table)
	subquery := Select(col2).From(table).AsSubQuery()

	t.Run("ANY condition columns", func(t *testing.T) {
		cond := col1.EqAny(subquery)
		columns := cond.Columns()
		require.Len(t, columns, 1)
		require.Equal(t, "col1", columns[0].Name())
	})

	t.Run("ALL condition columns", func(t *testing.T) {
		cond := col1.GtAll(subquery)
		columns := cond.Columns()
		require.Len(t, columns, 1)
		require.Equal(t, "col1", columns[0].Name())
	})
}

// Test ComparableParam interface compliance
func TestCol_ComparableParam_InterfaceCompliance(t *testing.T) {
	var _ ComparableParam[int] = Col[int]{}
	var _ ComparableParam[string] = Col[string]{}
	var _ ComparableParam[float64] = Col[float64]{}
	var _ ComparableParam[bool] = Col[bool]{}
}
