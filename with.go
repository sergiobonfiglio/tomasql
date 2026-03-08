package tomasql

import "strings"

type CommonTableExpression struct {
	name  string
	query SQLable
}

func CTE(name string, query SQLable) *CommonTableExpression {
	if name == "" {
		panic("CTE: empty name")
	}
	if query == nil {
		panic("CTE: nil query")
	}
	return &CommonTableExpression{name: name, query: query}
}

func (c *CommonTableExpression) Name() string {
	return c.name
}

func (c *CommonTableExpression) Table() Table {
	return &cteTable{name: c.name}
}

func (c *CommonTableExpression) SqlWithParams(params ParamsMap, _ RenderContext) (string, ParamsMap) {
	innerSQL, params := c.query.SqlWithParams(params, OutputContext)
	return c.name + " AS (" + innerSQL + ")", params
}

type cteTable struct {
	name  string
	alias *string
}

func (t *cteTable) TableName() string {
	return t.name
}

func (t *cteTable) Alias() *string {
	return t.alias
}

func (t *cteTable) SqlWithParams(params ParamsMap, _ RenderContext) (string, ParamsMap) {
	if t.alias != nil {
		return t.name + " AS " + *t.alias, params
	}
	return t.name, params
}

var _ Table = &cteTable{}

type withBuilder struct {
	ctes []*CommonTableExpression
}

var _ WithBuilder = &withBuilder{}

func With(ctes ...*CommonTableExpression) WithBuilder {
	for _, cte := range ctes {
		if cte == nil {
			panic("With: nil CTE")
		}
	}
	return &withBuilder{ctes: ctes}
}

func (w *withBuilder) Select(first ParametricSql, columns ...ParametricSql) BuilderWithSelect {
	return newBuilderWithSelectWithCTEs(false, w.ctes, first, columns...)
}

func (w *withBuilder) SelectCols(first Column, columns ...Column) BuilderWithSelect {
	convertedColumns := make([]ParametricSql, len(columns))
	for i, col := range columns {
		convertedColumns[i] = col
	}
	return newBuilderWithSelectWithCTEs(false, w.ctes, first, convertedColumns...)
}

func (w *withBuilder) SelectDistinct(first ParametricSql, columns ...ParametricSql) BuilderWithSelect {
	return newBuilderWithSelectWithCTEs(true, w.ctes, first, columns...)
}

func (w *withBuilder) SelectDistinctCols(first Column, columns ...Column) BuilderWithSelect {
	convertedColumns := make([]ParametricSql, len(columns))
	for i, col := range columns {
		convertedColumns[i] = col
	}
	return newBuilderWithSelectWithCTEs(true, w.ctes, first, convertedColumns...)
}

func (w *withBuilder) SelectAll() BuilderWithSelect {
	return newBuilderWithSelectAllWithCTEs(false, w.ctes)
}

func (w *withBuilder) SelectDistinctAll() BuilderWithSelect {
	return newBuilderWithSelectAllWithCTEs(true, w.ctes)
}

type cteCarrier interface {
	getCTEs() []*CommonTableExpression
}

func ctesFrom(sql ParametricSql) []*CommonTableExpression {
	if carrier, ok := sql.(cteCarrier); ok {
		return carrier.getCTEs()
	}
	return nil
}

func renderWithCTEs(params ParamsMap, ctes []*CommonTableExpression, ctx RenderContext, render func(ParamsMap) (string, ParamsMap)) (string, ParamsMap) {
	if len(ctes) == 0 || isWithSuppressed(ctx) {
		return render(params)
	}

	parts := make([]string, 0, len(ctes))
	for _, cte := range ctes {
		var cteSQL string
		cteSQL, params = cte.SqlWithParams(params, DefinitionContext)
		parts = append(parts, cteSQL)
	}

	body, params := render(params)
	return "WITH " + strings.Join(parts, ", ") + " " + body, params
}
