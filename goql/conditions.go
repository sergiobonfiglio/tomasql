package goql

import (
	"fmt"
	"strings"
)

type ParamsMap map[any]int

// ToSlice returns a slice of parameters' values respecting their order as placeholders
func (p ParamsMap) ToSlice() []any {
	out := make([]any, len(p))
	for pVal, pOrder := range p {
		out[pOrder-1] = pVal
	}
	return out
}

func (p ParamsMap) AddAll(toAdd ParamsMap) ParamsMap {

	for val := range toAdd {
		if _, ok := p[val]; !ok {
			p[val] = len(p) + 1
		}
	}
	return p
}

type Condition interface {
	SQL(ParamsMap) string
	And(Condition) Condition
	Or(Condition) Condition
}

type comparerType string

const (
	comparerEq = comparerType("=")
	comparerGt = comparerType(">")
	comparerGe = comparerType(">=")
	comparerLt = comparerType("<")
	comparerLe = comparerType("<=")
)

type BinaryCondition struct {
	left     Column
	right    Column
	comparer comparerType // symbol to compare the columns
}

var _ Condition = &BinaryCondition{} // Ensure BinaryCondition implements Condition

func newBinaryCondition(left Column, right Column, comparer comparerType) *BinaryCondition {
	return &BinaryCondition{left: left, right: right, comparer: comparer}
}
func (b *BinaryCondition) SQL(_ ParamsMap) string {
	return fmt.Sprintf("%s %s %s", b.left.getRef(), b.comparer, b.right.getRef())
}

func (b *BinaryCondition) And(condition Condition) Condition {
	return newConcatCondition(AndCondConnector, b, condition)
}

func (b *BinaryCondition) Or(condition Condition) Condition {
	return newConcatCondition(OrCondConnector, b, condition)
}

type BinaryParamCondition[T any] struct {
	col      Column
	param    T
	comparer comparerType
}

var _ Condition = &BinaryParamCondition[any]{} // Ensure BinaryParamCondition implements Condition

func newBinaryParamCondition[T any](col Column, param T, comparer comparerType) *BinaryParamCondition[T] {
	return &BinaryParamCondition[T]{col: col, param: param, comparer: comparer}
}

func (b *BinaryParamCondition[T]) SQL(params ParamsMap) string {

	// If the parameter is not already in the map, add it
	if _, ok := params[b.param]; !ok {
		params[b.param] = len(params) + 1
	}

	return fmt.Sprintf("%s %s $%d", b.col.getRef(), b.comparer, params[b.param])
}

func (b *BinaryParamCondition[T]) And(condition Condition) Condition {
	return newConcatCondition(AndCondConnector, b, condition)
}

func (b *BinaryParamCondition[T]) Or(condition Condition) Condition {
	return newConcatCondition(OrCondConnector, b, condition)
}

type InArrayCondition[T any] struct {
	col   Column
	array []T
}

var _ Condition = &InArrayCondition[any]{} // Ensure InArrayCondition implements Condition

func newInArrayCondition[T any](col Column, array []T) *InArrayCondition[T] {
	return &InArrayCondition[T]{col: col, array: array}
}

func (i *InArrayCondition[T]) SQL(params ParamsMap) string {

	paramsStr := make([]string, len(i.array))
	for ix, pItem := range i.array {
		// If the parameter is not already in the map, add it
		if _, ok := params[pItem]; !ok {
			params[pItem] = len(params) + 1
		}
		order, _ := params[pItem]
		paramsStr[ix] = fmt.Sprintf("$%d", order)
	}
	allParams := strings.Join(paramsStr, ", ")
	return fmt.Sprintf("%s IN (%s)", i.col.getRef(), allParams)
}

func (i *InArrayCondition[T]) And(condition Condition) Condition {
	return newConcatCondition(AndCondConnector, i, condition)
}

func (i *InArrayCondition[T]) Or(condition Condition) Condition {
	return newConcatCondition(OrCondConnector, i, condition)
}

type InCondition struct {
	col     Column
	sqlable ParametricSql
}

var _ Condition = &InCondition{} // Ensure InArrayCondition implements Condition

func newInCondition(col Column, sqlable ParametricSql) *InCondition {
	return &InCondition{col: col, sqlable: sqlable}
}

func (i *InCondition) SQL(params ParamsMap) string {
	subquerySql, _ := i.sqlable.sqlWithParams(params)
	sql := fmt.Sprintf("%s IN %s", i.col.getRef(), subquerySql)
	return sql
}

func (i *InCondition) And(condition Condition) Condition {
	return newConcatCondition(AndCondConnector, i, condition)
}

func (i *InCondition) Or(condition Condition) Condition {
	return newConcatCondition(OrCondConnector, i, condition)
}
