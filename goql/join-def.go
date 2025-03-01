package goql

type JoinType string

const (
	LeftJoin  = JoinType("LEFT")
	InnerJoin = JoinType("")
	RightJoin = JoinType("RIGHT")
)

type joinDef struct {
	params        ParamsMap
	joinType      JoinType
	joinTable     Table
	joinCondition Condition
}

func newJoinDef(joinType JoinType, joinTable Table, cond Condition, params ParamsMap) *joinDef {
	return &joinDef{
		joinType:      joinType,
		joinTable:     joinTable,
		joinCondition: cond,
		params:        params,
	}
}

func (j *joinDef) SQL() string {
	joinStr := ""
	if j.joinType != "" {
		joinStr += string(j.joinType) + " "
	}
	joinStr += "JOIN " + j.joinTable.TableName()

	if j.joinTable.Alias() != nil {
		joinStr += " AS " + *j.joinTable.Alias()
	}

	if j.joinCondition != nil {
		joinStr += " ON " + j.joinCondition.SQL(j.params)
	}
	return joinStr
}
