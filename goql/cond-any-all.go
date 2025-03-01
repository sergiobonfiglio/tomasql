package goql

import (
	"fmt"
	"github.com/lib/pq"
)

type quantifiedOperator string

const (
	anyOperator = quantifiedOperator("ANY")
	allOperator = quantifiedOperator("ALL")
)

type anyAllArrayCondition struct {
	col      Column
	comparer comparerType
	operator quantifiedOperator
	sqlable  ParametricSql
	params   ParamsMap
}

var _ Condition = &anyAllArrayCondition{}

func newAnyArrayCondition(col Column, comparer comparerType, sqlable ParametricSql) Condition {
	return &anyAllArrayCondition{
		col:      col,
		comparer: comparer,
		operator: anyOperator,
		sqlable:  sqlable,
		params:   ParamsMap{},
	}
}

func newAllArrayCondition(col Column, comparer comparerType, sqlable ParametricSql) Condition {
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
	sql, _ := a.sqlable.sqlWithParams(a.params)

	return fmt.Sprintf("%s %s %s%s", a.col.getRef(), a.comparer, a.operator, sql)
}

func (a *anyAllArrayCondition) And(condition Condition) Condition {
	return newConcatCondition(AndCondConnector, a, condition)
}

func (a *anyAllArrayCondition) Or(condition Condition) Condition {
	return newConcatCondition(OrCondConnector, a, condition)
}

type sqlableArray struct {
	array  any
	params ParamsMap
}

var _ ParametricSql = &sqlableArray{}

func SQLArray[T any](array []T) ParametricSql {
	return newSQLableArray(array)
}

func newSQLableArray(array any) ParametricSql {
	return &sqlableArray{array: pq.Array(array), params: ParamsMap{}}
}

//func (s *sqlableArray) SQL() (sql string, params []any) {
//	sql, paramsMap := s.sqlWithParams(s.params)
//	return sql, paramsMap.ToSlice()
//}

func (s *sqlableArray) sqlWithParams(params ParamsMap) (string, ParamsMap) {
	s.params = params.AddAll(s.params)
	if _, ok := params[s.array]; !ok {
		params[s.array] = len(params) + 1
	}
	return fmt.Sprintf("($%d)", params[s.array]), params
}

//func (s *sqlableArray) AsNamedSubQuery(alias string) SQLable {
//	panic("unsupported")
//}
//
//func (s *sqlableArray) AsSubQuery() SQLable {
//	panic("unsupported")
//}
