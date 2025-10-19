package tomasql

import "strings"

type builderWithGroupBy struct {
	prevStage ParametricSql
	groupBy   []ParametricSql
	having    Condition
	params    ParamsMap
}

var _ BuilderWithGroupBy = &builderWithGroupBy{}

func newBuilderWithGroupBy(prev ParametricSql, groupBy []ParametricSql, having Condition) BuilderWithGroupBy {
	b := &builderWithGroupBy{
		prevStage: prev,
		groupBy:   groupBy,
		having:    having,
		params:    ParamsMap{},
	}
	return b
}

func (b *builderWithGroupBy) Having(condition Condition) BuilderWithHaving {
	b.having = condition
	return b
}

func (b *builderWithGroupBy) OrderBy(first SortColumn, columns ...SortColumn) BuilderWithOrderBy {
	return newBuilderWithOrderBy(b, append([]SortColumn{first}, columns...))
}

func (b *builderWithGroupBy) SQL() (sql string, params []any) {
	sql, paramsMap := b.SqlWithParams(b.params)
	return sql, paramsMap.ToSlice()
}

func (b *builderWithGroupBy) SqlWithParams(paramsMap ParamsMap) (string, ParamsMap) {
	b.params = paramsMap.AddAll(b.params)

	var sql string
	sql, b.params = b.prevStage.SqlWithParams(b.params)
	var groupBySql []string
	for _, col := range b.groupBy {
		var colSql string
		colSql, b.params = col.SqlWithParams(b.params)
		groupBySql = append(groupBySql, colSql)
	}

	havingSql := ""
	if b.having != nil {
		havingSql = " HAVING " + b.having.SQL(b.params)
	}

	return sql + " GROUP BY " + strings.Join(groupBySql, ", ") + havingSql, b.params
}

func (b *builderWithGroupBy) AsNamedSubQuery(alias string) Table {
	return newWithOptionalAlias(b, &alias)
}

func (b *builderWithGroupBy) AsSubQuery() SQLable {
	return newWithOptionalAlias(b, nil)
}
