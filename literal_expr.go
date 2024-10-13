package click

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type BinaryExpression struct {
	Operator     Operator
	LeftOperand  Expression
	RightOperand Expression
}

func (e BinaryExpression) Expression() string {
	var sb strings.Builder
	sb.WriteByte('(')
	sb.WriteString(e.LeftOperand.Expression())
	sb.WriteByte(' ')
	sb.WriteString(e.Operator.String())
	sb.WriteByte(' ')
	sb.WriteString(e.RightOperand.Expression())
	sb.WriteByte(')')
	return sb.String()
}

// LiteralExpression converts a Go value to a SQL string, formatting that value to string with fmt.Sprint.
func LiteralExpression[T any](v T) Expression {
	if vv, ok := interface{}(v).(Expression); ok {
		return vv
	}
	return literalExpr[T]{val: v}
}

// LiteralExpressionQuoted creates literal expression, treating the argument as a string value, not a SQL expression
func LiteralExpressionQuoted[T any](v T) Expression {
	return literalExpr[T]{val: v, quoteString: true}
}

func LiteralExpressions[T any](v []T, quoteString bool) (ret []Expression) {
	ret = make([]Expression, len(v))
	for i := range ret {
		ret[i] = literalExpr[T]{val: v[i], quoteString: quoteString}
	}
	return ret
}

type alias string

func (a alias) Expression() string {
	return string(a)
}

func Alias(s string) Expression {
	return alias(s)
}

type literalExpr[T any] struct {
	val         T
	quoteString bool
}

func (e literalExpr[T]) Expression() string {
	if typ := reflect.TypeOf(e.val); typ != nil && typ.Kind() == reflect.String {
		// fast access to the underneath string value
		s := *convInto[T, string](&e.val)
		if e.quoteString {
			return "'" + clickhouseStringEscapeReplacer.Replace(s) + "'"
		}
		return s
	}
	if ts, ok := interface{}(e.val).(interface{ Unix() int64 }); ok {
		// compatible with both time.Time and *time.Time
		// ClickHouse DateTime has 1 second resolution, so Unix() is enough.
		return strconv.FormatInt(ts.Unix(), 10)
	}
	return fmt.Sprint(e.val)
}

var clickhouseStringEscapeReplacer = strings.NewReplacer(
	`'`, `\'`,
	`\`, `\\`,
)
