package click

import (
	"strings"
)

func And(sub ...Expression) Expression {
	return Concatenate(OpAnd, sub...)
}

func Or(sub ...Expression) Expression {
	return Concatenate(OpOr, sub...)
}

func Concatenate(op Operator, sub ...Expression) Expression {
	if len(sub) == 0 {
		panic("empty subexpressions")
	}
	for len(sub) == 1 {
		return sub[0]
	}
	return concatenatedExpression{
		Op:   op,
		Expr: sub,
	}
}

type concatenatedExpression struct {
	Op   Operator
	Expr []Expression // 长度必须大于1
}

func (c concatenatedExpression) Expression() string {
	sb := strings.Builder{}
	sb.WriteByte('(')
	for i, ex := range c.Expr {
		if i != 0 {
			sb.WriteByte(' ')
			sb.WriteString(string(c.Op))
			sb.WriteByte(' ')
		}
		sb.WriteString(ex.Expression())
	}
	sb.WriteByte(')')
	return sb.String()
}

func (c concatenatedExpression) SelectExpression() Expression {
	return c
}

func Equal(l Expression, r Expression) Expression {
	return BinaryExpression{
		Operator:     "=",
		LeftOperand:  l,
		RightOperand: r,
	}
}

func GreaterThan(l Expression, r Expression) Expression {
	return BinaryExpression{
		Operator:     ">",
		LeftOperand:  l,
		RightOperand: r,
	}
}

func GreaterOrEqualThan(l Expression, r Expression) Expression {
	return BinaryExpression{
		Operator:     ">=",
		LeftOperand:  l,
		RightOperand: r,
	}
}

func LessOrEqualThan(l Expression, r Expression) Expression {
	return BinaryExpression{
		Operator:     "<=",
		LeftOperand:  l,
		RightOperand: r,
	}
}

func LessThan(l Expression, r Expression) Expression {
	return BinaryExpression{
		Operator:     "<",
		LeftOperand:  l,
		RightOperand: r,
	}
}

const errEmptyTuple = "ClickHouse tuple must have at least one element"

// Tuple is ClickHouse tuple object, and must contain at least one element.
// See https://clickhouse.com/docs/sql-reference/data-types/tuple
type Tuple []Expression

func (t Tuple) mustValid() {
	if len(t) == 0 {
		panic(errEmptyTuple)
	}
}

func (t Tuple) SelectExpression() string {
	t.mustValid()
	return t.Expression()
}

func (t Tuple) Expression() string {
	t.mustValid()
	var sb strings.Builder
	sb.WriteByte('(')
	for i, expr := range t {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(expr.Expression())
	}
	sb.WriteByte(')')
	return sb.String()
}

func In(v Expression, ary Tuple) Expression {
	return BinaryExpression{
		Operator:     "IN",
		LeftOperand:  v,
		RightOperand: ary,
	}
}

func NotIn(v Expression, ary Tuple) Expression {
	return BinaryExpression{
		Operator:     "NOT IN",
		LeftOperand:  v,
		RightOperand: ary,
	}
}

func NotEqual(l Expression, r Expression) Expression {
	return BinaryExpression{
		Operator:     "!=",
		LeftOperand:  l,
		RightOperand: r,
	}
}
