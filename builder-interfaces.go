package tomasql

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
	GroupBy(ParametricSql, ...ParametricSql) BuilderWithGroupBy
	OrderBy(SortColumn, ...SortColumn) BuilderWithOrderBy
}

type BuilderWithGroupBy interface {
	BuilderWithHaving
	Having(Condition) BuilderWithHaving
}

type BuilderWithHaving interface {
	SubQueryable
	OrderBy(SortColumn, ...SortColumn) BuilderWithOrderBy
}

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

type ParametricSql interface {
	// SqlWithParams renders SQL with awareness of the context (SELECT, WHERE, etc.)
	SqlWithParams(ParamsMap, RenderContext) (string, ParamsMap)
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
	ParametricSql

	// Column returns the underlying Column that this SortColumn represents or nil if it is not a Column (e.g. subquery).
	Column() Column
}
