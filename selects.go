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
	var sb strings.Builder
	if len(s.selects) == 0 {
		return "", errors.New("no selects")
	}
	sb.WriteString("SELECT ")
	for i := range s.selects {
		if i != 0 {
			sb.WriteString(", ")
		}
		if expr, ok := s.selects[i].(SelectExpression); ok {
			sb.WriteString(expr.SelectExpression())
		} else {
			sb.WriteString(s.selects[i].Expression())
		}
	}
	if s.from != "" {
		sb.WriteString(" FROM ")
		sb.WriteString(s.from)
	}
	if s.sample > 0 {
		if s.from == "" {
			return "", errors.New("SAMPLE is present while FROM is absent")
		}
		sb.WriteString(" SAMPLE ")
		sb.WriteString(strconv.FormatFloat(s.sample, 'f', -1, 64))
	}
	if s.where != nil {
		sb.WriteString(" WHERE ")
		sb.WriteString(s.where.Expression())
	}
	if len(s.groupBy) > 0 {
		sb.WriteString(" GROUP BY ")
		for i := range s.groupBy {
			if i != 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(s.groupBy[i].Expression())
		}
	}
	if s.having != nil {
		sb.WriteString(" HAVING ")
		sb.WriteString(s.having.Expression())
	}
	if len(s.orderBy) > 0 {
		sb.WriteString(" ORDER BY ")
		for i := range s.orderBy {
			if i != 0 {
				sb.WriteString(", ")
			}
			if expr, ok := s.orderBy[i].(OrderByExpression); ok {
				sb.WriteString(expr.OrderByExpression())
			} else {
				sb.WriteString(s.orderBy[i].Expression())
			}
		}
	}
	if s.hasLimit {
		sb.WriteString(" LIMIT ")
		sb.WriteString(strconv.Itoa(s.limit))
	}
	if s.offset > 0 {
		sb.WriteString(" OFFSET ")
		sb.WriteString(strconv.Itoa(s.offset))
	}
	if s.format != "" {
		sb.WriteString(" FORMAT ")
		sb.WriteString(string(s.format))
	}
	return sb.String(), nil
}

func (s *SelectBuilder) Build() (SelectQuery, error) {
	_, err := s.BuildString()
	if err != nil {
		return nil, err
	}
	return (*sealedSelect)(s), nil
}

// sealedSelect is complete, valid and unmodifiable SelectBuilder.
type sealedSelect SelectBuilder

func (s sealedSelect) String() string {
	str, _ := (*SelectBuilder)(&s).BuildString()
	return str
}
