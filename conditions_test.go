package tomasql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCondition_Columns(t *testing.T) {
	type test struct {
		name string
		impl Condition
		want []Column
	}

	tests := []test{
		{
			name: "identity condition columns",
			impl: IdentityCond,
			want: []Column{NewCol[string]("1", nil), NewCol[string]("1", nil)},
		},
		{
			name: "binary condition columns",
			impl: newBinaryCondition(NewCol[int64]("col1", nil), NewCol[int64]("col2", nil), comparerEq),
			want: []Column{NewCol[int64]("col1", nil), NewCol[int64]("col2", nil)},
		},
		{
			name: "binary condition with param columns",
			impl: newBinaryParamCondition(NewCol[int64]("col1", nil), int64(42), comparerEq),
			want: []Column{NewCol[int64]("col1", nil)},
		},
		// {
		// 	name: "in array condition columns",
		// 	impl: newInArrayCondition(NewCol[int64]("col1", nil), []int64{1, 2, 3}),
		// 	want: []Column{NewCol[int64]("col1", nil)},
		// },
		{
			name: "in subquery condition columns",
			impl: newInCondition(NewCol[int64]("col1", nil), Select(NewCol[int64]("col2", nil))),
			want: []Column{NewCol[int64]("col1", nil)},
		},
		{
			name: "is condition columns",
			impl: newIsCondition(NewCol[int64]("col1", nil), comparerNull),
			want: []Column{NewCol[int64]("col1", nil)},
		},
		{
			name: "grouped condition columns",
			impl: Grouped(newBinaryCondition(NewCol[int64]("col1", nil), NewCol[int64]("col2", nil), comparerEq)),
			want: []Column{NewCol[int64]("col1", nil), NewCol[int64]("col2", nil)},
		},
		{
			name: "exists condition columns",
			impl: NewExistsCondition(Select(NewCol[int64]("col1", nil))),
			want: []Column{},
		},
		{
			name: "and condition columns",
			impl: newBinaryCondition(NewCol[int64]("col1", nil), NewCol[int64]("col2", nil), comparerEq).
				And(newBinaryCondition(NewCol[int64]("col3", nil), NewCol[int64]("col4", nil), comparerEq)),
			want: []Column{
				NewCol[int64]("col1", nil),
				NewCol[int64]("col2", nil),
				NewCol[int64]("col3", nil),
				NewCol[int64]("col4", nil),
			},
		},
		{
			name: "or condition columns",
			impl: newBinaryCondition(NewCol[int64]("col1", nil), NewCol[int64]("col2", nil), comparerEq).
				Or(newBinaryCondition(NewCol[int64]("col3", nil), NewCol[int64]("col4", nil), comparerEq)),
			want: []Column{
				NewCol[int64]("col1", nil),
				NewCol[int64]("col2", nil),
				NewCol[int64]("col3", nil),
				NewCol[int64]("col4", nil),
			},
		},
	}

	for _, testItem := range tests {
		got := testItem.impl.Columns()
		name := testItem.impl.SQL(ParamsMap{})
		t.Run(testItem.name+"_"+name, func(tt *testing.T) {
			require.ElementsMatch(tt, testItem.want, got)
		})
	}
}
