package tomasql

import "fmt"

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
	// ILike(other ParametricSql) Condition
	// ILikeParam(pattern string) Condition
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
	getRef() string

	SQLable
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

func (c Col[T]) SQL() (sql string, params []any) {
	sql, paramsMap := c.SqlWithParams(nil)
	return sql, paramsMap.ToSlice()
}

func (c Col[T]) SqlWithParams(params ParamsMap) (string, ParamsMap) {
	colRef := c.getRef()
	if c.Alias() != nil {
		colRef += " AS " + *c.Alias()
	}
	return colRef, params
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

func (c Col[T]) getRef() string {
	if c.Table() == nil {
		// this is the case for "*"
		return c.Name()
	}

	tRef := c.Table().TableName()
	if c.Table().Alias() != nil {
		tRef = *c.Table().Alias()
	}
	return tRef + "." + c.Name()
}

type SortCol[T any] struct {
	col       Column
	direction SortDirection
}

var _ SortColumn = &SortCol[int]{} // Ensure SortCol implements SortColumn

func (s *SortCol[T]) Name() string {
	return s.col.Name()
}

func (s *SortCol[T]) Column() Column {
	return s.col
}

func (s *SortCol[T]) Table() Table {
	return s.col.Table()
}

func (s *SortCol[T]) getRef() string {
	return s.col.getRef()
}

func (s *SortCol[T]) Direction() SortDirection {
	return s.direction
}

func (s *SortCol[T]) SQL() string {
	return s.getRef() + " " + string(s.direction)
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
