package tomasql

type builderWithFrom struct {
	prevStage ParametricSql
	fromTable ParametricSql
	params    ParamsMap
}

var _ BuilderWithTables = &builderWithFrom{}

func newBuilderWithFrom(prev, from ParametricSql) BuilderWithTables {
	b := &builderWithFrom{
		prevStage: prev,
		fromTable: from,
	}

	return b
}

func (b *builderWithFrom) Joins(joinItems ...*JoinItem) BuilderWithTables {
	return _addJoins(b, joinItems...)
}

func _addJoins(b BuilderWithTables, joinItems ...*JoinItem) BuilderWithTables {
	if len(joinItems) == 0 {
		return b
	}
	for _, joinItem := range joinItems {
		if joinItem == nil {
			continue
		}
		b = b.Join(joinItem.Target).On(joinItem.OnCondition)
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

func (b *builderWithFrom) AsNamedSubQuery(alias string) Table {
	return newWithOptionalAlias(b, &alias)
}

func (b *builderWithFrom) AsSubQuery() SQLable {
	return newWithOptionalAlias(b, nil)
}

func (b *builderWithFrom) SqlWithParams(params ParamsMap, ctx RenderContext) (string, ParamsMap) {
	b.params = params.AddAll(b.params)
	var sql string
	sql, b.params = b.prevStage.SqlWithParams(b.params, ctx)
	var sqlTable string
	sqlTable, b.params = b.fromTable.SqlWithParams(b.params, DefinitionContext)
	return sql + " FROM " + sqlTable, b.params
}

func (b *builderWithFrom) SQL() (sql string, params []any) {
	sql, paramsMap := b.SqlWithParams(b.params, OutputContext)
	return sql, paramsMap.ToSlice()
}
