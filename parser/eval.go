package parser

import (
	"fmt"
	"math"
)

// Eval walks the AST from its leaves to compute the result.
func Eval(n *Node) (float64, error) {
	if n == nil {
		return 0, fmt.Errorf("ASTがありません")
	}
	if n.Kind == "number" {
		return n.Number()
	}
	right, err := Eval(n.Right)
	if err != nil {
		return 0, err
	}
	if n.Kind == "unary" {
		switch n.Value {
		case "+":
			return right, nil
		case "-":
			return -right, nil
		}
	}
	if n.Kind != "binary" {
		return 0, fmt.Errorf("不正なノード: %s", n.Kind)
	}
	left, err := Eval(n.Left)
	if err != nil {
		return 0, err
	}
	var result float64
	switch n.Value {
	case "+":
		result = left + right
	case "-":
		result = left - right
	case "*":
		result = left * right
	case "/":
		if right == 0 {
			return 0, fmt.Errorf("ゼロで割ることはできません")
		}
		result = left / right
	default:
		return 0, fmt.Errorf("不正な演算子: %s", n.Value)
	}
	if math.IsInf(result, 0) || math.IsNaN(result) {
		return 0, fmt.Errorf("計算結果が数値の範囲を超えました")
	}
	return result, nil
}
