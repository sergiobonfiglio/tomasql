package pgres

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/sergiobonfiglio/tomasql"
)

type sqlableArray struct {
	array  sql.Scanner
	params tomasql.ParamsMap
}

var _ tomasql.ParametricSql = &sqlableArray{}

func SQLArray[T any](array []T) tomasql.ParametricSql {
	return newSQLableArray(array)
}

func newSQLableArray(array any) tomasql.ParametricSql {
	return &sqlableArray{array: pq.Array(array), params: tomasql.ParamsMap{}}
}

func (s *sqlableArray) SqlWithParams(params tomasql.ParamsMap) (string, tomasql.ParamsMap) {
	s.params = params.AddAll(s.params)

	// Use a pointer to the slice as a key to ensure it's hashable. Note: we'll have to pass the array
	// multiple times in the query params even if it's the same array, but it should be fine for most use cases.
	arrayKey := &s.array
	if _, ok := s.params[arrayKey]; !ok {
		s.params[arrayKey] = len(s.params) + 1
	}
	return fmt.Sprintf("(%s)", tomasql.GetDialect().Placeholder(s.params[arrayKey])), s.params
}
