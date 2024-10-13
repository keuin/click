package click

import (
	"strings"
)

// popular SQL functions

func Fn(name string, args ...Expression) Expression {
	return fnCall{
		name: name,
		args: args,
	}
}

type fnCall struct {
	name string
	args []Expression
}

func (f fnCall) Expression() string {
	var sb strings.Builder
	sb.WriteString(f.name)
	sb.WriteByte('(')
	for i := range f.args {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(f.args[i].Expression())
	}
	sb.WriteByte(')')
	return sb.String()
}

func Sum(v Expression) Expression           { return Fn("sum", v) }
func Avg(v Expression) Expression           { return Fn("avg", v) }
func CountIf(v Expression) Expression       { return Fn("countIf", v) }
func If(cond, v1, v2 Expression) Expression { return Fn("if", cond, v1, v2) }
func IsNotNull(v Expression) Expression     { return Fn("isNotNull", v) }
