package click

type LiteralOperand struct {
	lit string
}

func (o LiteralOperand) String() string {
	return o.lit
}
