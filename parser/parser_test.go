package parser

import (
	"math"
	"strings"
	"testing"
)

func TestEvaluate(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  float64
	}{
		{"1 + 2 * 3", 7}, {"(1 + 2) * 3", 9},
		{"10 - 3 - 2", 5}, {"20 / 2 / 5", 2},
		{"-4 + 10 / 2", 1}, {"-(2 + 3) * +2", -10},
		{"1--2", 3}, {".5 + 1.25", 1.75}, {"1.", 1},
		{"\t2 *\n 3\r", 6}, {"1 / 4", .25}, {"0", 0},
	} {
		t.Run(tc.input, func(t *testing.T) {
			ast, _, err := Parse(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			got, err := Eval(ast)
			if err != nil || math.Abs(got-tc.want) > 1e-12 {
				t.Fatalf("got %v, %v; want %v", got, err, tc.want)
			}
		})
	}
}

func TestASTPrecedenceAndGrouping(t *testing.T) {
	ast, tokens, err := Parse("1 + 2 * 3")
	if err != nil {
		t.Fatal(err)
	}
	if ast.Value != "+" || ast.Left.Value != "1" || ast.Right.Value != "*" || ast.Right.Left.Value != "2" || ast.Right.Right.Value != "3" {
		t.Fatalf("incorrect precedence: %#v", ast)
	}
	if len(tokens) != 6 || tokens[2].Pos != 4 || tokens[5].Kind != "EOF" {
		t.Fatalf("incorrect tokens: %#v", tokens)
	}
	ast, _, err = Parse("(1 + 2) * 3")
	if err != nil {
		t.Fatal(err)
	}
	if ast.Value != "*" || ast.Left.Value != "+" || ast.Right.Value != "3" {
		t.Fatalf("incorrect grouping: %#v", ast)
	}
}

func TestInvalidExpressions(t *testing.T) {
	for _, input := range []string{"", " ", "1+", "*2", "()", "(1+2", "1+2)", "1 2", "2(3)", "1..2", ".", "abc", "２+３", "1e3", strings.Repeat("1", 400), strings.Repeat("-", MaxInputBytes) + "1"} {
		t.Run(input, func(t *testing.T) {
			if _, _, err := Parse(input); err == nil {
				t.Fatal("expected parse error")
			}
		})
	}
}

func TestEvaluationErrors(t *testing.T) {
	for _, input := range []string{"1/0", "0/0", "1/(2-2)", "1/-0", strings.Repeat("9", 200) + "*" + strings.Repeat("9", 200)} {
		ast, _, err := Parse(input)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := Eval(ast); err == nil {
			t.Fatalf("expected evaluation error: %s", input)
		}
	}
}
