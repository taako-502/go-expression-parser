package parser

import (
	"fmt"
	"strconv"
)

// Node represents a number, unary operator, or binary operator.
// Parentheses do not need their own node: the tree preserves grouping.
type Node struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
	Left  *Node  `json:"left,omitempty"`
	Right *Node  `json:"right,omitempty"`
}

type reader struct {
	tokens []Token
	pos    int
}

// Parse returns both the AST and tokens so callers can inspect each stage.
func Parse(input string) (*Node, []Token, error) {
	tokens, err := Tokenize(input)
	if err != nil {
		return nil, nil, err
	}
	p := reader{tokens: tokens}
	node, err := p.expression()
	if err == nil && p.peek().Kind != "EOF" {
		err = fmt.Errorf("位置%d: 予期しないトークン %q", p.peek().Pos+1, p.peek().Text)
	}
	if err != nil {
		return nil, tokens, err
	}
	return node, tokens, nil
}

func (p *reader) peek() Token { return p.tokens[p.pos] }

// expression = term { ("+" | "-") term }
func (p *reader) expression() (*Node, error) {
	left, err := p.term()
	if err != nil {
		return nil, err
	}
	for p.peek().Kind == "+" || p.peek().Kind == "-" {
		op := p.peek().Kind
		p.pos++
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		left = &Node{Kind: "binary", Value: op, Left: left, Right: right}
	}
	return left, nil
}

// term = unary { ("*" | "/") unary }
// Parsing this level first gives multiplication/division higher precedence.
func (p *reader) term() (*Node, error) {
	left, err := p.unary()
	if err != nil {
		return nil, err
	}
	for p.peek().Kind == "*" || p.peek().Kind == "/" {
		op := p.peek().Kind
		p.pos++
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		left = &Node{Kind: "binary", Value: op, Left: left, Right: right}
	}
	return left, nil
}

// unary = ("+" | "-") unary | primary
func (p *reader) unary() (*Node, error) {
	if p.peek().Kind == "+" || p.peek().Kind == "-" {
		op := p.peek().Kind
		p.pos++
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &Node{Kind: "unary", Value: op, Right: right}, nil
	}
	return p.primary()
}

// primary = number | "(" expression ")"
func (p *reader) primary() (*Node, error) {
	t := p.peek()
	switch t.Kind {
	case "number":
		p.pos++
		return &Node{Kind: "number", Value: t.Text}, nil
	case "(":
		p.pos++
		node, err := p.expression()
		if err != nil {
			return nil, err
		}
		if p.peek().Kind != ")" {
			return nil, fmt.Errorf("位置%d: 閉じ括弧 ) が必要です", p.peek().Pos+1)
		}
		p.pos++
		return node, nil
	default:
		return nil, fmt.Errorf("位置%d: 数値または開き括弧 ( が必要です", t.Pos+1)
	}
}

// Number converts a number node's literal to float64.
func (n *Node) Number() (float64, error) { return strconv.ParseFloat(n.Value, 64) }
