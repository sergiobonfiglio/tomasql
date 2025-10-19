package tomasql

import (
	"fmt"
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

	if len(p) == 0 {
		// If no parameters were added, return a new map to avoid returning nil
		return ParamsMap{}
	}

	return p
}

type Condition interface {
	SQL(ParamsMap) string
	And(Condition) Condition
	Or(Condition) Condition
	// Columns returns the columns used in this condition, if applicable.
	Columns() []Column
}

type comparerType string

const (
	comparerEq      = comparerType("=")
	comparerNeq     = comparerType("<>")
	comparerGt      = comparerType(">")
	comparerGe      = comparerType(">=")
	comparerLt      = comparerType("<")
	comparerLe      = comparerType("<=")
	comparerNull    = comparerType("NULL")
	comparerNotNull = comparerType("NOT NULL")
	comparerLike    = comparerType("LIKE")
)

// IdentityCond represents a default condition that always evaluates to true (1 = 1). Could be useful in cases where
// it's necessary to handle an array of conditions and at least one condition is required.
var IdentityCond = NewBinaryCondition(NewCol[string]("1", nil), NewCol[string]("1", nil), comparerEq)

type BinaryCondition struct {
	left     ParametricSql
	right    ParametricSql
	comparer comparerType // symbol to compare the columns
}

func (b *BinaryCondition) Columns() []Column {
	var cols []Column
	if leftCol, ok := b.left.(Column); ok {
		cols = append(cols, leftCol)
	}
	if rightCol, ok := b.right.(Column); ok {
		cols = append(cols, rightCol)
	}
	return cols
}

var _ Condition = &BinaryCondition{} // Ensure BinaryCondition implements Condition

func NewBinaryCondition(left, right ParametricSql, comparer comparerType) *BinaryCondition {
	return &BinaryCondition{left: left, right: right, comparer: comparer}
}

func (b *BinaryCondition) SQL(p ParamsMap) string {
	var leftSql, rightSql string
	params := p
	leftSql, params = b.left.SqlWithParams(params)
	rightSql, _ = b.right.SqlWithParams(params)
	return fmt.Sprintf("%s %s %s", leftSql, b.comparer, rightSql)
}

func (b *BinaryCondition) And(condition Condition) Condition {
	return NewConcatCondition(AndCondConnector, b, condition)
}

func (b *BinaryCondition) Or(condition Condition) Condition {
	return NewConcatCondition(OrCondConnector, b, condition)
}

type BinaryParamCondition[T any] struct {
	col      ParametricSql
	param    T
	comparer comparerType
}

func (b *BinaryParamCondition[T]) Columns() []Column {
	var cols []Column
	if col, ok := b.col.(Column); ok {
		cols = append(cols, col)
	}
	return cols
}

var _ Condition = &BinaryParamCondition[any]{} // Ensure BinaryParamCondition implements Condition

func NewBinaryParamCondition[T any](col ParametricSql, param T, comparer comparerType) *BinaryParamCondition[T] {
	return &BinaryParamCondition[T]{col: col, param: param, comparer: comparer}
}

func (b *BinaryParamCondition[T]) SQL(params ParamsMap) string {
	// If the parameter is not already in the map, add it
	if _, ok := params[b.param]; !ok {
		params[b.param] = len(params) + 1
	}

	colSql, _ := b.col.SqlWithParams(params)
	placeholder := GetDialect().Placeholder(params[b.param])
	return fmt.Sprintf("%s %s %s", colSql, b.comparer, placeholder)
}

func (b *BinaryParamCondition[T]) And(condition Condition) Condition {
	return NewConcatCondition(AndCondConnector, b, condition)
}

func (b *BinaryParamCondition[T]) Or(condition Condition) Condition {
	return NewConcatCondition(OrCondConnector, b, condition)
}

type InCondition struct {
	col     ParametricSql
	sqlable ParametricSql
}

func (i *InCondition) Columns() []Column {
	var cols []Column
	if col, ok := i.col.(Column); ok {
		cols = append(cols, col)
	}
	if sqlableCol, ok := i.sqlable.(Column); ok {
		cols = append(cols, sqlableCol)
	}
	return cols
}

var _ Condition = &InCondition{} // Ensure InArrayCondition implements Condition

func newInCondition(col, sqlable ParametricSql) *InCondition {
	return &InCondition{col: col, sqlable: sqlable}
}

func (i *InCondition) SQL(params ParamsMap) string {
	subquerySql, _ := i.sqlable.SqlWithParams(params)
	colSql, _ := i.col.SqlWithParams(params)
	sql := fmt.Sprintf("%s IN %s", colSql, subquerySql)
	return sql
}

func (i *InCondition) And(condition Condition) Condition {
	return NewConcatCondition(AndCondConnector, i, condition)
}

func (i *InCondition) Or(condition Condition) Condition {
	return NewConcatCondition(OrCondConnector, i, condition)
}

type IsCondition struct {
	col      ParametricSql
	comparer comparerType
}

func (i IsCondition) Columns() []Column {
	var cols []Column
	if col, ok := i.col.(Column); ok {
		cols = append(cols, col)
	}
	return cols
}

var _ Condition = &IsCondition{}

func newIsCondition(col ParametricSql, comparer comparerType) *IsCondition {
	return &IsCondition{col: col, comparer: comparer}
}

func (i IsCondition) SQL(params ParamsMap) string {
	colSql, _ := i.col.SqlWithParams(params)
	return fmt.Sprintf("%s IS %s", colSql, i.comparer)
}

func (i IsCondition) And(condition Condition) Condition {
	return NewConcatCondition(AndCondConnector, i, condition)
}

func (i IsCondition) Or(condition Condition) Condition {
	return NewConcatCondition(OrCondConnector, i, condition)
}

type ExistsCondition struct {
	inner ParametricSql
}

func (e *ExistsCondition) Columns() []Column {
	var cols []Column
	if innerCol, ok := e.inner.(Column); ok {
		cols = append(cols, innerCol)
	}
	return cols
}

var _ Condition = &ExistsCondition{}

func NewExistsCondition(inner ParametricSql) *ExistsCondition {
	return &ExistsCondition{inner: inner}
}

func (e *ExistsCondition) SQL(paramsMap ParamsMap) string {
	innerSql, _ := e.inner.SqlWithParams(paramsMap)
	return fmt.Sprintf("EXISTS(%s)", innerSql)
}

func (e *ExistsCondition) And(condition Condition) Condition {
	return NewConcatCondition(AndCondConnector, e, condition)
}

func (e *ExistsCondition) Or(condition Condition) Condition {
	return NewConcatCondition(OrCondConnector, e, condition)
}
