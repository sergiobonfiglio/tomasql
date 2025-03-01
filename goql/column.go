package goql

type Table interface {
	TableName() string
	Columns() []Column

	// cannot define As() here because we need it to return a concrete table to be
	// able to access specific columns and generics won't work due to type erasure
	//As(x string)
	Alias() *string
	SQLable
}

type sqlableTable struct {
	table Table
}

func newSqlableTable(t Table) *sqlableTable {
	return &sqlableTable{table: t}
}

func (s *sqlableTable) SQL() (string, []any) {
	sql, params := s.sqlWithParams(ParamsMap{})
	return sql, params.ToSlice()
}

func (s *sqlableTable) sqlWithParams(params ParamsMap) (string, ParamsMap) {
	tRef := s.table.TableName()
	if s.table.Alias() != nil {
		tRef += " AS " + *s.table.Alias()
	}
	return tRef, params
}

func (s *sqlableTable) AsNamedSubQuery(alias string) SQLable {
	return newWithOptionalAlias(s, &alias)
}

func (b *sqlableTable) AsSubQuery() SQLable {
	return newWithOptionalAlias(b, nil)
}

var _ SQLable = &sqlableTable{}

type Comparable[T any] interface {
	Eq(other Column) Condition
	EqParam(other T) Condition
	Gt(other Column) Condition
	GtParam(other T) Condition
	Ge(other Column) Condition
	GeParam(other T) Condition
	Lt(other Column) Condition
	LtParam(other T) Condition
	Le(other Column) Condition
	LeParam(other T) Condition
}

type SetComparable[T any] interface {
	InArray(array []T) Condition
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
}

type Column interface {
	Name() string
	Table() Table

	//As(x string) Column
	Alias() *string

	Asc() SortColumn
	Desc() SortColumn

	getType() colTypeTag
	getRef() string
}

type ParametricSql interface {
	sqlWithParams(ParamsMap) (string, ParamsMap)
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
)

type Col[T any] struct {
	name         string
	table        Table
	concreteType colTypeTag
	alias        *string
	Comparable[T]
}

func (c Col[T]) SQL() (string, []any) {
	sql, params := c.sqlWithParams(nil)
	return sql, params.ToSlice()
}

func (c Col[T]) sqlWithParams(_ ParamsMap) (string, ParamsMap) {
	colRef := c.getRef()
	if c.Alias() != nil {
		colRef += " AS " + *c.Alias()
	}
	return colRef, ParamsMap{}
}

var _ Column = &Col[int]{}
var _ ParametricSql = &Col[int]{} // Ensure Col implements ParametricSql

func (c Col[T]) Alias() *string {
	return c.alias
}

func (c Col[T]) AsNamedSubQuery(x string) SQLable {
	return c.As(x)
}

func (c Col[T]) AsSubQuery() SQLable {
	return newWithOptionalAlias(c, nil)
}

func (c Col[T]) As(x string) SQLable {
	c.alias = &x
	return c
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
		//this is the case for "*"
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

func newCol[T any](name string, table Table) *Col[T] {
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
		panic("unknown type")
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
