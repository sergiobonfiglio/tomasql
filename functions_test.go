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
				sql, _ := Count().sqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "COUNT(1) AS c1",
			got: func() string {
				sql, _ := Count().As("c1").sqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "EXISTS(SELECT 1)",
			got: func() string {
				sql, _ := Exists(NewBuilder().Select(NewFixedCol(1, nil))).sqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "EXISTS(SELECT 1) AS e1",
			got: func() string {
				sql, _ := Exists(NewBuilder().Select(NewFixedCol(1, nil))).As("e1").sqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "SUM(col1)",
			got: func() string {
				sql, _ := Sum[int](NewCol[int]("col1", nil)).sqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "AVG(col2)",
			got: func() string {
				sql, _ := Avg[float64](NewCol[float64]("col2", nil)).sqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "MIN(col3)",
			got: func() string {
				sql, _ := Min[int](NewCol[int]("col3", nil)).sqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "MAX(col4)",
			got: func() string {
				sql, _ := Max[int](NewCol[int]("col4", nil)).sqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "UPPER(col5)",
			got: func() string {
				sql, _ := Upper(NewCol[string]("col5", nil)).sqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "LOWER(col6)",
			got: func() string {
				sql, _ := Lower(NewCol[string]("col6", nil)).sqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "LENGTH(col7)",
			got: func() string {
				sql, _ := Length(NewCol[string]("col7", nil)).sqlWithParams(ParamsMap{})
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
				).sqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "ROUND(col11, 2)",
			got: func() string {
				sql, _ := Round(NewCol[float64]("col11", nil), 2).sqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "ABS(col12)",
			got: func() string {
				sql, _ := Abs[int](NewCol[int]("col12", nil)).sqlWithParams(ParamsMap{})
				return sql
			},
		},
		{
			want: "TRIM(col13)",
			got: func() string {
				sql, _ := Trim(NewCol[string]("col13", nil)).sqlWithParams(ParamsMap{})
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
