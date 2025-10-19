package tomasql

var _ ComparableParam[int] = Col[int]{}

func (c Col[T]) Eq(other ParametricSql) Condition {
	return newBinaryCondition(c, other, comparerEq)
}

func (c Col[T]) EqParam(other T) Condition {
	return newBinaryParamCondition(c, other, comparerEq)
}

func (c Col[T]) Gt(other ParametricSql) Condition {
	return newBinaryCondition(c, other, comparerGt)
}

func (c Col[T]) GtParam(other T) Condition {
	return newBinaryParamCondition(c, other, comparerGt)
}

func (c Col[T]) Ge(other ParametricSql) Condition {
	return newBinaryCondition(c, other, comparerGe)
}

func (c Col[T]) GeParam(other T) Condition {
	return newBinaryParamCondition(c, other, comparerGe)
}

func (c Col[T]) Lt(other ParametricSql) Condition {
	return newBinaryCondition(c, other, comparerLt)
}

func (c Col[T]) LtParam(other T) Condition {
	return newBinaryParamCondition(c, other, comparerLt)
}

func (c Col[T]) Le(other ParametricSql) Condition {
	return newBinaryCondition(c, other, comparerLe)
}

func (c Col[T]) LeParam(other T) Condition {
	return newBinaryParamCondition(c, other, comparerLe)
}

func (c Col[T]) Like(other ParametricSql) Condition {
	return newBinaryCondition(c, other, comparerLike)
}

func (c Col[T]) LikeParam(other string) Condition {
	return newBinaryParamCondition(c, other, comparerLike)
}

// func (c Col[T]) ILike(other ParametricSql) Condition {
// 	return newBinaryCondition(c, other, comparerILike)
// }

// func (c Col[T]) ILikeParam(other string) Condition {
// 	return newBinaryParamCondition(c, other, comparerILike)
// }

// func (c Col[T]) InArray(array []T) Condition {
// 	return newInArrayCondition(c, array)
// }

func (c Col[T]) In(sqlable ParametricSql) Condition {
	return newInCondition(c, sqlable)
}

func (c Col[T]) EqAny(sqlable ParametricSql) Condition {
	return newAnyArrayCondition(c, comparerEq, sqlable)
}

func (c Col[T]) EqAll(sqlable ParametricSql) Condition {
	return newAllArrayCondition(c, comparerEq, sqlable)
}

func (c Col[T]) GtAny(sqlable ParametricSql) Condition {
	return newAnyArrayCondition(c, comparerGt, sqlable)
}

func (c Col[T]) GtAll(sqlable ParametricSql) Condition {
	return newAllArrayCondition(c, comparerGt, sqlable)
}

func (c Col[T]) GeAny(sqlable ParametricSql) Condition {
	return newAnyArrayCondition(c, comparerGe, sqlable)
}

func (c Col[T]) GeAll(sqlable ParametricSql) Condition {
	return newAllArrayCondition(c, comparerGe, sqlable)
}

func (c Col[T]) LtAny(sqlable ParametricSql) Condition {
	return newAnyArrayCondition(c, comparerLt, sqlable)
}

func (c Col[T]) LtAll(sqlable ParametricSql) Condition {
	return newAllArrayCondition(c, comparerLt, sqlable)
}

func (c Col[T]) LeAny(sqlable ParametricSql) Condition {
	return newAnyArrayCondition(c, comparerLe, sqlable)
}

func (c Col[T]) LeAll(sqlable ParametricSql) Condition {
	return newAllArrayCondition(c, comparerLe, sqlable)
}

func (c Col[T]) IsNull() Condition {
	return newIsCondition(c, comparerNull)
}

func (c Col[T]) IsNotNull() Condition {
	return newIsCondition(c, comparerNotNull)
}
