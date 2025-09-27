package goql

type builder1 struct{}

func (b *builder1) SelectCols(first Column, columns ...Column) BuilderWithSelect {
	convertedColumns := make([]ParametricSql, len(columns))
	for i, col := range columns {
		convertedColumns[i] = col
	}

	return newBuilderWithSelect(false, first, convertedColumns...)
}

func (b *builder1) Select(first ParametricSql, columns ...ParametricSql) BuilderWithSelect {
	return newBuilderWithSelect(false, first, columns...)
}

func (b *builder1) SelectDistinctCols(first Column, columns ...Column) BuilderWithSelect {
	convertedColumns := make([]ParametricSql, len(columns))
	for i, col := range columns {
		convertedColumns[i] = col
	}
	return newBuilderWithSelect(true, first, convertedColumns...)
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

func NewBuilder() Builder1 {
	return &builder1{}
}
