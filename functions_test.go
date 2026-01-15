package tomasql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFuncColumn(t *testing.T) {
	tests := []struct {
		got  func() string
		want string
	}{
		{
			want: "COUNT(1)",
			got: func() string {
				sql, _ := Count().SqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "COUNT(1) AS c1",
			got: func() string {
				sql, _ := Count().As("c1").SqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "COUNT(col1)",
			got: func() string {
				sql, _ := Count(NewCol[int]("col1", nil)).SqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "COUNT(col1) AS c2",
			got: func() string {
				sql, _ := Count(NewCol[int]("col1", nil)).As("c2").SqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "COUNT(DISTINCT col1)",
			got: func() string {
				sql, _ := CountDistinct(NewCol[int]("col1", nil)).SqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "COUNT(DISTINCT col1) AS cd1",
			got: func() string {
				sql, _ := CountDistinct(NewCol[int]("col1", nil)).As("cd1").SqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "COUNT(DISTINCT col1, col2)",
			got: func() string {
				sql, _ := CountDistinct(NewCol[int]("col1", nil), NewCol[int]("col2", nil)).SqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "COUNT(DISTINCT col1, col2, col3) AS cd2",
			got: func() string {
				sql, _ := CountDistinct(NewCol[int]("col1", nil), NewCol[int]("col2", nil), NewCol[int]("col3", nil)).As("cd2").SqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "EXISTS(SELECT 1)",
			got: func() string {
				sql, _ := Exists(Select(NewFixedCol(1, nil))).SqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "EXISTS(SELECT 1) AS e1",
			got: func() string {
				sql, _ := Exists(Select(NewFixedCol(1, nil))).As("e1").SqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "SUM(col1)",
			got: func() string {
				sql, _ := Sum[int](NewCol[int]("col1", nil)).SqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "AVG(col2)",
			got: func() string {
				sql, _ := Avg[float64](NewCol[float64]("col2", nil)).SqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "MIN(col3)",
			got: func() string {
				sql, _ := Min[int](NewCol[int]("col3", nil)).SqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "MAX(col4)",
			got: func() string {
				sql, _ := Max[int](NewCol[int]("col4", nil)).SqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "UPPER(col5)",
			got: func() string {
				sql, _ := Upper(NewCol[string]("col5", nil)).SqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "LOWER(col6)",
			got: func() string {
				sql, _ := Lower(NewCol[string]("col6", nil)).SqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "LENGTH(col7)",
			got: func() string {
				sql, _ := Length(NewCol[string]("col7", nil)).SqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "COALESCE(col8, col9, col10)",
			got: func() string {
				sql, _ := Coalesce[string](
					NewCol[string]("col8", nil),
					NewCol[string]("col9", nil),
					NewCol[string]("col10", nil),
				).SqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "ROUND(col11, 2)",
			got: func() string {
				sql, _ := Round(NewCol[float64]("col11", nil), 2).SqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "ABS(col12)",
			got: func() string {
				sql, _ := Abs[int](NewCol[int]("col12", nil)).SqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "TRIM(col13)",
			got: func() string {
				sql, _ := Trim(NewCol[string]("col13", nil)).SqlWithParams(ParamsMap{})
				return sql
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			require.Equal(t, tt.want, tt.got())
		})
	}
}

func TestCountPanicOnMultipleColumns(t *testing.T) {
	require.Panics(t, func() {
		Count(NewCol[int]("col1", nil), NewCol[int]("col2", nil))
	}, "Count() should panic when passed more than 1 column")
}

// TestFuncCol_Comparisons tests comparison operators on function columns
func TestFuncCol_Comparisons(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		testFn   func(countFunc *FuncCol[int]) Condition
	}{
		{"Eq", "=", func(f *FuncCol[int]) Condition { return f.Eq(NewCol[int]("col1", nil)) }},
		{"Gt", ">", func(f *FuncCol[int]) Condition { return f.Gt(NewCol[int]("col1", nil)) }},
		{"Ge", ">=", func(f *FuncCol[int]) Condition { return f.Ge(NewCol[int]("col1", nil)) }},
		{"Lt", "<", func(f *FuncCol[int]) Condition { return f.Lt(NewCol[int]("col1", nil)) }},
		{"Le", "<=", func(f *FuncCol[int]) Condition { return f.Le(NewCol[int]("col1", nil)) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			countFunc := Count(NewCol[int]("col1", nil))
			cond := tt.testFn(countFunc)

			sql := cond.SQL(ParamsMap{})
			require.Equal(t, "COUNT(col1) "+tt.operator+" col1", sql)
		})
	}
}

// TestFuncCol_ComparisonsParam tests parameterized comparisons on function columns
func TestFuncCol_ComparisonsParam(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		testFn   func(countFunc *FuncCol[int]) Condition
	}{
		{"EqParam", "=", func(f *FuncCol[int]) Condition { return f.EqParam(10) }},
		{"GtParam", ">", func(f *FuncCol[int]) Condition { return f.GtParam(10) }},
		{"GeParam", ">=", func(f *FuncCol[int]) Condition { return f.GeParam(10) }},
		{"LtParam", "<", func(f *FuncCol[int]) Condition { return f.LtParam(10) }},
		{"LeParam", "<=", func(f *FuncCol[int]) Condition { return f.LeParam(10) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			countFunc := Count(NewCol[int]("col1", nil))
			cond := tt.testFn(countFunc)

			params := ParamsMap{}
			sql := cond.SQL(params)
			require.Equal(t, "COUNT(col1) "+tt.operator+" "+GetDialect().Placeholder(1), sql)
			require.Equal(t, ParamsMap{10: 1}, params)
		})
	}
}

// TestFuncCol_NullChecks tests NULL checks on function columns
func TestFuncCol_NullChecks(t *testing.T) {
	t.Run("IsNull", func(t *testing.T) {
		countFunc := Count(NewCol[int]("col1", nil))
		cond := countFunc.IsNull()

		sql := cond.SQL(ParamsMap{})
		require.Equal(t, "COUNT(col1) IS NULL", sql)
	})

	t.Run("IsNotNull", func(t *testing.T) {
		countFunc := Count(NewCol[int]("col1", nil))
		cond := countFunc.IsNotNull()

		sql := cond.SQL(ParamsMap{})
		require.Equal(t, "COUNT(col1) IS NOT NULL", sql)
	})
}

// TestFuncCol_In tests IN condition with function columns
func TestFuncCol_In(t *testing.T) {
	countFunc := Count(NewCol[int]("col1", nil))
	table := &simpleTable{name: "table1"}
	col2 := NewCol[int]("col2", table)
	subquery := Select(col2).From(table).AsSubQuery()

	cond := countFunc.In(subquery)

	params := ParamsMap{}
	sql := cond.SQL(params)
	require.Equal(t, "COUNT(col1) IN (SELECT table1.col2 FROM table1)", sql)
}

// TestFuncCol_Like tests LIKE on string function columns
func TestFuncCol_Like(t *testing.T) {
	t.Run("Like with column", func(t *testing.T) {
		upperFunc := Upper(NewCol[string]("col1", nil))
		cond := upperFunc.Like(NewCol[string]("col2", nil))

		sql := cond.SQL(ParamsMap{})
		require.Equal(t, "UPPER(col1) LIKE col2", sql)
	})

	t.Run("LikeParam", func(t *testing.T) {
		upperFunc := Upper(NewCol[string]("col1", nil))
		cond := upperFunc.LikeParam("%test%")

		params := ParamsMap{}
		sql := cond.SQL(params)
		require.Equal(t, "UPPER(col1) LIKE "+GetDialect().Placeholder(1), sql)
		require.Equal(t, ParamsMap{"%test%": 1}, params)
	})
}

// TestFuncCol_Alias tests alias functionality
func TestFuncCol_Alias(t *testing.T) {
	t.Run("Alias not set", func(t *testing.T) {
		countFunc := Count(NewCol[int]("col1", nil))
		require.Nil(t, countFunc.Alias())
	})

	t.Run("Alias set and retrieved", func(t *testing.T) {
		countFunc := Count(NewCol[int]("col1", nil)).As("cnt")
		require.NotNil(t, countFunc.Alias())
		require.Equal(t, "cnt", *countFunc.Alias())
	})

	t.Run("Alias in SQL", func(t *testing.T) {
		countFunc := Count(NewCol[int]("col1", nil)).As("cnt")
		sql, _ := countFunc.SqlWithParams(ParamsMap{})
		require.Equal(t, "COUNT(col1) AS cnt", sql)
	})
}

// TestFuncCol_Sorting tests sorting (Asc, Desc) on function columns
func TestFuncCol_Sorting(t *testing.T) {
	t.Run("Asc", func(t *testing.T) {
		countFunc := Count(NewCol[int]("col1", nil)).As("cnt")
		sortCol := countFunc.Asc()

		sql, _ := sortCol.SqlWithParams(ParamsMap{})
		require.Equal(t, "cnt ASC", sql)
	})

	t.Run("Desc", func(t *testing.T) {
		countFunc := Count(NewCol[int]("col1", nil)).As("cnt")
		sortCol := countFunc.Desc()

		sql, _ := sortCol.SqlWithParams(ParamsMap{})
		require.Equal(t, "cnt DESC", sql)
	})

	t.Run("Sorting without alias uses function expression", func(t *testing.T) {
		countFunc := Count(NewCol[int]("col1", nil))
		sortCol := countFunc.Asc()

		sql, _ := sortCol.SqlWithParams(ParamsMap{})
		require.Equal(t, "COUNT(col1) ASC", sql)
	})
}

// TestFuncCol_QuantifiedComparisons tests ANY/ALL operators on function columns
func TestFuncCol_QuantifiedComparisons(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		testFn   func(f *FuncCol[int], subquery SQLable) Condition
	}{
		{"EqAny", "= ANY", func(f *FuncCol[int], sq SQLable) Condition { return f.EqAny(sq) }},
		{"EqAll", "= ALL", func(f *FuncCol[int], sq SQLable) Condition { return f.EqAll(sq) }},
		{"GtAny", "> ANY", func(f *FuncCol[int], sq SQLable) Condition { return f.GtAny(sq) }},
		{"GtAll", "> ALL", func(f *FuncCol[int], sq SQLable) Condition { return f.GtAll(sq) }},
		{"GeAny", ">= ANY", func(f *FuncCol[int], sq SQLable) Condition { return f.GeAny(sq) }},
		{"GeAll", ">= ALL", func(f *FuncCol[int], sq SQLable) Condition { return f.GeAll(sq) }},
		{"LtAny", "< ANY", func(f *FuncCol[int], sq SQLable) Condition { return f.LtAny(sq) }},
		{"LtAll", "< ALL", func(f *FuncCol[int], sq SQLable) Condition { return f.LtAll(sq) }},
		{"LeAny", "<= ANY", func(f *FuncCol[int], sq SQLable) Condition { return f.LeAny(sq) }},
		{"LeAll", "<= ALL", func(f *FuncCol[int], sq SQLable) Condition { return f.LeAll(sq) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := &simpleTable{name: "table1"}
			col1 := NewCol[int]("col1", table)
			col2 := NewCol[int]("col2", table)
			countFunc := Count(col1)
			subquery := Select(col2).From(table).AsSubQuery()

			cond := tt.testFn(countFunc, subquery)

			params := ParamsMap{}
			sql := cond.SQL(params)
			require.Equal(t, "COUNT(table1.col1) "+tt.operator+"(SELECT table1.col2 FROM table1)", sql)
		})
	}
}
