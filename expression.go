package click

// Expression is a ClickHouse SQL expression AST object.
type Expression interface {
	Expression() string
}

type SelectOrderByExpression interface {
	Expression
	SelectExpression
	OrderByExpression
}

// Column is a ClickHouse table column name, without quotations.
type Column string

func (e Column) Expression() string {
	return string(e)
}
