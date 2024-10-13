package click

import "strings"

// OrderByExpression is an Expression in ORDER BY clause.
// Any Expression that can be selected may implement OrderByExpression,
// customizing how it will look like when being ordered by.
type OrderByExpression interface {
	Expression
	OrderByExpression() string
}

type OrderDirection int

const (
	OrderDefault OrderDirection = iota
	OrderAscending
	OrderDescending
)

type orderByExpression struct {
	expression     Expression
	orderDirection OrderDirection
}

func (o orderByExpression) Expression() string {
	return o.expression.Expression()
}

func (o orderByExpression) OrderByExpression() string {
	var sb strings.Builder
	sb.WriteString(o.expression.Expression())
	switch o.orderDirection {
	case OrderDefault:
		// default order, add nothing
	case OrderAscending:
		sb.WriteString(" ASC")
	case OrderDescending:
		sb.WriteString(" DESC")
	default:
		panic("invalid order direction")
	}
	return sb.String()
}

func Desc(v Expression) OrderByExpression {
	return orderByExpression{
		expression:     v,
		orderDirection: OrderDescending,
	}
}

func Asc(v Expression) OrderByExpression {
	return orderByExpression{
		expression:     v,
		orderDirection: OrderAscending,
	}
}
