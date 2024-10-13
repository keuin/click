package click

import (
	"testing"
)

func TestLiteralExpressions(t *testing.T) {
	v := LiteralExpressions([]string{"a", "b", "c"}, true)
	expectedValues := []string{"'a'", "'b'", "'c'"}
	for i := range v {
		if v[i].Expression() != expectedValues[i] {
			t.Fatal("expected: ", expectedValues[i], ", got: ", v[i].Expression())
		}
	}
}
