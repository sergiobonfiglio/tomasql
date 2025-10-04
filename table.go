package tomasql

type Table interface {
	TableName() string
	// Columns() []Column

	// cannot define As() here because we need it to return a concrete table to be
	// able to access specific columns and generics won't work due to type erasure
	// As(x string)
	Alias() *string
	SQLable
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

func (s *sqlableTable) SQL() (sql string, params []any) {
	sql, paramsMap := s.sqlWithParams(ParamsMap{})
	return sql, paramsMap.ToSlice()
}

func (s *sqlableTable) sqlWithParams(params ParamsMap) (string, ParamsMap) {
	tRef := s.table.TableName()
	if s.table.Alias() != nil {
		tRef += " AS " + *s.table.Alias()
	}
	return tRef, params
}

var _ SQLable = &sqlableTable{}

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
