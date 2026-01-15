package tomasql

import (
	"fmt"
	"strings"
)

type builderWithOrderBy struct {
	prevStage ParametricSql
	orderBy   []SortColumn
	limit     *int
	offset    *int
	params    ParamsMap
}

var (
	_ BuilderWithOrderBy = &builderWithOrderBy{}
	_ BuilderWithLimit   = &builderWithOrderBy{}
)

func newBuilderWithOrderBy(prev ParametricSql, orderBy []SortColumn) BuilderWithOrderBy {
	b := &builderWithOrderBy{
		prevStage: prev,
		orderBy:   orderBy,
		params:    ParamsMap{},
	}

	return b
}

func (b *builderWithOrderBy) AsNamedSubQuery(alias string) Table {
	return newWithOptionalAlias(b, &alias)
}

func (b *builderWithOrderBy) AsSubQuery() SQLable {
	return newWithOptionalAlias(b, nil)
}

func (b *builderWithOrderBy) Limit(i int) BuilderWithLimit {
	b.limit = &i
	return b
}

func (b *builderWithOrderBy) Offset(i int) SQLable {
	b.offset = &i
	return b
}

func (b *builderWithOrderBy) SqlWithParams(params ParamsMap) (string, ParamsMap) {
	b.params = params.AddAll(b.params)
	var out string
	out, b.params = b.prevStage.SqlWithParams(b.params)

	if len(b.orderBy) > 0 {
		out += " ORDER BY "
		var orderStr []string
		for _, col := range b.orderBy {
			var sortStr string
			sortStr, b.params = col.SqlWithParams(b.params)
			orderStr = append(orderStr, sortStr)
		}
		out += strings.Join(orderStr, ", ")
	}

	if b.limit != nil {
		out += fmt.Sprintf(" LIMIT %d", *b.limit)
	}

	if b.offset != nil {
		out += fmt.Sprintf(" OFFSET %d", *b.offset)
	}

	return out, b.params
}

func (b *builderWithOrderBy) SQL() (sql string, params []any) {
	sql, paramsMap := b.SqlWithParams(b.params)
	return sql, paramsMap.ToSlice()
}
