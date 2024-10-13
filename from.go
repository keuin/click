package click

type Table string

func (t Table) FromExpression(_ RenderStyle) (string, error) {
	return string(t), nil
}

// FromExpression is an Expression in FROM clause.
// Any Expression that can be used in nested query may implement FromExpression,
// customizing how it will look like when being selected from.
type FromExpression interface {
	FromExpression(style RenderStyle) (string, error)
}
