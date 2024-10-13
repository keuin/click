package click

import (
	"errors"
	"time"
)

type SimpleQuery struct {
	IsTimeSeriesQuery   bool
	TimeColumn          Column
	GranularityFunction string
	StartTime           time.Time
	EndTime             time.Time

	Select  []Expression
	From    string // From is table name
	Where   Expression
	GroupBy []Expression
	OrderBy []Expression
	Having  Expression
	Limit   int // Limit is valid only with positive values
	Offset  int
}

func (q SimpleQuery) Build() (SelectQuery, error) {
	b := SelectBuilder{}
	if len(q.Select) == 0 {
		return nil, errors.New("no selects")
	}
	b.Select(q.Select...)
	if q.From == "" {
		return nil, errors.New("no from")
	}
	b.From(Table(q.From))
	if q.Where != nil {
		b.Where(q.Where)
	}
	if len(q.GroupBy) > 0 {
		b.GroupBy(q.GroupBy...)
	}
	if q.Having != nil {
		b.Having(q.Having)
	}
	if len(q.OrderBy) > 0 {
		b.OrderBy(q.OrderBy...)
	}
	if q.Limit > 0 {
		b.Limit(q.Limit)
	}
	if q.Offset > 0 {
		b.Offset(q.Offset)
	}
	if q.IsTimeSeriesQuery {
		// time series query:
		//  - select & group-by & order-by: add time granularity
		//  - where: add time range filter

		// add time granularity
		timeOffset := Fn(q.GranularityFunction, q.TimeColumn)
		b.Select(timeOffset)
		b.GroupBy(timeOffset)
		b.OrderBy(timeOffset)

		// add time range filter
		var wheres []Expression
		if q.Where != nil {
			wheres = append(wheres, q.Where)
		}
		if !q.StartTime.IsZero() {
			wheres = append(wheres, GreaterOrEqualThan(q.TimeColumn, LiteralExpression(q.StartTime)))
		}
		if !q.EndTime.IsZero() {
			wheres = append(wheres, LessThan(q.TimeColumn, LiteralExpression(q.EndTime)))
		}
		if len(wheres) > 0 {
			b.Where(And(wheres...))
		}
	}
	return sealedSelect(b), nil
}

func (q SimpleQuery) BuildString() (string, error) {
	query, err := q.Build()
	if err != nil {
		return "", err
	}
	return query.String(), nil
}
