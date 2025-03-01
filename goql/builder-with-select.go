package goql

import "strings"

type builderWithSelect struct {
	selectColumns []ParametricSql
	distinct      bool
	params        ParamsMap
}

var _ BuilderWithSelect = &builderWithSelect{}

func newBuilderWithSelect(distinct bool, first ParametricSql, columns ...ParametricSql) BuilderWithSelect {
	b := &builderWithSelect{
		selectColumns: append([]ParametricSql{first}, columns...),
		distinct:      distinct,
		params:        ParamsMap{},
	}
	return b
}

func (b *builderWithSelect) AsNamedSubQuery(alias string) SQLable {
	return newWithOptionalAlias(b, &alias)
}

func (b *builderWithSelect) AsSubQuery() SQLable {
	return newWithOptionalAlias(b, nil)
}

func (b *builderWithSelect) From(t SQLable) BuilderWithFrom {
	return newBuilderWithFrom(b, t)
}

func (b *builderWithSelect) sqlWithParams(params ParamsMap) (string, ParamsMap) {
	var colStr []string

	b.params = params.AddAll(b.params)

	for _, col := range b.selectColumns {
		// potentially, here we could check if the column is a subquery without alias that
		// would cause a grammar error in SQL
		var sql string
		sql, b.params = col.sqlWithParams(b.params)

		colStr = append(colStr, sql)
	}
	distinctStr := ""
	if b.distinct {
		distinctStr = "DISTINCT "
	}

	return "SELECT " + distinctStr + strings.Join(colStr, ", "), ParamsMap{}
}

func (b *builderWithSelect) SQL() (string, []any) {
	sql, params := b.sqlWithParams(b.params)
	return sql, params.ToSlice()
}

type builderWithSelectAll struct {
	*builderWithSelect
}

var _ BuilderWithSelect = &builderWithSelectAll{}

func newBuilderWithSelectAll(distinct bool) BuilderWithSelect {
	return &builderWithSelectAll{
		builderWithSelect: &builderWithSelect{distinct: distinct},
	}
}

func (b *builderWithSelectAll) From(t SQLable) BuilderWithFrom {
	return newBuilderWithFrom(b, t)
}

func (b *builderWithSelectAll) sqlWithParams(params ParamsMap) (string, ParamsMap) {
	distinctStr := ""
	if b.distinct {
		distinctStr = "DISTINCT "
	}
	return "SELECT " + distinctStr + "*", params
}
func (b *builderWithSelectAll) SQL() (string, []any) {
	sql, params := b.sqlWithParams(b.params)
	return sql, params.ToSlice()
}

type withOptionalAlias struct {
	SQLable
	alias *string
}

func newWithOptionalAlias(sqlable SQLable, alias *string) SQLable {
	return &withOptionalAlias{SQLable: sqlable, alias: alias}
}

func (b *withOptionalAlias) As(alias *string) SQLable {
	b.alias = alias
	return b
}

func (b *withOptionalAlias) sqlWithParams(params ParamsMap) (string, ParamsMap) {
	sql, params := b.SQLable.sqlWithParams(params)
	if b.alias == nil {
		return "(" + sql + ")", params
	}
	return "(" + sql + ") AS " + *b.alias, params
}
