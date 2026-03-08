package tomasql

type builderWithWhere struct {
	prevStage ParametricSql
	params    ParamsMap
	where     Condition
	ctes      []*CommonTableExpression
}

var _ BuilderWithWhere = &builderWithWhere{}

func newBuilderWithWhere(prev *builderWithJoin, where Condition) BuilderWithWhere {
	b := &builderWithWhere{
		prevStage: prev,
		where:     where,
		ctes:      ctesFrom(prev),
	}
	return b
}

func (b *builderWithWhere) AsNamedSubQuery(alias string) Table {
	return newWithOptionalAlias(b, &alias)
}

func (b *builderWithWhere) AsSubQuery() SQLable {
	return newWithOptionalAlias(b, nil)
}

func (b *builderWithWhere) GroupBy(column ParametricSql, columns ...ParametricSql) BuilderWithGroupBy {
	return newBuilderWithGroupBy(b, append([]ParametricSql{column}, columns...), nil)
}

func (b *builderWithWhere) OrderBy(column SortColumn, column2 ...SortColumn) BuilderWithOrderBy {
	return newBuilderWithOrderBy(b, append([]SortColumn{column}, column2...))
}

func (b *builderWithWhere) getCTEs() []*CommonTableExpression {
	return b.ctes
}

func (b *builderWithWhere) SqlWithParams(params ParamsMap, ctx RenderContext) (string, ParamsMap) {
	return renderWithCTEs(params, b.ctes, ctx, func(params ParamsMap) (string, ParamsMap) {
		b.params = params.AddAll(b.params)
		var sql string
		sql, b.params = b.prevStage.SqlWithParams(b.params, suppressWith(ctx))
		whereStr := ""
		if b.where != nil {
			whereStr = " WHERE " + b.where.SQL(b.params)
		}
		return sql + whereStr, b.params
	})
}

func (b *builderWithWhere) SQL() (sql string, params []any) {
	sql, paramsMap := b.SqlWithParams(b.params, OutputContext)
	return sql, paramsMap.ToSlice()
}
