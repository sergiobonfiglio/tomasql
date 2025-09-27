package goql

type mockTable string

func (t mockTable) Columns() []Column {
	panic("implement me")
}

func (t mockTable) Alias() *string {
	return nil
}

func (t mockTable) sqlWithParams(paramsMap ParamsMap) (string, ParamsMap) {
	panic("implement me")
}

func (t mockTable) SQL() (sql string, params []any) {
	panic("implement me")
}

func (t mockTable) AsNamedSubQuery(s string) SQLable {
	panic("implement me")
}

func (t mockTable) AsSubQuery() SQLable {
	panic("implement me")
}
func (t mockTable) TableName() string { return string(t) }
