package click

import (
	"errors"
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
	String() string
}

// SelectBuilder implements builder pattern for constructing SELECT SQLs.
// Its zero value is a ready-to-use empty builder.
// It's recommended to use Select as a shortcut.
type SelectBuilder struct {
	selects  []Expression // Expression | SelectExpression
	from     string
	where    Expression
	groupBy  []Expression
	orderBy  []Expression // Expression | OrderByExpression
	having   Expression
	sample   float64
	limit    int
	hasLimit bool
	offset   int
	format   Format
	pretty   bool
}

func (s *SelectBuilder) Select(values ...Expression) *SelectBuilder {
	s.selects = append(s.selects, values...)
	return s
}

func (s *SelectBuilder) From(table string) *SelectBuilder {
	s.from = table
	return s
}

func (s *SelectBuilder) Sampling(v float64) *SelectBuilder {
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
	var p sqlPrinter
	if len(s.selects) == 0 {
		return "", errors.New("no selects")
	}
	p.BeginClause("SELECT", s.pretty)
	for i := range s.selects {
		var v string
		if expr, ok := s.selects[i].(SelectExpression); ok {
			v = expr.SelectExpression()
		} else {
			v = s.selects[i].Expression()
		}
		p.AddClauseArgument(v, s.pretty, i == len(s.selects)-1)
	}
	if s.from != "" {
		p.BeginClause("FROM", s.pretty)
		p.AddClauseArgument(s.from, s.pretty, true)
	}
	if s.sample > 0 {
		if s.from == "" {
			return "", errors.New("SAMPLE is present while FROM is absent")
		}
		p.BeginClause("SAMPLE", s.pretty)
		p.AddClauseArgument(strconv.FormatFloat(s.sample, 'f', -1, 64), s.pretty, true)
	}
	if s.where != nil {
		p.BeginClause("WHERE", s.pretty)
		p.AddClauseArgument(s.where.Expression(), s.pretty, true)
	}
	if len(s.groupBy) > 0 {
		p.BeginClause("GROUP BY", s.pretty)
		for i := range s.groupBy {
			p.AddClauseArgument(s.groupBy[i].Expression(), s.pretty, i == len(s.groupBy)-1)
		}
	}
	if s.having != nil {
		p.BeginClause("HAVING", s.pretty)
		p.AddClauseArgument(s.having.Expression(), s.pretty, true)
	}
	if len(s.orderBy) > 0 {
		p.BeginClause("ORDER BY", s.pretty)
		for i := range s.orderBy {
			var v string
			if expr, ok := s.orderBy[i].(OrderByExpression); ok {
				v = expr.OrderByExpression()
			} else {
				v = s.orderBy[i].Expression()
			}
			p.AddClauseArgument(v, s.pretty, i == len(s.orderBy)-1)
		}
	}
	if s.hasLimit {
		p.BeginClause("LIMIT", s.pretty)
		p.AddClauseArgument(strconv.Itoa(s.limit), s.pretty, true)
	}
	if s.offset > 0 {
		p.BeginClause("OFFSET", s.pretty)
		p.AddClauseArgument(strconv.Itoa(s.offset), s.pretty, true)
	}
	if s.format != "" {
		p.BeginClause("FORMAT", s.pretty)
		p.AddClauseArgument(string(s.format), s.pretty, true)
	}
	return p.String(), nil
}

func (s *SelectBuilder) PrettyPrint(b ...bool) *SelectBuilder {
	if len(b) > 0 {
		s.pretty = b[0]
	} else {
		s.pretty = true
	}
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
	sb strings.Builder
}

func (p *sqlPrinter) BeginClause(name string, newline bool) {
	if newline {
		p.sb.WriteString(name)
		p.sb.WriteByte('\n')
	} else {
		p.sb.WriteByte(' ')
		p.sb.WriteString(name)
		p.sb.WriteByte(' ')
	}
}

func (p *sqlPrinter) AddClauseArgument(v string, newline, lastElem bool) {
	if newline {
		p.sb.WriteString("\t")
		p.sb.WriteString(v)
		if !lastElem {
			p.sb.WriteString(",")
		}
		p.sb.WriteString("\n")
	} else {
		p.sb.WriteString(v)
		if !lastElem {
			p.sb.WriteString(", ")
		}
	}
}

func (p *sqlPrinter) String() string {
	return strings.TrimSpace(p.sb.String())
}

// sealedSelect is complete, valid and unmodifiable SelectBuilder.
type sealedSelect SelectBuilder

func (s sealedSelect) String() string {
	str, _ := (*SelectBuilder)(&s).BuildString()
	return str
}
