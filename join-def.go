package tomasql

type JoinType string

const (
	LeftJoin  = JoinType("LEFT")
	InnerJoin = JoinType("")
	RightJoin = JoinType("RIGHT")
)

type joinDef struct {
	joinType      JoinType
	joinTable     Table
	joinCondition Condition
}

var _ ParametricSql = &joinDef{}

func newJoinDef(joinType JoinType, joinTable Table, cond Condition) *joinDef {
	return &joinDef{
		joinType:      joinType,
		joinTable:     joinTable,
		joinCondition: cond,
	}
}

func (j *joinDef) SqlWithParams(paramsMap ParamsMap, ctx RenderContext) (string, ParamsMap) {
	joinStr := ""
	if j.joinType != "" {
		joinStr += string(j.joinType) + " "
	}
	joinTableSql, paramsMap := j.joinTable.SqlWithParams(paramsMap, ReferenceContext)
	joinStr += "JOIN " + joinTableSql
	if j.joinCondition != nil {
		joinStr += " ON " + j.joinCondition.SQL(paramsMap)
	}
	return joinStr, paramsMap
}
