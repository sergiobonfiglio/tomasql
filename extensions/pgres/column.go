package pgres

import "github.com/sergiobonfiglio/tomasql"

type PGCol[T any] struct {
	*tomasql.Col[T]
}

var (
	_ tomasql.Column        = &PGCol[any]{}
	_ tomasql.ParametricSql = &PGCol[int]{} // Ensure Col implements ParametricSql
)

func Wrap[T any](c *tomasql.Col[T]) *PGCol[T] {
	return &PGCol[T]{Col: c}
}

const comparerILike = "ILIKE" // case-insensitive LIKE

func (c PGCol[T]) ILike(other tomasql.ParametricSql) tomasql.Condition {
	return tomasql.NewBinaryCondition(c, other, comparerILike)
}

func (c PGCol[T]) ILikeParam(other string) tomasql.Condition {
	return tomasql.NewBinaryParamCondition(c, other, comparerILike)
}

// func (f *funcCol[T]) ILike(other ParametricSql) Condition {
// 	return NewBinaryCondition(f, other, comparerILike)
// }

// func (f *funcCol[T]) ILikeParam(pattern string) Condition {
// 	return NewBinaryParamCondition(f, pattern, comparerILike)
// }

// func (f *funcCol[T]) InArray(array []T) Condition {
// 	return newInArrayCondition(f, array)
// }

func (c PGCol[T]) InArray(array []T) tomasql.Condition {
	return newInArrayCondition(c, array)
}
