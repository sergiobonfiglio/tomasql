package tomasql

import (
	"fmt"
)

type quantifiedOperator string

const (
	anyOperator = quantifiedOperator("ANY")
	allOperator = quantifiedOperator("ALL")
)

type anyAllCondition struct {
	col      ParametricSql
	comparer comparerType
	operator quantifiedOperator
	sqlable  ParametricSql
	params   ParamsMap
}

func (a *anyAllCondition) Columns() []Column {
	var cols []Column
	if col, ok := a.col.(Column); ok {
		cols = append(cols, col)
	}
	if sqlableCol, ok := a.sqlable.(Column); ok {
		cols = append(cols, sqlableCol)
	}
	return cols
}

var _ Condition = &anyAllCondition{}

func newAnyCondition(col ParametricSql, comparer comparerType, sqlable ParametricSql) Condition {
	return &anyAllCondition{
		col:      col,
		comparer: comparer,
		operator: anyOperator,
		sqlable:  sqlable,
		params:   ParamsMap{},
	}
}

func newAllCondition(col ParametricSql, comparer comparerType, sqlable ParametricSql) Condition {
	return &anyAllCondition{
		col:      col,
		comparer: comparer,
		operator: allOperator,
		sqlable:  sqlable,
		params:   ParamsMap{},
	}
}

func (a *anyAllCondition) SQL(params ParamsMap) string {
	a.params = params.AddAll(a.params)
	sqlWithParams, _ := a.sqlable.SqlWithParams(a.params, ReferenceContext)
	colSql, _ := a.col.SqlWithParams(a.params, ReferenceContext)

	return fmt.Sprintf("%s %s %s%s", colSql, a.comparer, a.operator, sqlWithParams)
}

func (a *anyAllCondition) And(condition Condition) Condition {
	return NewConcatCondition(AndCondConnector, a, condition)
}

func (a *anyAllCondition) Or(condition Condition) Condition {
	return NewConcatCondition(OrCondConnector, a, condition)
}
