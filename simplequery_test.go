package click

import (
	"testing"
	"time"
)

func TestSimpleQuery_BuildString(t *testing.T) {
	q := SimpleQuery{
		IsTimeSeriesQuery:   true,
		TimeColumn:          "ts",
		GranularityFunction: "toStartOfDay",
		StartTime:           time.Unix(1704038400, 0),
		EndTime:             time.Unix(1706716800, 0),
		Select:              []Expression{Count()},
		From:                "tbl",
		Where:               LiteralExpression(1),
		GroupBy:             []Expression{Column("b")},
		OrderBy:             []Expression{Column("c")},
		Having:              LiteralExpression("d"),
		Limit:               0,
		Offset:              0,
	}
	s, err := q.BuildString()
	if err != nil {
		t.Fatal(err)
	}
	if s != "SELECT count(), toStartOfDay(ts) FROM tbl WHERE (1 AND (ts >= 1704038400) AND (ts < 1706716800)) GROUP BY b, toStartOfDay(ts) ORDER BY c, toStartOfDay(ts) HAVING d" {
		t.Fatal(s)
	}
}
