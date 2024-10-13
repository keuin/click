package click

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func Select(values ...Expression) *SelectBuilder {
	// clone to avoid unexpected modification on the argument
	// not using slices.Clone to keep compatible with 1.18
	selects := make([]Expression, len(values))
	for i := range selects {
		selects[i] = values[i]
	}
	return &SelectBuilder{
		selects: selects,
	}
}

type SelectQuery interface {
	FromExpression
	String() string
}

// SelectBuilder implements builder pattern for constructing SELECT SQLs.
// Its zero value is a ready-to-use empty builder.
// It's recommended to use Select as a shortcut.
type SelectBuilder struct {
	selects  []Expression // Expression | SelectExpression
	from     FromExpression
	where    Expression
	groupBy  []Expression
	orderBy  []Expression // Expression | OrderByExpression
	having   Expression
	sample   float64
	limit    int
	hasLimit bool
	offset   int
	format   Format
	style    RenderStyle
	styleSet bool
}

func (s *SelectBuilder) FromExpression(style RenderStyle) (string, error) {
	style.IndentLevel++
	fromExpr, err := s.buildString(style)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	sb.WriteString("(\n")
	for i := 0; i < style.IndentLevel; i++ {
		sb.WriteString(style.Indent)
	}
	sb.WriteString(fromExpr)
	sb.WriteString("\n")
	style.IndentLevel--
	for i := 0; i < style.IndentLevel; i++ {
		sb.WriteString(style.Indent)
	}
	sb.WriteString(")")
	return sb.String(), nil
}

func (s *SelectBuilder) Select(values ...Expression) *SelectBuilder {
	s.selects = append(s.selects, values...)
	return s
}

func (s *SelectBuilder) From(table FromExpression) *SelectBuilder {
	s.from = table
	return s
}

func (s *SelectBuilder) Sample(v float64) *SelectBuilder {
	s.sample = v
	return s
}

func (s *SelectBuilder) Where(where Expression) *SelectBuilder {
	s.where = where
	return s
}

func (s *SelectBuilder) GroupBy(values ...Expression) *SelectBuilder {
	s.groupBy = append(s.groupBy, values...)
	return s
}

func (s *SelectBuilder) OrderBy(values ...Expression) *SelectBuilder {
	s.orderBy = append(s.orderBy, values...)
	return s
}

func (s *SelectBuilder) Having(value Expression) *SelectBuilder {
	s.having = value
	return s
}

func (s *SelectBuilder) Limit(n int) *SelectBuilder {
	s.limit = n
	s.hasLimit = true
	return s
}

func (s *SelectBuilder) Offset(n int) *SelectBuilder {
	s.offset = n
	return s
}

func (s *SelectBuilder) Format(f Format) *SelectBuilder {
	s.format = f
	return s
}

func (s *SelectBuilder) BuildString() (string, error) {
	style := defaultStyle
	if s.styleSet {
		style = s.style
	}
	return s.buildString(style)
}

// buildString ignores style settings in SelectBuilder itself, using the RenderStyle in argument.
// This is used in nested query rendering. Since RenderStyle is only needed when rendering the entire SELECT expression,
// other Expression interfaces (like Expression, SelectExpression, OrderByExpression) does not accept RenderStyle as argument.
// buildString is a simple replacement for a common and public interface which accepts RenderStyle.
// This may be exposed in the future, or move RenderStyle in sqlPrinter, and pass sqlPrinter in the entire process.
func (s *SelectBuilder) buildString(style RenderStyle) (string, error) {
	p := sqlPrinter{
		Style: style,
	}
	if len(s.selects) == 0 {
		return "", errors.New("no selects")
	}
	p.BeginClause("SELECT")
	for i := range s.selects {
		var v string
		if expr, ok := s.selects[i].(SelectExpression); ok {
			v = expr.SelectExpression()
		} else {
			v = s.selects[i].Expression()
		}
		p.AddClauseArgument(v, i == len(s.selects)-1)
	}
	if s.from != nil {
		p.BeginClause("FROM")
		fromExpr, err := s.from.FromExpression(style)
		if err != nil {
			return "", fmt.Errorf("build FROM clause: %w", err)
		}
		_, isTable := s.from.(Table)
		p.AddClauseArgumentPrefix(fromExpr, true, isTable)
	}
	if s.sample > 0 {
		if s.from == nil {
			return "", errors.New("SAMPLE is present while FROM is absent")
		}
		p.BeginClause("SAMPLE")
		p.AddClauseArgument(strconv.FormatFloat(s.sample, 'f', -1, 64), true)
	}
	if s.where != nil {
		p.BeginClause("WHERE")
		p.AddClauseArgument(s.where.Expression(), true)
	}
	if len(s.groupBy) > 0 {
		p.BeginClause("GROUP BY")
		for i := range s.groupBy {
			p.AddClauseArgument(s.groupBy[i].Expression(), i == len(s.groupBy)-1)
		}
	}
	if s.having != nil {
		p.BeginClause("HAVING")
		p.AddClauseArgument(s.having.Expression(), true)
	}
	if len(s.orderBy) > 0 {
		p.BeginClause("ORDER BY")
		for i := range s.orderBy {
			var v string
			if expr, ok := s.orderBy[i].(OrderByExpression); ok {
				v = expr.OrderByExpression()
			} else {
				v = s.orderBy[i].Expression()
			}
			p.AddClauseArgument(v, i == len(s.orderBy)-1)
		}
	}
	if s.hasLimit {
		p.BeginClause("LIMIT")
		p.AddClauseArgument(strconv.Itoa(s.limit), true)
	}
	if s.offset > 0 {
		p.BeginClause("OFFSET")
		p.AddClauseArgument(strconv.Itoa(s.offset), true)
	}
	if s.format != "" {
		p.BeginClause("FORMAT")
		p.AddClauseArgument(string(s.format), true)
	}
	return p.String(), nil
}

func (s *SelectBuilder) PrettyPrint(b ...bool) *SelectBuilder {
	if len(b) == 0 || b[0] {
		s.style = prettyStyle
	} else {
		s.style = defaultStyle
	}
	s.styleSet = true
	return s
}

func (s *SelectBuilder) Build() (SelectQuery, error) {
	_, err := s.BuildString()
	if err != nil {
		return nil, err
	}
	return (*sealedSelect)(s), nil
}

type sqlPrinter struct {
	sb    strings.Builder
	Style RenderStyle
}

func (p *sqlPrinter) BeginClause(name string) {
	for i := 0; i < p.Style.IndentLevel; i++ {
		p.sb.WriteString(p.Style.Indent)
	}
	p.sb.WriteString(p.Style.ClauseNamePrefix)
	p.sb.WriteString(name)
	p.sb.WriteString(p.Style.ClauseNameSuffix)
}

func (p *sqlPrinter) AddClauseArgument(v string, lastElem bool) {
	p.AddClauseArgumentPrefix(v, lastElem, true)
}

func (p *sqlPrinter) AddClauseArgumentPrefix(v string, lastElem, prefix bool) {
	for i := 0; i < p.Style.IndentLevel; i++ {
		p.sb.WriteString(p.Style.Indent)
	}
	if prefix {
		p.sb.WriteString(p.Style.ArgumentPrefix)
	}
	p.sb.WriteString(v)
	if !lastElem {
		p.sb.WriteString(p.Style.ArgumentDelimiter)
	}
	p.sb.WriteString(p.Style.ArgumentSuffix)
}

func (p *sqlPrinter) String() string {
	return strings.TrimSpace(p.sb.String())
}

// sealedSelect is complete, valid and unmodifiable SelectBuilder.
type sealedSelect SelectBuilder

func (s sealedSelect) FromExpression(style RenderStyle) (string, error) {
	return (*SelectBuilder)(&s).FromExpression(style)
}

func (s sealedSelect) String() string {
	str := must((*SelectBuilder)(&s).BuildString())
	return str
}
