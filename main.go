package gogrammer

import (
	"encoding/json"
	"strings"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
)

// GGrammer is the root of the grammer
type GGrammer struct {
	FetchAll bool      `@("*"|"all")`
	Filters  []*Filter `| @@ { "and" @@ }`
}

// Filter represents a filtering unit
type Filter struct {
	Labeled    *LabelExpression `@@`
	Expression *Expression      `| @@`
}

// LabelExpression x
type LabelExpression struct {
	Token      string        `"label" "(" @(Ident|String|RawString) `
	IsNegation bool          `[ (@("!" "=")|"=")`
	Operand    *ComplexValue `@@ ] ")"`
}

// Expression represents an atomic operation
type Expression struct {
	Token    string   `@Ident`
	SubToken string   `[ "." @Ident ]`
	Operand  *Operand `@@`
}

// Operand represents the right hand side of an expression
type Operand struct {
	Numeric *NumericalOperand `@@`
	Literal *LiteralOperand   `| @@`
}

// LiteralOperand literal operand
type LiteralOperand struct {
	IsNegation bool          `["!"]`
	Value      *ComplexValue `"=" @@`
}

// ComplexValue represents an operand in string format
type ComplexValue struct {
	Value *SimpleValue `@@`
	Regex string       `| @(String|RawString)`
}

// NumericalOperand represents a numeric operand, alongside an operator
type NumericalOperand struct {
	Operator string  `@( ">" ["="] | "<" ["="] )`
	Val      float64 `@(Float|Int)`
}

// SimpleValue is a leaf value
type SimpleValue struct {
	String *string  `@Ident`
	Number *float64 `| @(Float|Int)`
}

// Parse parses an input string for a grammer representation
func Parse(input string) (*GGrammer, error) {
	ggrammer := &GGrammer{}
	parser, err := createANewGrammer()
	if err == nil {
		err = parser.ParseString(input, ggrammer)
	}
	return ggrammer, err
}

// ParseToJSON parses an inout string to Json
func ParseToJSON(input string) (string, error) {
	var result = ""
	parser, err := createANewGrammer()
	if err == nil {
		ggrammer := &GGrammer{}
		err = parser.ParseString(input, ggrammer)
		if err == nil {
			var res []byte
			res, err = json.Marshal(ggrammer)
			result = string(res)
		}
	}
	return result, err
}

func createANewGrammer() (*participle.Parser, error) {
	return participle.Build(&GGrammer{}, participle.Map(toLowercase))
}

const ident rune = -2

func toLowercase(token lexer.Token) lexer.Token {
	var result = token
	if token.Type == ident {
		result = lexer.Token{token.Type, strings.ToLower(token.Value), token.Pos}
	}
	return result
}
