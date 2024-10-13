package click

type Operator string

func (o Operator) String() string {
	return string(o)
}

const (
	OpAnd Operator = "AND"
	OpOr  Operator = "OR"
	OpNot Operator = "NOT"
)
