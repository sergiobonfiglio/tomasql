package tomasql

import "fmt"

type Table interface {
	TableName() string
	// Columns() []Column

	// cannot define As() here because we need it to return a concrete table to be
	// able to access specific columns and generics won't work due to type erasure
	// As(x string)
	Alias() *string
	ParametricSql
}

type sqlableTable struct {
	table Table
}

// SqlableTable is the exported version of sqlableTable for use in generated code
type SqlableTable = sqlableTable

func newSqlableTable(t Table) *sqlableTable {
	return &sqlableTable{table: t}
}

// NewSqlableTable creates a new sqlable table wrapper - exported for use in generated code
func NewSqlableTable(t Table) *sqlableTable {
	return newSqlableTable(t)
}

func (s *sqlableTable) SqlWithParams(params ParamsMap, ctx RenderContext) (string, ParamsMap) {
	switch ctx {
	case DefinitionContext:
		tRef := s.table.TableName()
		if s.table.Alias() != nil {
			tRef += " AS " + *s.table.Alias()
		}
		return tRef, params
	case ReferenceContext:
		tRef := s.table.TableName()
		if s.table.Alias() != nil {
			tRef += " AS " + *s.table.Alias()
		}
		return tRef, params
	case OrderByContext:
		tRef := s.table.TableName()
		if s.table.Alias() != nil {
			tRef += " AS " + *s.table.Alias()
		}
		return tRef, params
	default:
		panic(fmt.Sprintf("sqlableTable.SqlWithParams: unexpected RenderContext %s", ctx))
	}
}

var _ ParametricSql = &sqlableTable{}

type tableDef struct {
	*withOptionalAlias
	columns []Column
}

var _ Table = &tableDef{}

func NewTableFromSubQuery(subQuery SQLable, alias string, columns []Column) (Table, []Column) {
	table := &tableDef{
		withOptionalAlias: newWithOptionalAlias(subQuery, &alias),
		columns:           []Column{},
	}

	for _, col := range columns {
		table.columns = append(table.columns, NewCol[any](col.Name(), table))
	}

	return table, table.columns
}

// tableRefWrapper is a simple wrapper around a Table that renders a table
// in reference usages (e.g. in column references or JOIN clauses).
type tableRefWrapper struct {
	table Table
}

// Alias implements Table.
func (t *tableRefWrapper) Alias() *string {
	return t.table.Alias()
}

// SqlWithParams implements Table.
func (t *tableRefWrapper) SqlWithParams(paramsMap ParamsMap, ctx RenderContext) (string, ParamsMap) {
	switch ctx {
	case DefinitionContext:
		if t.table.Alias() != nil {
			return *t.table.Alias(), paramsMap
		}
		return t.table.TableName(), paramsMap
	case ReferenceContext:
		if t.table.Alias() != nil {
			return *t.table.Alias(), paramsMap
		}
		return t.table.TableName(), paramsMap
	case OrderByContext:
		if t.table.Alias() != nil {
			return *t.table.Alias(), paramsMap
		}
		return t.table.TableName(), paramsMap
	default:
		panic(fmt.Sprintf("tableRefWrapper.SqlWithParams: unexpected RenderContext %s", ctx))
	}
}

// TableName implements Table.
func (t *tableRefWrapper) TableName() string {
	return t.table.TableName()
}

// TODO: do we need to implement the interface methods other than SqlWithParams?
var _ Table = &tableRefWrapper{}
