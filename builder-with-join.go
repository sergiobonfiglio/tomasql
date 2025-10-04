package tomasql

import "strings"

type builderWithJoin struct {
	prevStage ParametricSql
	joins     []*joinDef
	params    ParamsMap
}

var (
	_ BuilderWithJoin   = &builderWithJoin{}
	_ BuilderWithTables = &builderWithJoin{}
)

func newBuilderWithJoin(prev ParametricSql, joinType JoinType, joinTable Table) BuilderWithJoin {
	params := ParamsMap{}
	var joins []*joinDef
	if joinTable != nil {
		joins = append(joins, newJoinDef(joinType, joinTable, nil))
	}
	b := &builderWithJoin{
		prevStage: prev,
		joins:     joins,
		params:    params,
	}
	return b
}

func (b *builderWithJoin) AsNamedSubQuery(alias string) Table {
	return newWithOptionalAlias(b, &alias)
}

func (b *builderWithJoin) AsSubQuery() SQLable {
	return newWithOptionalAlias(b, nil)
}

func (b *builderWithJoin) On(condition Condition) BuilderWithTables {
	lastJoin := b.joins[len(b.joins)-1]
	lastJoin.joinCondition = condition
	return b
}

func (b *builderWithJoin) _join(joinType JoinType, t Table) BuilderWithJoin {
	if t != nil {
		b.joins = append(b.joins, newJoinDef(joinType, t, nil))
	}
	return b
}

func (b *builderWithJoin) Join(t Table) BuilderWithJoin {
	return b._join(InnerJoin, t)
}

func (b *builderWithJoin) LeftJoin(t Table) BuilderWithJoin {
	return b._join(LeftJoin, t)
}

func (b *builderWithJoin) RightJoin(t Table) BuilderWithJoin {
	return b._join(RightJoin, t)
}

func (b *builderWithJoin) Joins(joinItems ...*JoinItem) BuilderWithTables {
	return _addJoins(b, joinItems...)
}

func (b *builderWithJoin) Where(cond Condition) BuilderWithWhere {
	return newBuilderWithWhere(b, cond)
}

func (b *builderWithJoin) GroupBy(column ParametricSql, columns ...ParametricSql) BuilderWithGroupBy {
	return newBuilderWithGroupBy(b, append([]ParametricSql{column}, columns...), nil)
}

func (b *builderWithJoin) OrderBy(column SortColumn, columns ...SortColumn) BuilderWithOrderBy {
	return newBuilderWithOrderBy(b, append([]SortColumn{column}, columns...))
}

func (b *builderWithJoin) sqlWithParams(params ParamsMap) (string, ParamsMap) {
	b.params = params.AddAll(b.params)

	var out string
	out, b.params = b.prevStage.sqlWithParams(b.params)

	if len(b.joins) > 0 {
		var joinStr []string
		for _, join := range b.joins {
			var jstr string
			jstr, b.params = join.sqlWithParams(b.params)
			joinStr = append(joinStr, jstr)
		}
		join := strings.Join(joinStr, " ")
		out += " " + join
	}
	return out, b.params
}

func (b *builderWithJoin) SQL() (sql string, params []any) {
	sql, paramsMap := b.sqlWithParams(b.params)
	return sql, paramsMap.ToSlice()
}
