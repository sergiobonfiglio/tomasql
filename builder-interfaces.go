package tomasql

type Builder1 interface {
	// SelectCols selects distinct columns from the table. Equivalent to SelectDistinct but avoids
	// the need to convert Column to ParametricSql.
	SelectCols(first Column, columns ...Column) BuilderWithSelect
	Select(ParametricSql, ...ParametricSql) BuilderWithSelect

	// SelectDistinctCols selects distinct columns from the table. Equivalent to SelectDistinct but avoids
	// the need to convert Column to ParametricSql.
	SelectDistinctCols(Column, ...Column) BuilderWithSelect
	SelectDistinct(ParametricSql, ...ParametricSql) BuilderWithSelect
	SelectDistinctAll() BuilderWithSelect
	SelectAll() BuilderWithSelect
}

type BuilderWithSelect interface {
	SubQueryable
	From(Table) BuilderWithTables
}

type BuilderWithTables interface {
	SubQueryable
	LeftJoin(Table) BuilderWithJoin
	Join(Table) BuilderWithJoin
	RightJoin(Table) BuilderWithJoin

	Joins(...*JoinItem) BuilderWithTables

	Where(Condition) BuilderWithWhere
	GroupBy(ParametricSql, ...ParametricSql) BuilderWithGroupBy
	OrderBy(SortColumn, ...SortColumn) BuilderWithOrderBy
}

type BuilderWithJoin interface {
	SubQueryable
	On(Condition) BuilderWithTables
}

type BuilderWithWhere interface {
	SubQueryable
	OrderBy(SortColumn, ...SortColumn) BuilderWithOrderBy
}

type BuilderWithGroupBy interface {
	BuilderWithHaving
	Having(Condition) BuilderWithHaving
}

type BuilderWithHaving BuilderWithWhere

type BuilderWithOrderBy interface {
	SubQueryable
	Limit(int) BuilderWithLimit
}

type BuilderWithLimit interface {
	SubQueryable
	Offset(int) SQLable
}

type SQLable interface {
	ParametricSql
	SQL() (sql string, params []any)
}

type SubQueryable interface {
	SQLable
	AsNamedSubQuery(string) Table
	AsSubQuery() SQLable
}

type SortDirection string

const (
	OrderByAsc  SortDirection = "ASC"
	OrderByDesc SortDirection = "DESC"
)

type SortColumn interface {
	Name() string
	Table() Table
	Direction() SortDirection
	SQL() string

	// Column returns the underlying Column that this SortColumn represents or nil if it is not a Column (e.g. subquery).
	Column() Column

	getRef() string
}

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}
