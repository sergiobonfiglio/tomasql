package goql

type builder1 struct {
}

func (b *builder1) Select(first ParametricSql, columns ...ParametricSql) BuilderWithSelect {
	return newBuilderWithSelect(false, first, columns...)
}

func (b *builder1) SelectDistinct(first ParametricSql, column ...ParametricSql) BuilderWithSelect {
	return newBuilderWithSelect(true, first, column...)
}

func (b *builder1) SelectAll() BuilderWithSelect {
	return newBuilderWithSelectAll(false)
}

func (b *builder1) SelectDistinctAll() BuilderWithSelect {
	return newBuilderWithSelectAll(true)
}

func newBuilder() Builder1 {
	return &builder1{}
}
