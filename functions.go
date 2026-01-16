package tomasql

import (
	"fmt"
	"strings"
)

func Count(col ...ParametricSql) *FuncCol[int] {
	if len(col) > 1 {
		panic("Count() accepts at most 1 column")
	}
	if len(col) == 0 {
		return newFuncCol[int]("COUNT", NewCol[int]("1", nil))
	}
	return newFuncCol[int]("COUNT", col[0])
}

func CountDistinct(col ParametricSql, otherCols ...ParametricSql) *FuncCol[int] {
	allCols := append([]ParametricSql{col}, otherCols...)
	colsPart := newMultiParametricSql(", ", allCols...)
	return newFuncCol[int]("COUNT", newMultiParametricSql("", []ParametricSql{NewFixedCol("DISTINCT ", nil), colsPart}...))
}

func Exists(subQuery ParametricSql) *FuncCol[bool] {
	return newFuncCol[bool]("EXISTS", subQuery)
}

func Sum[T any](col ParametricSql) *FuncCol[T] {
	return newFuncCol[T]("SUM", col)
}

func Avg[T any](col ParametricSql) *FuncCol[T] {
	return newFuncCol[T]("AVG", col)
}

func Min[T any](col ParametricSql) *FuncCol[T] {
	return newFuncCol[T]("MIN", col)
}

func Max[T any](col ParametricSql) *FuncCol[T] {
	return newFuncCol[T]("MAX", col)
}

func Upper(col ParametricSql) *FuncCol[string] {
	return newFuncCol[string]("UPPER", col)
}

func Lower(col ParametricSql) *FuncCol[string] {
	return newFuncCol[string]("LOWER", col)
}

func Length(col ParametricSql) *FuncCol[int] {
	return newFuncCol[int]("LENGTH", col)
}

func Coalesce[T any](col1 ParametricSql, other ...ParametricSql) *FuncCol[T] {
	return newFuncCol[T]("COALESCE", newMultiParametricSql(", ", append([]ParametricSql{col1}, other...)...))
}

func Round(col ParametricSql, decimals int) *FuncCol[float64] {
	return newFuncCol[float64]("ROUND", newMultiParametricSql(", ",
		[]ParametricSql{col, NewFixedCol(decimals, nil)}...))
}

func Abs[T any](col ParametricSql) *FuncCol[T] {
	return newFuncCol[T]("ABS", col)
}

func Trim(col ParametricSql) *FuncCol[string] {
	return newFuncCol[string]("TRIM", col)
}

type MultiParametricSql struct {
	separator string
	sqlables  []ParametricSql
}

func (m *MultiParametricSql) SqlWithParams(paramsMap ParamsMap, ctx RenderContext) (string, ParamsMap) {
	switch ctx {
	case DefinitionContext:
		var sqls []string
		for _, pSql := range m.sqlables {
			sql, pm := pSql.SqlWithParams(paramsMap, ctx)
			sqls = append(sqls, sql)
			paramsMap = pm
		}
		return strings.Join(sqls, m.separator), paramsMap
	case ReferenceContext:
		var sqls []string
		for _, pSql := range m.sqlables {
			sql, pm := pSql.SqlWithParams(paramsMap, ctx)
			sqls = append(sqls, sql)
			paramsMap = pm
		}
		return strings.Join(sqls, m.separator), paramsMap
	case OrderByContext:
		var sqls []string
		for _, pSql := range m.sqlables {
			sql, pm := pSql.SqlWithParams(paramsMap, ctx)
			sqls = append(sqls, sql)
			paramsMap = pm
		}
		return strings.Join(sqls, m.separator), paramsMap
	default:
		panic(fmt.Sprintf("MultiParametricSql.SqlWithParams: unexpected RenderContext %s", ctx))
	}
}

var _ ParametricSql = &MultiParametricSql{}

func newMultiParametricSql(separator string, sqlables ...ParametricSql) ParametricSql {
	return &MultiParametricSql{
		separator: separator,
		sqlables:  sqlables,
	}
}

type FuncColumn interface {
	As(string) FuncColumn
	Alias() *string
	Asc() SortColumn
	Desc() SortColumn

	ParametricSql
	Comparable
	SetComparable
}

type FuncCol[T any] struct {
	alias    *string // alias for the function column
	funcName string  // name of the function, e.g. "COUNT", "SUM", etc.
	inner    ParametricSql
	ComparableParam[T]
}

var _ FuncColumn = &FuncCol[any]{}

func (f *FuncCol[T]) As(s string) FuncColumn {
	f.alias = &s
	return f
}

func (f *FuncCol[T]) Alias() *string {
	return f.alias
}

func (f *FuncCol[T]) Asc() SortColumn {
	return &SortCol[T]{
		col:       nil,
		subQuery:  funcColRefWrapper[T]{funcCol: f},
		direction: OrderByAsc,
	}
}
func (f *FuncCol[T]) Desc() SortColumn {
	return &SortCol[T]{
		col:       nil,
		subQuery:  funcColRefWrapper[T]{funcCol: f},
		direction: OrderByDesc,
	}
}

func (f *FuncCol[T]) Eq(other ParametricSql) Condition {
	return NewBinaryCondition(f, other, comparerEq)
}

func (f *FuncCol[T]) EqParam(other T) Condition {
	return NewBinaryParamCondition(f, other, comparerEq)
}

func (f *FuncCol[T]) Gt(other ParametricSql) Condition {
	return NewBinaryCondition(f, other, comparerGt)
}

func (f *FuncCol[T]) GtParam(other T) Condition {
	return NewBinaryParamCondition(f, other, comparerGt)
}

func (f *FuncCol[T]) Ge(other ParametricSql) Condition {
	return NewBinaryCondition(f, other, comparerGe)
}

func (f *FuncCol[T]) GeParam(other T) Condition {
	return NewBinaryParamCondition(f, other, comparerGe)
}

func (f *FuncCol[T]) Lt(other ParametricSql) Condition {
	return NewBinaryCondition(f, other, comparerLt)
}

func (f *FuncCol[T]) LtParam(other T) Condition {
	return NewBinaryParamCondition(f, other, comparerLt)
}

func (f *FuncCol[T]) Le(other ParametricSql) Condition {
	return NewBinaryCondition(f, other, comparerLe)
}

func (f *FuncCol[T]) LeParam(other T) Condition {
	return NewBinaryParamCondition(f, other, comparerLe)
}

func (f *FuncCol[T]) Like(other ParametricSql) Condition {
	return NewBinaryCondition(f, other, comparerLike)
}

func (f *FuncCol[T]) LikeParam(pattern string) Condition {
	return NewBinaryParamCondition(f, pattern, comparerLike)
}

func (f *FuncCol[T]) IsNull() Condition {
	return newIsCondition(f, comparerNull)
}

func (f *FuncCol[T]) IsNotNull() Condition {
	return newIsCondition(f, comparerNotNull)
}

func (f *FuncCol[T]) In(sqlable ParametricSql) Condition {
	return newInCondition(f, sqlable)
}

func (f *FuncCol[T]) EqAny(sqlable ParametricSql) Condition {
	return newAnyCondition(f, comparerEq, sqlable)
}

func (f *FuncCol[T]) EqAll(sqlable ParametricSql) Condition {
	return newAllCondition(f, comparerEq, sqlable)
}

func (f *FuncCol[T]) GtAny(sqlable ParametricSql) Condition {
	return newAnyCondition(f, comparerGt, sqlable)
}

func (f *FuncCol[T]) GtAll(sqlable ParametricSql) Condition {
	return newAllCondition(f, comparerGt, sqlable)
}

func (f *FuncCol[T]) GeAny(sqlable ParametricSql) Condition {
	return newAnyCondition(f, comparerGe, sqlable)
}

func (f *FuncCol[T]) GeAll(sqlable ParametricSql) Condition {
	return newAllCondition(f, comparerGe, sqlable)
}

func (f *FuncCol[T]) LtAny(sqlable ParametricSql) Condition {
	return newAnyCondition(f, comparerLt, sqlable)
}

func (f *FuncCol[T]) LtAll(sqlable ParametricSql) Condition {
	return newAllCondition(f, comparerLt, sqlable)
}

func (f *FuncCol[T]) LeAny(sqlable ParametricSql) Condition {
	return newAnyCondition(f, comparerLe, sqlable)
}

func (f *FuncCol[T]) LeAll(sqlable ParametricSql) Condition {
	return newAllCondition(f, comparerLe, sqlable)
}

var _ FuncColumn = &FuncCol[any]{}

func newFuncCol[T any](funcName string, inner ParametricSql) *FuncCol[T] {
	return &FuncCol[T]{
		funcName: funcName,
		inner:    inner,
	}
}

func (f *FuncCol[T]) SqlWithParams(paramsMap ParamsMap, ctx RenderContext) (string, ParamsMap) {
	switch ctx {
	case DefinitionContext:
		sql := f.funcName + "("
		var innerSql string
		innerSql, paramsMap = f.inner.SqlWithParams(paramsMap, ctx)
		sql += innerSql + ")"
		// Only include alias in SELECT context
		if f.Alias() != nil {
			sql += " AS " + *f.Alias()
		}
		return sql, paramsMap
	case ReferenceContext:
		sql := f.funcName + "("
		var innerSql string
		innerSql, paramsMap = f.inner.SqlWithParams(paramsMap, ctx)
		sql += innerSql + ")"
		return sql, paramsMap
	case OrderByContext:
		// In ORDER BY context, if there's an alias, return just the alias
		if f.Alias() != nil {
			return *f.Alias(), paramsMap
		}
		// Otherwise return the full function expression
		sql := f.funcName + "("
		var innerSql string
		innerSql, paramsMap = f.inner.SqlWithParams(paramsMap, ctx)
		sql += innerSql + ")"
		return sql, paramsMap
	default:
		panic(fmt.Sprintf("FuncCol.SqlWithParams: unexpected RenderContext %s", ctx))
	}
}

// funcColRefWrapper renders a function column reference (just the alias if present, or the full function expression)
type funcColRefWrapper[T any] struct {
	funcCol *FuncCol[T]
}

func (fcrw funcColRefWrapper[T]) SqlWithParams(paramsMap ParamsMap, ctx RenderContext) (string, ParamsMap) {
	switch ctx {
	case DefinitionContext:
		if fcrw.funcCol.Alias() != nil {
			return *fcrw.funcCol.Alias(), paramsMap
		}
		// If no alias, render the full function expression without " AS ..."
		sql := fcrw.funcCol.funcName + "("
		var innerSql string
		innerSql, paramsMap = fcrw.funcCol.inner.SqlWithParams(paramsMap, OrderByContext)
		sql += innerSql + ")"
		return sql, paramsMap
	case ReferenceContext:
		if fcrw.funcCol.Alias() != nil {
			return *fcrw.funcCol.Alias(), paramsMap
		}
		// If no alias, render the full function expression without " AS ..."
		sql := fcrw.funcCol.funcName + "("
		var innerSql string
		innerSql, paramsMap = fcrw.funcCol.inner.SqlWithParams(paramsMap, OrderByContext)
		sql += innerSql + ")"
		return sql, paramsMap
	case OrderByContext:
		if fcrw.funcCol.Alias() != nil {
			return *fcrw.funcCol.Alias(), paramsMap
		}
		// If no alias, render the full function expression without " AS ..."
		sql := fcrw.funcCol.funcName + "("
		var innerSql string
		innerSql, paramsMap = fcrw.funcCol.inner.SqlWithParams(paramsMap, OrderByContext)
		sql += innerSql + ")"
		return sql, paramsMap
	default:
		panic(fmt.Sprintf("funcColRefWrapper.SqlWithParams: unexpected RenderContext %s", ctx))
	}
}
