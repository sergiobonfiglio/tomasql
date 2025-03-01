package goql

import "strings"

type builderWithJoin struct {
	prevStage ParametricSql
	joins     []*joinDef
	params    ParamsMap
}

var _ BuilderWithJoin = &builderWithJoin{}
var _ BuilderWithOn = &builderWithJoin{}

func newBuilderWithJoin(prev ParametricSql, joinType JoinType, joinTable Table) BuilderWithJoin {
	params := ParamsMap{}
	var joins []*joinDef
	if joinTable != nil {
		joins = append(joins, newJoinDef(joinType, joinTable, nil, params))
	}
	b := &builderWithJoin{
		prevStage: prev,
		joins:     joins,
		params:    params,
	}
	return b
}

func (b *builderWithJoin) AsNamedSubQuery(alias string) SQLable {
	return newWithOptionalAlias(b, &alias)
}

func (b *builderWithJoin) AsSubQuery() SQLable {
	return newWithOptionalAlias(b, nil)
}

func (b *builderWithJoin) On(condition Condition) BuilderWithOn {
	lastJoin := b.joins[len(b.joins)-1]
	lastJoin.joinCondition = condition
	return b
}

func (b *builderWithJoin) _join(joinType JoinType, t Table) BuilderWithJoin {
	if t != nil {
		b.joins = append(b.joins, newJoinDef(joinType, t, nil, b.params))
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

	out, params := b.prevStage.sqlWithParams(b.params)

	if len(b.joins) > 0 {
		var joinStr []string
		for _, join := range b.joins {
			joinStr = append(joinStr, join.SQL())
		}
		join := strings.Join(joinStr, " ")
		out += " " + join
	}
	return out, params
}

func (b *builderWithJoin) SQL() (string, []any) {
	sql, params := b.sqlWithParams(b.params)
	return sql, params.ToSlice()
}
