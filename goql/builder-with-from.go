package goql

type builderWithFrom struct {
	prevStage ParametricSql
	fromTable ParametricSql
	params    ParamsMap
}

var _ BuilderWithFrom = &builderWithFrom{}

func newBuilderWithFrom(prev ParametricSql, from ParametricSql) BuilderWithFrom {
	b := &builderWithFrom{
		prevStage: prev,
		fromTable: from,
	}

	return b
}

func (b *builderWithFrom) Join(t Table) BuilderWithJoin {
	return newBuilderWithJoin(b, InnerJoin, t)
}

func (b *builderWithFrom) LeftJoin(t Table) BuilderWithJoin {
	return newBuilderWithJoin(b, LeftJoin, t)
}

func (b *builderWithFrom) RightJoin(t Table) BuilderWithJoin {
	return newBuilderWithJoin(b, RightJoin, t)
}

func (b *builderWithFrom) Where(cond Condition) BuilderWithWhere {
	withJoin := (b.Join(nil)).(*builderWithJoin)
	return newBuilderWithWhere(withJoin, cond)
}

func (b *builderWithFrom) GroupBy(column ParametricSql, columns ...ParametricSql) BuilderWithGroupBy {
	return newBuilderWithGroupBy(b, append([]ParametricSql{column}, columns...), nil)
}

func (b *builderWithFrom) OrderBy(column SortColumn, columns ...SortColumn) BuilderWithOrderBy {
	return newBuilderWithOrderBy(b, append([]SortColumn{column}, columns...))
}
func (b *builderWithFrom) AsNamedSubQuery(alias string) SQLable {
	return newWithOptionalAlias(b, &alias)
}

func (b *builderWithFrom) AsSubQuery() SQLable {
	return newWithOptionalAlias(b, nil)
}

func (b *builderWithFrom) sqlWithParams(params ParamsMap) (string, ParamsMap) {
	b.params = params.AddAll(b.params)
	sql, params := b.prevStage.sqlWithParams(b.params)

	var sqlTable string
	sqlTable, params = b.fromTable.sqlWithParams(params)

	return sql + " FROM " + sqlTable, params
}

func (b *builderWithFrom) SQL() (string, []any) {
	sql, params := b.sqlWithParams(b.params)
	return sql, params.ToSlice()
}
