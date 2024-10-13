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
		t.Fatal(v)
	}
}

func TestSelect_Simple_Pretty(t *testing.T) {
	s := Select(LiteralExpression(1))
	v := must(s.PrettyPrint().BuildString())
	if v != "SELECT\n\t1" {
		t.Fatal(v)
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

func TestSelect_Alias_Pretty(t *testing.T) {
	s := Select(As(LiteralExpression(1), LiteralExpression("a")), LiteralExpression(2), LiteralExpression(3))
	v := must(s.PrettyPrint().BuildString())
	if v != "SELECT\n\t1 AS a,\n\t2,\n\t3" {
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

func TestSelect_AliasReference_Pretty(t *testing.T) {
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
	v := must(s.PrettyPrint().BuildString())
	if v != "SELECT\n\tavg(score) AS avg_score\nFROM\n\ttbl\nSAMPLING\n\t0.1\nWHERE\n\t((date >= '2025-01-01') AND (date < '2025-02-01'))\nGROUP BY\n\tdate\nORDER BY\n\tdate,\n\tdate ASC,\n\tavg_score DESC\nHAVING\n\t(avg_score > 60)\nLIMIT\n\t5\nOFFSET\n\t10\nFORMAT\n\tCSV" {
		t.Fatal(v)
	}
	t.Log(v)
}
