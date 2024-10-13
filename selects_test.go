package click

import (
	"testing"
)

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func TestSelect_Simple(t *testing.T) {
	s := Select(LiteralExpression(1))
	v := must(s.BuildString())
	if v != "SELECT 1" {
		t.Fail()
	}
}

func TestSelect_SimpleMultiple(t *testing.T) {
	s := Select(LiteralExpression(1), LiteralExpression(2), LiteralExpression(3))
	v := must(s.BuildString())
	if v != "SELECT 1, 2, 3" {
		t.Fail()
	}
}

func TestSelect_Alias(t *testing.T) {
	s := Select(As(LiteralExpression(1), LiteralExpression("a")), LiteralExpression(2), LiteralExpression(3))
	v := must(s.BuildString())
	if v != "SELECT 1 AS a, 2, 3" {
		t.Fail()
	}
}

func TestSelect_Aggregate(t *testing.T) {
	s := Select(Avg(Column("score"))).
		From("tbl").
		Where(And(
			GreaterOrEqualThan(Column("date"), LiteralExpressionQuoted("2025-01-01")),
			LessThan(Column("date"), LiteralExpressionQuoted("2025-02-01")),
		)).
		GroupBy(Column("date"))
	v := must(s.BuildString())
	if v != "SELECT avg(score) FROM tbl WHERE ((date >= '2025-01-01') AND (date < '2025-02-01')) GROUP BY date" {
		t.Fatal(v)
	}
}

func TestSelect_AliasReference(t *testing.T) {
	avgScore := As(Avg(Column("score")), LiteralExpression("avg_score"))
	s := Select(avgScore).
		From("tbl").
		Sampling(0.1).
		Where(And(
			GreaterOrEqualThan(Column("date"), LiteralExpressionQuoted("2025-01-01")),
			LessThan(Column("date"), LiteralExpressionQuoted("2025-02-01")),
		)).
		GroupBy(Column("date")).
		OrderBy(Column("date"), Asc(Column("date")), Desc(avgScore)).
		Having(GreaterThan(avgScore, LiteralExpression(60))).
		Limit(5).Offset(10).
		Format(FormatCSV)
	v := must(s.BuildString())
	if v != "SELECT avg(score) AS avg_score FROM tbl SAMPLE 0.1 WHERE ((date >= '2025-01-01') AND (date < '2025-02-01')) GROUP BY date HAVING (avg_score > 60) ORDER BY date, date ASC, avg_score DESC LIMIT 5 OFFSET 10 FORMAT CSV" {
		t.Fatal(v)
	}
}
