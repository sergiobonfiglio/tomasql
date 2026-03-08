package tomasql

import "strings"

type builderWithSelect struct {
	selectColumns []ParametricSql
	distinct      bool
	params        ParamsMap
	ctes          []*CommonTableExpression
}

var _ BuilderWithSelect = &builderWithSelect{}

func newBuilderWithSelect(distinct bool, first ParametricSql, columns ...ParametricSql) BuilderWithSelect {
	return newBuilderWithSelectWithCTEs(distinct, nil, first, columns...)
}

func newBuilderWithSelectWithCTEs(distinct bool, ctes []*CommonTableExpression, first ParametricSql, columns ...ParametricSql) BuilderWithSelect {
	b := &builderWithSelect{
		selectColumns: append([]ParametricSql{first}, columns...),
		distinct:      distinct,
		params:        ParamsMap{},
		ctes:          ctes,
	}
	return b
}

func (b *builderWithSelect) AsNamedSubQuery(alias string) Table {
	return newWithOptionalAlias(b, &alias)
}

func (b *builderWithSelect) AsSubQuery() SQLable {
	return newWithOptionalAlias(b, nil)
}

func (b *builderWithSelect) From(t Table) BuilderWithTables {
	return newBuilderWithFrom(b, t)
}

func (b *builderWithSelect) getCTEs() []*CommonTableExpression {
	return b.ctes
}

func (b *builderWithSelect) SqlWithParams(params ParamsMap, ctx RenderContext) (string, ParamsMap) {
	return renderWithCTEs(params, b.ctes, ctx, func(params ParamsMap) (string, ParamsMap) {
		var colStr []string
		b.params = params.AddAll(b.params)
		for _, col := range b.selectColumns {
			var sql string
			sql, b.params = col.SqlWithParams(b.params, DefinitionContext)
			colStr = append(colStr, sql)
		}
		distinctStr := ""
		if b.distinct {
			distinctStr = "DISTINCT "
		}
		return "SELECT " + distinctStr + strings.Join(colStr, ", "), b.params
	})
}

func (b *builderWithSelect) SQL() (sql string, params []any) {
	sql, paramsMap := b.SqlWithParams(b.params, OutputContext)
	return sql, paramsMap.ToSlice()
}

type builderWithSelectAll struct {
	*builderWithSelect
}

var _ BuilderWithSelect = &builderWithSelectAll{}

func newBuilderWithSelectAll(distinct bool) BuilderWithSelect {
	return newBuilderWithSelectAllWithCTEs(distinct, nil)
}

func newBuilderWithSelectAllWithCTEs(distinct bool, ctes []*CommonTableExpression) BuilderWithSelect {
	return &builderWithSelectAll{
		builderWithSelect: &builderWithSelect{distinct: distinct, ctes: ctes},
	}
}

func (b *builderWithSelectAll) From(t Table) BuilderWithTables {
	return newBuilderWithFrom(b, t)
}

func (b *builderWithSelectAll) SqlWithParams(params ParamsMap, ctx RenderContext) (string, ParamsMap) {
	return renderWithCTEs(params, b.ctes, ctx, func(params ParamsMap) (string, ParamsMap) {
		distinctStr := ""
		if b.distinct {
			distinctStr = "DISTINCT "
		}
		return "SELECT " + distinctStr + "*", params
	})
}

func (b *builderWithSelectAll) SQL() (sql string, params []any) {
	sql, paramsMap := b.SqlWithParams(b.params, OutputContext)
	return sql, paramsMap.ToSlice()
}

type withOptionalAlias struct {
	SQLable
	alias *string
}

// TableName implements the Table interface for withOptionalAlias.
// Returns the alias if present, otherwise an empty string.
func (b *withOptionalAlias) TableName() string {
	if b.alias != nil {
		return *b.alias
	}
	return ""
}

func (b *withOptionalAlias) Alias() *string {
	return b.alias
}

var _ Table = &withOptionalAlias{}

func newWithOptionalAlias(sqlable SQLable, alias *string) *withOptionalAlias {
	return &withOptionalAlias{SQLable: sqlable, alias: alias}
}

func (b *withOptionalAlias) As(alias *string) SQLable {
	b.alias = alias
	return b
}

func (b *withOptionalAlias) SqlWithParams(params ParamsMap, ctx RenderContext) (string, ParamsMap) {
	sql, params := b.SQLable.SqlWithParams(params, ctx)
	if b.alias == nil {
		return "(" + sql + ")", params
	}
	// Subquery aliases should always be rendered (they're table aliases, not column aliases)
	return "(" + sql + ") AS " + *b.alias, params
}
