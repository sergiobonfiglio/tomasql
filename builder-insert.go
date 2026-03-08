package tomasql

import "strings"

func InsertInto(table Table) BuilderInsertColumns {
	return &insertColumnsBuilder{
		table: table,
	}
}

type insertColumnsBuilder struct {
	table   Table
	columns []Column
}

func (b *insertColumnsBuilder) Columns(columns ...Column) BuilderInsertValues {
	b.columns = columns
	return &insertValuesBuilder{
		table:   b.table,
		columns: b.columns,
		params:  ParamsMap{},
		rows:    [][]any{},
	}
}

type insertValuesBuilder struct {
	table   Table
	columns []Column
	params  ParamsMap
	rows    [][]any
}

func (b *insertValuesBuilder) ValuesParam(row []any) BuilderInsertValues {
	if len(row) != len(b.columns) {
		panic("number of values must match number of columns")
	}
	b.rows = append(b.rows, row)
	return b
}

func (b *insertValuesBuilder) SqlWithParams(params ParamsMap, ctx RenderContext) (string, ParamsMap) {
	if len(b.rows) == 0 {
		panic("at least one row must be provided")
	}

	b.params = params.AddAll(b.params)

	columnNames := make([]string, len(b.columns))
	for i, col := range b.columns {
		columnNames[i] = col.Name()
	}

	var rowValues []string
	paramIndex := len(b.params) + 1
	for _, row := range b.rows {
		placeholders := make([]string, len(row))
		for j, val := range row {
			placeholders[j] = GetDialect().Placeholder(paramIndex)
			b.params[val] = paramIndex
			paramIndex++
		}
		rowValues = append(rowValues, "("+strings.Join(placeholders, ", ")+")")
	}

	sql := "INSERT INTO " + b.table.TableName() + " (" + strings.Join(columnNames, ", ") + ") VALUES " + strings.Join(rowValues, ", ")

	return sql, b.params
}

func (b *insertValuesBuilder) SQL() (string, []any) {
	sql, paramsMap := b.SqlWithParams(b.params, OutputContext)
	return sql, paramsMap.ToSlice()
}

var _ BuilderInsertColumns = &insertColumnsBuilder{}
var _ BuilderInsertValues = &insertValuesBuilder{}
