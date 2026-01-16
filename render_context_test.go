package tomasql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestColSqlWithParams_RenderContexts tests Col.SqlWithParams with different RenderContext values
func TestColSqlWithParams_RenderContexts(t *testing.T) {
	table := &simpleTable{name: "users"}
	col := NewCol[string]("username", table)
	params := ParamsMap{}

	t.Run("without alias", func(t *testing.T) {
		tests := []struct {
			name     string
			ctx      RenderContext
			expected string
		}{
			{name: "DefinitionContext", ctx: DefinitionContext, expected: "users.username"},
			{name: "ReferenceContext", ctx: ReferenceContext, expected: "users.username"},
			{name: "OrderByContext", ctx: OrderByContext, expected: "users.username"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				sql, _ := col.SqlWithParams(params, tt.ctx)
				require.Equal(t, tt.expected, sql)
			})
		}
	})

	t.Run("with alias", func(t *testing.T) {
		colWithAlias := col.As("user_name")
		tests := []struct {
			name     string
			ctx      RenderContext
			expected string
		}{
			{name: "DefinitionContext (includes AS)", ctx: DefinitionContext, expected: "users.username AS user_name"},
			{name: "ReferenceContext (no alias)", ctx: ReferenceContext, expected: "users.username"},
			{name: "OrderByContext (uses alias)", ctx: OrderByContext, expected: "user_name"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				sql, _ := colWithAlias.SqlWithParams(params, tt.ctx)
				require.Equal(t, tt.expected, sql)
			})
		}
	})

	t.Run("with table alias", func(t *testing.T) {
		tableWithAlias := &simpleTable{name: "users", alias: strPtr("u")}
		colWithTableAlias := NewCol[string]("username", tableWithAlias)
		tests := []struct {
			name     string
			ctx      RenderContext
			expected string
		}{
			{name: "DefinitionContext", ctx: DefinitionContext, expected: "u.username"},
			{name: "ReferenceContext", ctx: ReferenceContext, expected: "u.username"},
			{name: "OrderByContext", ctx: OrderByContext, expected: "u.username"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				sql, _ := colWithTableAlias.SqlWithParams(params, tt.ctx)
				require.Equal(t, tt.expected, sql)
			})
		}
	})

	t.Run("with table alias and column alias", func(t *testing.T) {
		tableWithAlias := &simpleTable{name: "users", alias: strPtr("u")}
		colWithBothAliases := NewCol[string]("username", tableWithAlias).As("user_name")
		tests := []struct {
			name     string
			ctx      RenderContext
			expected string
		}{
			{name: "DefinitionContext (includes AS)", ctx: DefinitionContext, expected: "u.username AS user_name"},
			{name: "ReferenceContext (no alias)", ctx: ReferenceContext, expected: "u.username"},
			{name: "OrderByContext (uses alias)", ctx: OrderByContext, expected: "user_name"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				sql, _ := colWithBothAliases.SqlWithParams(params, tt.ctx)
				require.Equal(t, tt.expected, sql)
			})
		}
	})

	t.Run("column without table", func(t *testing.T) {
		colNoTable := NewCol[string]("*", nil)
		sql, _ := colNoTable.SqlWithParams(params, DefinitionContext)
		require.Equal(t, "*", sql)
	})
}

// TestTableSqlWithParams_RenderContexts tests Table.SqlWithParams with different RenderContext values
func TestTableSqlWithParams_RenderContexts(t *testing.T) {
	params := ParamsMap{}

	t.Run("sqlableTable", func(t *testing.T) {
		tests := []struct {
			name     string
			table    Table
			expected string
		}{
			{
				name:     "without alias",
				table:    &simpleTable{name: "accounts"},
				expected: "accounts",
			},
			{
				name:     "with alias",
				table:    &simpleTable{name: "accounts", alias: strPtr("a")},
				expected: "accounts AS a",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				sqlTable := newSqlableTable(tt.table)
				// sqlableTable behavior is the same across all contexts
				for _, ctx := range []RenderContext{DefinitionContext, ReferenceContext, OrderByContext} {
					sql, _ := sqlTable.SqlWithParams(params, ctx)
					require.Equal(t, tt.expected, sql)
				}
			})
		}
	})

	t.Run("tableRefWrapper", func(t *testing.T) {
		tests := []struct {
			name     string
			table    Table
			expected string
		}{
			{
				name:     "without alias",
				table:    &simpleTable{name: "users"},
				expected: "users",
			},
			{
				name:     "with alias (returns only alias)",
				table:    &simpleTable{name: "users", alias: strPtr("u")},
				expected: "u",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				refWrapper := tableRefWrapper{table: tt.table}
				// tableRefWrapper behavior is the same across all contexts
				for _, ctx := range []RenderContext{DefinitionContext, ReferenceContext, OrderByContext} {
					sql, _ := refWrapper.SqlWithParams(params, ctx)
					require.Equal(t, tt.expected, sql)
				}
			})
		}
	})
}

// TestFuncColSqlWithParams_RenderContexts tests FuncCol.SqlWithParams with different RenderContext values
// Focus on the OrderByContext behavior which differs when alias is present
func TestFuncColSqlWithParams_RenderContexts(t *testing.T) {
	table := &simpleTable{name: "orders"}
	col := NewCol[int]("amount", table)
	params := ParamsMap{}

	t.Run("function with alias - OrderByContext returns only alias", func(t *testing.T) {
		funcCol := Sum[int](col).As("total_amount")
		sql, _ := funcCol.SqlWithParams(params, OrderByContext)
		require.Equal(t, "total_amount", sql)
	})

	t.Run("function without alias - all contexts return full function", func(t *testing.T) {
		funcCol := Sum[int](col)
		expected := "SUM(orders.amount)"
		for _, ctx := range []RenderContext{DefinitionContext, ReferenceContext, OrderByContext} {
			sql, _ := funcCol.SqlWithParams(params, ctx)
			require.Equal(t, expected, sql)
		}
	})

	t.Run("function with alias in DefinitionContext includes AS clause", func(t *testing.T) {
		funcCol := Sum[int](col).As("total_amount")
		sql, _ := funcCol.SqlWithParams(params, DefinitionContext)
		require.Equal(t, "SUM(orders.amount) AS total_amount", sql)
	})

	t.Run("function with alias in ReferenceContext omits alias", func(t *testing.T) {
		funcCol := Sum[int](col).As("total_amount")
		sql, _ := funcCol.SqlWithParams(params, ReferenceContext)
		require.Equal(t, "SUM(orders.amount)", sql)
	})

	t.Run("COUNT variations", func(t *testing.T) {
		tests := []struct {
			name     string
			funcCol  FuncColumn
			ctx      RenderContext
			expected string
		}{
			{
				name:     "COUNT() in DefinitionContext",
				funcCol:  Count(),
				ctx:      DefinitionContext,
				expected: "COUNT(1)",
			},
			{
				name:     "COUNT() in OrderByContext",
				funcCol:  Count(),
				ctx:      OrderByContext,
				expected: "COUNT(1)",
			},
			{
				name:     "COUNT with alias in DefinitionContext",
				funcCol:  Count().As("total"),
				ctx:      DefinitionContext,
				expected: "COUNT(1) AS total",
			},
			{
				name:     "COUNT with alias in OrderByContext",
				funcCol:  Count().As("total"),
				ctx:      OrderByContext,
				expected: "total",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				sql, _ := tt.funcCol.SqlWithParams(params, tt.ctx)
				require.Equal(t, tt.expected, sql)
			})
		}
	})

	t.Run("various function types with alias in OrderByContext", func(t *testing.T) {
		tests := []struct {
			name    string
			funcCol FuncColumn
		}{
			{name: "AVG", funcCol: Avg[float64](col).As("avg_amount")},
			{name: "MAX", funcCol: Max[int](col).As("max_amount")},
			{name: "MIN", funcCol: Min[int](col).As("min_amount")},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				sql, _ := tt.funcCol.SqlWithParams(params, OrderByContext)
				// All should return just the alias in OrderByContext
				require.True(t, len(sql) < 20, "should be short alias, not full function expression")
			})
		}
	})
}

// Helper function to create string pointers
func strPtr(s string) *string {
	return &s
}
