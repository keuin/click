package click

import (
	"strings"
)

// SelectExpression is an expression in SELECT clause.
// Any Expression that can be selected must implement SelectExpression,
// defining how it will look like when being selected.
// Symmetrically, any Expression cannot be selected MUST NOT implement SelectExpression.
type SelectExpression interface {
	Expression
	SelectExpression() string
}
type asExpression struct {
	Left  Expression
	Right Expression
}

func (e asExpression) Expression() string {
	return e.Right.Expression()
}

func (e asExpression) SelectExpression() string {
	var sb strings.Builder
	sb.WriteString(e.Left.Expression())
	sb.WriteString(" AS ")
	sb.WriteString(e.Right.Expression())
	return sb.String()
}

func As(l Expression, r Expression) SelectExpression {
	if l == nil {
		panic("empty left value in AS operator")
	}
	if r == nil {
		panic("empty right value in AS operator")
	}
	return asExpression{
		Left:  l,
		Right: r,
	}
}
