package goql

type Builder1 interface {
	Select(ParametricSql, ...ParametricSql) BuilderWithSelect
	SelectDistinct(ParametricSql, ...ParametricSql) BuilderWithSelect
	SelectDistinctAll() BuilderWithSelect
	SelectAll() BuilderWithSelect
}

type BuilderWithSelect interface {
	SQLable
	From(SQLable) BuilderWithFrom
}

type BuilderWithFrom interface {
	SQLable
	LeftJoin(Table) BuilderWithJoin
	Join(Table) BuilderWithJoin
	RightJoin(Table) BuilderWithJoin
	Where(Condition) BuilderWithWhere
	GroupBy(ParametricSql, ...ParametricSql) BuilderWithGroupBy
	OrderBy(SortColumn, ...SortColumn) BuilderWithOrderBy
}

type BuilderWithJoin interface {
	SQLable
	On(Condition) BuilderWithOn
}

type BuilderWithOn interface {
	SQLable
	LeftJoin(Table) BuilderWithJoin
	Join(Table) BuilderWithJoin
	RightJoin(Table) BuilderWithJoin
	Where(Condition) BuilderWithWhere
	GroupBy(ParametricSql, ...ParametricSql) BuilderWithGroupBy
	OrderBy(SortColumn, ...SortColumn) BuilderWithOrderBy
}

type BuilderWithWhere interface {
	SQLable
	OrderBy(SortColumn, ...SortColumn) BuilderWithOrderBy
}

type BuilderWithGroupBy interface {
	BuilderWithHaving
	Having(Condition) BuilderWithHaving
}

type BuilderWithHaving BuilderWithWhere

type BuilderWithOrderBy interface {
	SQLable
	Limit(int) BuilderWithLimit
}

type BuilderWithLimit interface {
	SQLable
	Offset(int) SQLable
}

type SQLable interface {
	SQL() (sql string, params []any)
	sqlWithParams(ParamsMap) (string, ParamsMap)
	AsNamedSubQuery(string) SQLable
	AsSubQuery() SQLable
}

type SortDirection string

const (
	OrderByAsc  = "ASC"
	OrderByDesc = "DESC"
)

type SortColumn interface {
	Name() string
	Table() Table
	Direction() SortDirection
	SQL() string

	getRef() string
}
