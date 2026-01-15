package tomasql

import (
	"fmt"
)

type Comparable interface {
	Eq(other ParametricSql) Condition
	Gt(other ParametricSql) Condition
	Ge(other ParametricSql) Condition
	Lt(other ParametricSql) Condition
	Le(other ParametricSql) Condition
	IsNull() Condition
	IsNotNull() Condition
	Like(other ParametricSql) Condition
	LikeParam(pattern string) Condition
}

type ComparableParam[T any] interface {
	EqParam(other T) Condition
	GtParam(other T) Condition
	GeParam(other T) Condition
	LtParam(other T) Condition
	LeParam(other T) Condition
	// InArray(array []T) Condition
}

type SetComparable interface {
	In(sqlable ParametricSql) Condition
	EqAny(sqlable ParametricSql) Condition
	EqAll(sqlable ParametricSql) Condition
	GtAny(sqlable ParametricSql) Condition
	GtAll(sqlable ParametricSql) Condition
	GeAny(sqlable ParametricSql) Condition
	GeAll(sqlable ParametricSql) Condition
	LtAny(sqlable ParametricSql) Condition
	LtAll(sqlable ParametricSql) Condition
	LeAny(sqlable ParametricSql) Condition
	LeAll(sqlable ParametricSql) Condition
	IsNull() Condition
	IsNotNull() Condition
}

type Column interface {
	Name() string
	Table() Table

	As(x string) Column
	Alias() *string

	Asc() SortColumn
	Desc() SortColumn

	getType() colTypeTag

	ParametricSql
	Comparable
	SetComparable
}

type ParametricSql interface {
	SqlWithParams(ParamsMap) (string, ParamsMap)
}

type colTypeTag string

const (
	intTag     = colTypeTag("int")
	int32Tag   = colTypeTag("int32")
	int64Tag   = colTypeTag("int64")
	float64Tag = colTypeTag("float64")
	float32Tag = colTypeTag("float32")
	stringTag  = colTypeTag("string")
	boolTag    = colTypeTag("bool")
	anyTag     = colTypeTag("any")
)

type Col[T any] struct {
	name         string
	table        Table
	concreteType colTypeTag
	alias        *string
	ComparableParam[T]
}

func (c Col[T]) SqlWithParams(params ParamsMap) (string, ParamsMap) {

	if c.Table() == nil {
		// this is the case for "*"
		return c.Name(), params
	}

	table := c.getTableForRef()
	var tRef string
	tRef, params = table.SqlWithParams(params)

	if c.Alias() == nil {
		return tRef + "." + c.Name(), params
	}
	return tRef + "." + c.Name() + " AS " + *c.Alias(), params
}

func (c Col[T]) getTableForRef() tableRefWrapper {
	return tableRefWrapper{table: c.Table()}
}

var (
	_ Column        = &Col[any]{}
	_ ParametricSql = &Col[int]{} // Ensure Col implements ParametricSql
)

func (c Col[T]) Alias() *string {
	return c.alias
}

func (c Col[T]) As(x string) Column {
	c.alias = &x
	return c
}

func (c Col[T]) ToSorted(direction SortDirection) SortColumn {
	return &SortCol[T]{
		col:       c,
		direction: direction,
	}
}

func (c Col[T]) Asc() SortColumn {
	return &SortCol[T]{
		col:       c,
		direction: OrderByAsc,
	}
}

func (c Col[T]) Desc() SortColumn {
	return &SortCol[T]{
		col:       c,
		direction: OrderByDesc,
	}
}

// colRefWrapper renders a column reference (just the name or alias (plus table name or alias), no "AS")
type colRefWrapper struct {
	col Column
}

func (crw colRefWrapper) SqlWithParams(params ParamsMap) (string, ParamsMap) {
	if crw.col.Alias() != nil {
		return *crw.col.Alias(), params
	}
	return crw.col.Name(), params
}

func NewCol[T any](name string, table Table) *Col[T] {
	var emptyT T
	var tag colTypeTag
	switch any(emptyT).(type) {
	case int:
		tag = intTag
	case int32:
		tag = int32Tag
	case int64:
		tag = int64Tag
	case string:
		tag = stringTag
	case float64:
		tag = float64Tag
	case float32:
		tag = float32Tag
	case bool:
		tag = boolTag
	default:
		tag = anyTag
	}

	return &Col[T]{
		name:         name,
		table:        table,
		concreteType: tag,
	}
}

func (c Col[T]) Name() string {
	return c.name
}

func (c Col[T]) Table() Table {
	return c.table
}

func (c Col[T]) getType() colTypeTag {
	return c.concreteType
}

type Number interface {
	~int | ~int32 | ~int64 | ~float32 | ~float64
}

func NewFixedCol[T string | Number](val T, alias *string) Column {
	valStr := fmt.Sprintf("%v", val)
	var col Column = NewCol[T](valStr, nil)
	if alias != nil {
		col = col.As(*alias)
	}

	return col
}

type SortCol[T any] struct {
	// either col or subQuery will be set
	col       Column
	subQuery  ParametricSql
	direction SortDirection
}

var _ SortColumn = &SortCol[int]{} // Ensure SortCol implements SortColumn

func (s *SortCol[T]) Column() Column {
	return s.col
}

func (s *SortCol[T]) SqlWithParams(params ParamsMap) (string, ParamsMap) {

	if s.subQuery != nil {
		subQueryStr, pm := s.subQuery.SqlWithParams(params)
		return fmt.Sprintf("%s %s", subQueryStr, string(s.direction)), pm
	}

	sortRef := ""
	if s.col.Table() != nil {
		table := tableRefWrapper{table: s.col.Table()}
		sortRef, params = table.SqlWithParams(params)
		sortRef += "."
	}

	var colRef string
	colRef, params = colRefWrapper{col: s.col}.SqlWithParams(params)

	return sortRef + colRef + " " + string(s.direction), params
}
