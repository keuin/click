package click

import (
	"reflect"
	"testing"
)

func TestFn(t *testing.T) {
	type args struct {
		name string
		args []Expression
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "function-call",
			args: args{
				name: "avg",
				args: []Expression{LiteralExpression("a")},
			},
			want: "avg(a)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Fn(tt.args.name, tt.args.args...).Expression(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSum(t *testing.T) {
	type args struct {
		v Expression
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "function-sum",
			args: args{
				v: LiteralExpression("a"),
			},
			want: "sum(a)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sum(tt.args.v).Expression(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAvg(t *testing.T) {
	type args struct {
		v Expression
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "function-avg",
			args: args{
				v: LiteralExpression("a"),
			},
			want: "avg(a)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Avg(tt.args.v).Expression(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Avg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCountIf(t *testing.T) {
	type args struct {
		v Expression
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "function-countIf",
			args: args{
				v: LiteralExpression("a"),
			},
			want: "countIf(a)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CountIf(tt.args.v).Expression(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CountIf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIf(t *testing.T) {
	type args struct {
		cond Expression
		v1   Expression
		v2   Expression
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "function-if",
			args: args{
				cond: LiteralExpression("a"),
				v1:   LiteralExpression("b"),
				v2:   LiteralExpression("c"),
			},
			want: "if(a, b, c)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := If(tt.args.cond, tt.args.v1, tt.args.v2).Expression(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("If() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNotNull(t *testing.T) {
	type args struct {
		v Expression
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "function-isNotNull",
			args: args{
				v: LiteralExpression("a"),
			},
			want: "isNotNull(a)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotNull(tt.args.v).Expression(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IsNotNull() = %v, want %v", got, tt.want)
			}
		})
	}
}
