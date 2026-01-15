package tomasql

var _ ComparableParam[int] = Col[int]{}

func (c Col[T]) Eq(other ParametricSql) Condition {
	return NewBinaryCondition(c, other, comparerEq)
}

func (c Col[T]) EqParam(other T) Condition {
	return NewBinaryParamCondition(c, other, comparerEq)
}

func (c Col[T]) Neq(other ParametricSql) Condition {
	return NewBinaryCondition(c, other, comparerNeq)
}

func (c Col[T]) NeqParam(other T) Condition {
	return NewBinaryParamCondition(c, other, comparerNeq)
}

func (c Col[T]) Gt(other ParametricSql) Condition {
	return NewBinaryCondition(c, other, comparerGt)
}

func (c Col[T]) GtParam(other T) Condition {
	return NewBinaryParamCondition(c, other, comparerGt)
}

func (c Col[T]) Ge(other ParametricSql) Condition {
	return NewBinaryCondition(c, other, comparerGe)
}

func (c Col[T]) GeParam(other T) Condition {
	return NewBinaryParamCondition(c, other, comparerGe)
}

func (c Col[T]) Lt(other ParametricSql) Condition {
	return NewBinaryCondition(c, other, comparerLt)
}

func (c Col[T]) LtParam(other T) Condition {
	return NewBinaryParamCondition(c, other, comparerLt)
}

func (c Col[T]) Le(other ParametricSql) Condition {
	return NewBinaryCondition(c, other, comparerLe)
}

func (c Col[T]) LeParam(other T) Condition {
	return NewBinaryParamCondition(c, other, comparerLe)
}

func (c Col[T]) Like(other ParametricSql) Condition {
	return NewBinaryCondition(c, other, comparerLike)
}

func (c Col[T]) LikeParam(other string) Condition {
	return NewBinaryParamCondition(c, other, comparerLike)
}

func (c Col[T]) In(sqlable ParametricSql) Condition {
	return newInCondition(c, sqlable)
}

func (c Col[T]) EqAny(sqlable ParametricSql) Condition {
	return newAnyCondition(c, comparerEq, sqlable)
}

func (c Col[T]) EqAll(sqlable ParametricSql) Condition {
	return newAllCondition(c, comparerEq, sqlable)
}

func (c Col[T]) GtAny(sqlable ParametricSql) Condition {
	return newAnyCondition(c, comparerGt, sqlable)
}

func (c Col[T]) GtAll(sqlable ParametricSql) Condition {
	return newAllCondition(c, comparerGt, sqlable)
}

func (c Col[T]) GeAny(sqlable ParametricSql) Condition {
	return newAnyCondition(c, comparerGe, sqlable)
}

func (c Col[T]) GeAll(sqlable ParametricSql) Condition {
	return newAllCondition(c, comparerGe, sqlable)
}

func (c Col[T]) LtAny(sqlable ParametricSql) Condition {
	return newAnyCondition(c, comparerLt, sqlable)
}

func (c Col[T]) LtAll(sqlable ParametricSql) Condition {
	return newAllCondition(c, comparerLt, sqlable)
}

func (c Col[T]) LeAny(sqlable ParametricSql) Condition {
	return newAnyCondition(c, comparerLe, sqlable)
}

func (c Col[T]) LeAll(sqlable ParametricSql) Condition {
	return newAllCondition(c, comparerLe, sqlable)
}

func (c Col[T]) IsNull() Condition {
	return newIsCondition(c, comparerNull)
}

func (c Col[T]) IsNotNull() Condition {
	return newIsCondition(c, comparerNotNull)
}
