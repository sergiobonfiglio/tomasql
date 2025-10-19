package tomasql

type builderWithWhere struct {
	prevStage ParametricSql
	params    ParamsMap
	where     Condition
}

var _ BuilderWithWhere = &builderWithWhere{}

func newBuilderWithWhere(prev *builderWithJoin, where Condition) BuilderWithWhere {
	b := &builderWithWhere{
		prevStage: prev,
		where:     where,
	}
	return b
}

func (b *builderWithWhere) AsNamedSubQuery(alias string) Table {
	return newWithOptionalAlias(b, &alias)
}

func (b *builderWithWhere) AsSubQuery() SQLable {
	return newWithOptionalAlias(b, nil)
}

func (b *builderWithWhere) OrderBy(column SortColumn, column2 ...SortColumn) BuilderWithOrderBy {
	return newBuilderWithOrderBy(b, append([]SortColumn{column}, column2...))
}

func (b *builderWithWhere) SqlWithParams(params ParamsMap) (string, ParamsMap) {
	b.params = params.AddAll(b.params)
	var sql string
	sql, b.params = b.prevStage.SqlWithParams(b.params)

	whereStr := ""
	if b.where != nil {
		whereStr = " WHERE " + b.where.SQL(b.params)
	}

	return sql + whereStr, b.params
}

func (b *builderWithWhere) SQL() (sql string, params []any) {
	sql, paramsMap := b.SqlWithParams(b.params)
	return sql, paramsMap.ToSlice()
}
