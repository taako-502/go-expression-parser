// Package parser demonstrates tokenization, recursive descent parsing, and AST evaluation.
package parser

import (
	"fmt"
	"strconv"
)

// Token positions are zero-based byte offsets in the input.
type Token struct {
	Kind string `json:"kind"`
	Text string `json:"text"`
	Pos  int    `json:"pos"`
}

const MaxInputBytes = 1024

// Tokenize accepts decimal numbers, ASCII whitespace, and + - * / ( ).
func Tokenize(input string) ([]Token, error) {
	if len(input) > MaxInputBytes {
		return nil, fmt.Errorf("入力は%dバイト以内にしてください", MaxInputBytes)
	}
	tokens := []Token{}
	for i := 0; i < len(input); {
		c := input[i]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			i++
			continue
		}
		start := i
		if digit(c) || c == '.' {
			for i < len(input) && digit(input[i]) {
				i++
			}
			if i < len(input) && input[i] == '.' {
				i++
				for i < len(input) && digit(input[i]) {
					i++
				}
			}
			literal := input[start:i]
			if _, err := strconv.ParseFloat(literal, 64); err != nil {
				return nil, fmt.Errorf("位置%d: 数値 %q を解釈できません", start+1, literal)
			}
			tokens = append(tokens, Token{"number", literal, start})
			continue
		}
		switch c {
		case '+', '-', '*', '/', '(', ')':
			tokens = append(tokens, Token{string(c), string(c), i})
			i++
		default:
			return nil, fmt.Errorf("位置%d: 使用できない文字です", i+1)
		}
	}
	return append(tokens, Token{"EOF", "", len(input)}), nil
}

func digit(c byte) bool { return c >= '0' && c <= '9' }
