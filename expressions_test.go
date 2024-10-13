package click

import (
	"reflect"
	"testing"
)

func TestAnd(t *testing.T) {
	type args struct {
		sub []Expression
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "and 1 element",
			args: args{
				sub: []Expression{LiteralExpression(true)},
			},
			want: "true",
		},
		{
			name: "and 2 elements",
			args: args{
				sub: []Expression{LiteralExpression(true), LiteralExpression(false)},
			},
			want: "(true AND false)",
		},
		{
			name: "and 3 elements",
			args: args{
				sub: []Expression{
					LiteralExpression(true),
					LiteralExpression(false),
					LiteralExpression(false),
				},
			},
			want: "(true AND false AND false)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := And(tt.args.sub...).Expression(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("And() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOr(t *testing.T) {
	type args struct {
		sub []Expression
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "or 1 element",
			args: args{
				sub: []Expression{LiteralExpression(true)},
			},
			want: "true",
		},
		{
			name: "or 2 elements",
			args: args{
				sub: []Expression{LiteralExpression(true), LiteralExpression(false)},
			},
			want: "(true OR false)",
		},
		{
			name: "or 3 elements",
			args: args{
				sub: []Expression{
					LiteralExpression(true),
					LiteralExpression(false),
					LiteralExpression(false),
				},
			},
			want: "(true OR false OR false)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Or(tt.args.sub...).Expression(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Or() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	type args struct {
		l Expression
		r Expression
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test equal",
			args: args{
				l: LiteralExpression("a"),
				r: LiteralExpression("b"),
			},
			want: "(a = b)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Equal(tt.args.l, tt.args.r).Expression(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGreaterThan(t *testing.T) {
	type args struct {
		l Expression
		r Expression
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test greater than",
			args: args{
				l: LiteralExpression("a"),
				r: LiteralExpression("b"),
			},
			want: "(a > b)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GreaterThan(tt.args.l, tt.args.r).Expression(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GreaterThan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGreaterOrEqualThan(t *testing.T) {
	type args struct {
		l Expression
		r Expression
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test greater/equal than",
			args: args{
				l: LiteralExpression("a"),
				r: LiteralExpression("b"),
			},
			want: "(a >= b)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GreaterOrEqualThan(tt.args.l, tt.args.r).Expression(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GreaterOrEqualThan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLessOrEqualThan(t *testing.T) {
	type args struct {
		l Expression
		r Expression
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test less/equal than",
			args: args{
				l: LiteralExpression("a"),
				r: LiteralExpression("b"),
			},
			want: "(a <= b)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LessOrEqualThan(tt.args.l, tt.args.r).Expression(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LessOrEqualThan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLessThan(t *testing.T) {
	type args struct {
		l Expression
		r Expression
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test less than",
			args: args{
				l: LiteralExpression("a"),
				r: LiteralExpression("b"),
			},
			want: "(a < b)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LessThan(tt.args.l, tt.args.r).Expression(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LessThan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTuple_Expression(t *testing.T) {
	tests := []struct {
		name string
		t    Tuple
		want string
	}{
		{
			name: "single-element tuple",
			t:    Tuple{LiteralExpression(1)},
			want: "(1)",
		},
		{
			name: "2-elements tuple",
			t:    Tuple{LiteralExpression(1), LiteralExpression(2)},
			want: "(1, 2)",
		},
		{
			name: "3-elements tuple",
			t:    Tuple{LiteralExpression(1), LiteralExpression(2), LiteralExpression(3)},
			want: "(1, 2, 3)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.Expression(); got != tt.want {
				t.Errorf("Expression() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIn(t *testing.T) {
	type args struct {
		v   Expression
		ary Tuple
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "simple in",
			args: args{
				v:   LiteralExpression("a"),
				ary: Tuple{LiteralExpression(1), LiteralExpression(2)},
			},
			want: "(a IN (1, 2))",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := In(tt.args.v, tt.args.ary).Expression(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("In() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotIn(t *testing.T) {
	type args struct {
		v   Expression
		ary Tuple
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "simple not-in",
			args: args{
				v:   LiteralExpression("a"),
				ary: Tuple{LiteralExpression(1), LiteralExpression(2)},
			},
			want: "(a NOT IN (1, 2))",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NotIn(tt.args.v, tt.args.ary).Expression(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NotIn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotEqual(t *testing.T) {
	type args struct {
		l Expression
		r Expression
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test not-equal",
			args: args{
				l: LiteralExpression("a"),
				r: LiteralExpression("b"),
			},
			want: "(a != b)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NotEqual(tt.args.l, tt.args.r).Expression(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NotEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConcatenate_EmptySubExpression(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic, got nothing")
		}
	}()
	Concatenate(OpAnd)
}

func TestTuple_Expression_EmptyTuple(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic, got nothing")
		}
	}()
	Tuple{}.Expression()
}
