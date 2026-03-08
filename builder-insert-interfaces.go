package tomasql

type BuilderInsertColumns interface {
	Columns(columns ...Column) BuilderInsertValues
}

type BuilderInsertValues interface {
	ValuesParam(row []any) BuilderInsertValues
	SQL() (sql string, params []any)
}
