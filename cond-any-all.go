package tomasql

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

type quantifiedOperator string

const (
	anyOperator = quantifiedOperator("ANY")
	allOperator = quantifiedOperator("ALL")
)

type anyAllArrayCondition struct {
	col      ParametricSql
	comparer comparerType
	operator quantifiedOperator
	sqlable  ParametricSql
	params   ParamsMap
}

func (a *anyAllArrayCondition) Columns() []Column {
	var cols []Column
	if col, ok := a.col.(Column); ok {
		cols = append(cols, col)
	}
	if sqlableCol, ok := a.sqlable.(Column); ok {
		cols = append(cols, sqlableCol)
	}
	return cols
}

var _ Condition = &anyAllArrayCondition{}

func newAnyArrayCondition(col ParametricSql, comparer comparerType, sqlable ParametricSql) Condition {
	return &anyAllArrayCondition{
		col:      col,
		comparer: comparer,
		operator: anyOperator,
		sqlable:  sqlable,
		params:   ParamsMap{},
	}
}

func newAllArrayCondition(col ParametricSql, comparer comparerType, sqlable ParametricSql) Condition {
	return &anyAllArrayCondition{
		col:      col,
		comparer: comparer,
		operator: allOperator,
		sqlable:  sqlable,
		params:   ParamsMap{},
	}
}

func (a *anyAllArrayCondition) SQL(params ParamsMap) string {
	a.params = params.AddAll(a.params)
	sqlWithParams, _ := a.sqlable.sqlWithParams(a.params)
	colSql, _ := a.col.sqlWithParams(a.params)

	return fmt.Sprintf("%s %s %s%s", colSql, a.comparer, a.operator, sqlWithParams)
}

func (a *anyAllArrayCondition) And(condition Condition) Condition {
	return newConcatCondition(AndCondConnector, a, condition)
}

func (a *anyAllArrayCondition) Or(condition Condition) Condition {
	return newConcatCondition(OrCondConnector, a, condition)
}

type sqlableArray struct {
	array  sql.Scanner
	params ParamsMap
}

var _ ParametricSql = &sqlableArray{}

func SQLArray[T any](array []T) ParametricSql {
	return newSQLableArray(array)
}

func newSQLableArray(array any) ParametricSql {
	return &sqlableArray{array: pq.Array(array), params: ParamsMap{}}
}

func (s *sqlableArray) sqlWithParams(params ParamsMap) (string, ParamsMap) {
	s.params = params.AddAll(s.params)

	// Use a pointer to the slice as a key to ensure it's hashable. Note: we'll have to pass the array
	// multiple times in the query params even if it's the same array, but it should be fine for most use cases.
	arrayKey := &s.array
	if _, ok := s.params[arrayKey]; !ok {
		s.params[arrayKey] = len(s.params) + 1
	}
	return fmt.Sprintf("($%d)", s.params[arrayKey]), s.params
}
