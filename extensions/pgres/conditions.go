package pgres

import (
	"fmt"
	"strings"

	. "github.com/sergiobonfiglio/tomasql"
)

type InArrayCondition[T any] struct {
	col   ParametricSql
	array []T
}

func (i *InArrayCondition[T]) Columns() []Column {
	var cols []Column
	if col, ok := i.col.(Column); ok {
		cols = append(cols, col)
	}
	return cols
}

var _ Condition = &InArrayCondition[any]{} // Ensure InArrayCondition implements Condition

func newInArrayCondition[T any](col ParametricSql, array []T) *InArrayCondition[T] {
	return &InArrayCondition[T]{col: col, array: array}
}

func (i *InArrayCondition[T]) SQL(params ParamsMap) string {
	paramsStr := make([]string, len(i.array))
	for ix, pItem := range i.array {
		// If the parameter is not already in the map, add it
		if _, ok := params[pItem]; !ok {
			params[pItem] = len(params) + 1
		}
		order := params[pItem]
		paramsStr[ix] = GetDialect().Placeholder(order)
	}
	allParams := strings.Join(paramsStr, ", ")

	colSql, _ := i.col.SqlWithParams(params)

	return fmt.Sprintf("%s IN (%s)", colSql, allParams)
}

func (i *InArrayCondition[T]) And(condition Condition) Condition {
	return NewConcatCondition(AndCondConnector, i, condition)
}

func (i *InArrayCondition[T]) Or(condition Condition) Condition {
	return NewConcatCondition(OrCondConnector, i, condition)
}
