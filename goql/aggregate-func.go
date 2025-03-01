package goql

type AggregateFunc[T any] interface {
	As(string) SQLable
	ParametricSql
	Comparable[T]
	SetComparable[T]
}

func Count() AggregateFunc[int] {
	return newCol[int]("COUNT(1)", nil)
}
