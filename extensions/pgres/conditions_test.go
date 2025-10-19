package pgres

import (
	"testing"

	"github.com/sergiobonfiglio/tomasql"
	"github.com/stretchr/testify/require"
)

func TestCondition_Columns(t *testing.T) {
	type test struct {
		name string
		impl tomasql.Condition
		want []tomasql.Column
	}

	tests := []test{
		{
			name: "in array condition columns",
			impl: newInArrayCondition(tomasql.NewCol[int64]("col1", nil), []int64{1, 2, 3}),
			want: []tomasql.Column{tomasql.NewCol[int64]("col1", nil)},
		},
	}

	for _, testItem := range tests {
		got := testItem.impl.Columns()
		name := testItem.impl.SQL(tomasql.ParamsMap{})
		t.Run(testItem.name+"_"+name, func(tt *testing.T) {
			require.ElementsMatch(tt, testItem.want, got)
		})
	}
}
