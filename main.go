package gogrammer

import (
	"encoding/json"

	"github.com/alecthomas/participle"
)

// GGrammer is the root of the grammer
type GGrammer struct {
	FetchAll bool      `@("*"|"all"|"alL"|"aLL"|"aLl"|"All"|"AlL"|"ALl"|"ALL")`
	Filters  []*Filter `| @@ { ("AND"|"ANd"|"AnD"|"And"|"anD"|"aNd"|"aND"|"and") @@ }`
}

// Filter represents a filtering unit
type Filter struct {
	Labels     []*Value    `"label" "(" @@ [ "=" @@ ] ")"`
	Expression *Expression `| @@`
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
	Literal *LiteralOperand   `| "=" @@`
}

// LiteralOperand represents an operand in string format
type LiteralOperand struct {
	Equivalent *Value `@@`
	Regex      string `| @(String|RawString)`
}

// NumericalOperand represents a numeric operand, alongside an operator
type NumericalOperand struct {
	Operator string  `@(">"|">="|"<"|"<=")`
	Val      float64 `@(Float|Int)`
}

// Value is a leaf value
type Value struct {
	String *string  `@Ident`
	Number *float64 `| @(Float|Int)`
}

// Parse parses an input string for a grammer representation
func Parse(input string) (*GGrammer, error) {
	ggrammer := &GGrammer{}
	parser, err := participle.Build(&GGrammer{})
	if err == nil {
		err = parser.ParseString(input, ggrammer)
	}
	return ggrammer, err
}

// ParseToJSON parses an inout string to Json
func ParseToJSON(input string) (string, error) {
	var result = ""
	parser, err := participle.Build(&GGrammer{})
	if err == nil {
		ggrammer := &GGrammer{}
		err = parser.ParseString(input, ggrammer)
		if err == nil {
			res, _ := json.Marshal(ggrammer)
			result = string(res)
		}
	}
	return result, err
}
