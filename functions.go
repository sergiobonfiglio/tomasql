package tomasql

import "strings"

func Count() FuncColumn[int] {
	return newFuncCol[int]("COUNT", NewCol[int]("1", nil))
}

func Exists(subQuery ParametricSql) FuncColumn[bool] {
	return newFuncCol[bool]("EXISTS", subQuery)
}

func Sum[T any](col ParametricSql) FuncColumn[T] {
	return newFuncCol[T]("SUM", col)
}

func Avg[T any](col ParametricSql) FuncColumn[T] {
	return newFuncCol[T]("AVG", col)
}

func Min[T any](col ParametricSql) FuncColumn[T] {
	return newFuncCol[T]("MIN", col)
}

func Max[T any](col ParametricSql) FuncColumn[T] {
	return newFuncCol[T]("MAX", col)
}

func Upper(col ParametricSql) FuncColumn[string] {
	return newFuncCol[string]("UPPER", col)
}

func Lower(col ParametricSql) FuncColumn[string] {
	return newFuncCol[string]("LOWER", col)
}

func Length(col ParametricSql) FuncColumn[int] {
	return newFuncCol[int]("LENGTH", col)
}

func Coalesce[T any](col1 ParametricSql, other ...ParametricSql) FuncColumn[T] {
	return newFuncCol[T]("COALESCE", newMultiParametricSql(", ", append([]ParametricSql{col1}, other...)...))
}

func Round(col ParametricSql, decimals int) FuncColumn[float64] {
	return newFuncCol[float64]("ROUND", newMultiParametricSql(", ",
		[]ParametricSql{col, NewFixedCol(decimals, nil)}...))
}

func Abs[T any](col ParametricSql) FuncColumn[T] {
	return newFuncCol[T]("ABS", col)
}

func Trim(col ParametricSql) FuncColumn[string] {
	return newFuncCol[string]("TRIM", col)
}

type MultiParametricSql struct {
	separator string
	sqlables  []ParametricSql
}

func (m *MultiParametricSql) SqlWithParams(paramsMap ParamsMap) (string, ParamsMap) {
	var sqls []string
	for _, pSql := range m.sqlables {
		sql, pm := pSql.SqlWithParams(paramsMap)
		sqls = append(sqls, sql)
		paramsMap = pm
	}
	return strings.Join(sqls, m.separator), paramsMap
}

var _ ParametricSql = &MultiParametricSql{}

func newMultiParametricSql(separator string, sqlables ...ParametricSql) ParametricSql {
	return &MultiParametricSql{
		separator: separator,
		sqlables:  sqlables,
	}
}

type FuncColumn[T any] interface {
	As(string) ParametricSql
	Alias() *string
	ParametricSql
	Comparable
	SetComparable
	ComparableParam[T]
}

type funcCol[T any] struct {
	alias    *string // alias for the function column
	funcName string  // name of the function, e.g. "COUNT", "SUM", etc.
	inner    ParametricSql
}

func (f *funcCol[T]) As(s string) ParametricSql {
	f.alias = &s
	return f
}

func (f *funcCol[T]) Alias() *string {
	return f.alias
}

func (f *funcCol[T]) Eq(other ParametricSql) Condition {
	return newBinaryCondition(f, other, comparerEq)
}

func (f *funcCol[T]) EqParam(other T) Condition {
	return newBinaryParamCondition(f, other, comparerEq)
}

func (f *funcCol[T]) Gt(other ParametricSql) Condition {
	return newBinaryCondition(f, other, comparerGt)
}

func (f *funcCol[T]) GtParam(other T) Condition {
	return newBinaryParamCondition(f, other, comparerGt)
}

func (f *funcCol[T]) Ge(other ParametricSql) Condition {
	return newBinaryCondition(f, other, comparerGe)
}

func (f *funcCol[T]) GeParam(other T) Condition {
	return newBinaryParamCondition(f, other, comparerGe)
}

func (f *funcCol[T]) Lt(other ParametricSql) Condition {
	return newBinaryCondition(f, other, comparerLt)
}

func (f *funcCol[T]) LtParam(other T) Condition {
	return newBinaryParamCondition(f, other, comparerLt)
}

func (f *funcCol[T]) Le(other ParametricSql) Condition {
	return newBinaryCondition(f, other, comparerLe)
}

func (f *funcCol[T]) LeParam(other T) Condition {
	return newBinaryParamCondition(f, other, comparerLe)
}

func (f *funcCol[T]) Like(other ParametricSql) Condition {
	return newBinaryCondition(f, other, comparerLike)
}

func (f *funcCol[T]) LikeParam(pattern string) Condition {
	return newBinaryParamCondition(f, pattern, comparerLike)
}

func (f *funcCol[T]) ILike(other ParametricSql) Condition {
	return newBinaryCondition(f, other, comparerILike)
}

func (f *funcCol[T]) ILikeParam(pattern string) Condition {
	return newBinaryParamCondition(f, pattern, comparerILike)
}

func (f *funcCol[T]) IsNull() Condition {
	return newIsCondition(f, comparerNull)
}

func (f *funcCol[T]) IsNotNull() Condition {
	return newIsCondition(f, comparerNotNull)
}

// func (f *funcCol[T]) InArray(array []T) Condition {
// 	return newInArrayCondition(f, array)
// }

func (f *funcCol[T]) In(sqlable ParametricSql) Condition {
	return newInCondition(f, sqlable)
}

func (f *funcCol[T]) EqAny(sqlable ParametricSql) Condition {
	return newAnyArrayCondition(f, comparerEq, sqlable)
}

func (f *funcCol[T]) EqAll(sqlable ParametricSql) Condition {
	return newAllArrayCondition(f, comparerEq, sqlable)
}

func (f *funcCol[T]) GtAny(sqlable ParametricSql) Condition {
	return newAnyArrayCondition(f, comparerGt, sqlable)
}

func (f *funcCol[T]) GtAll(sqlable ParametricSql) Condition {
	return newAllArrayCondition(f, comparerGt, sqlable)
}

func (f *funcCol[T]) GeAny(sqlable ParametricSql) Condition {
	return newAnyArrayCondition(f, comparerGe, sqlable)
}

func (f *funcCol[T]) GeAll(sqlable ParametricSql) Condition {
	return newAllArrayCondition(f, comparerGe, sqlable)
}

func (f *funcCol[T]) LtAny(sqlable ParametricSql) Condition {
	return newAnyArrayCondition(f, comparerLt, sqlable)
}

func (f *funcCol[T]) LtAll(sqlable ParametricSql) Condition {
	return newAllArrayCondition(f, comparerLt, sqlable)
}

func (f *funcCol[T]) LeAny(sqlable ParametricSql) Condition {
	return newAnyArrayCondition(f, comparerLe, sqlable)
}

func (f *funcCol[T]) LeAll(sqlable ParametricSql) Condition {
	return newAllArrayCondition(f, comparerLe, sqlable)
}

var _ FuncColumn[any] = &funcCol[any]{}

func newFuncCol[T any](funcName string, inner ParametricSql) FuncColumn[T] {
	return &funcCol[T]{
		funcName: funcName,
		inner:    inner,
	}
}

func (f *funcCol[T]) SqlWithParams(paramsMap ParamsMap) (string, ParamsMap) {
	sql := f.funcName + "("
	var innerSql string
	innerSql, paramsMap = f.inner.SqlWithParams(paramsMap)
	sql += innerSql + ")"
	if f.Alias() != nil {
		sql += " AS " + *f.Alias()
	}
	return sql, paramsMap
}
