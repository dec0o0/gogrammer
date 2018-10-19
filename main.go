package main

import (
	"encoding/json"
	"os"

	"github.com/alecthomas/participle"
)

type GGrammer struct {
	FetchAll bool      `@"*"`
	Filters  []*Filter `| @@ { ("AND"|"ANd"|"AnD"|"And"|"anD"|"aNd"|"aND"|"and") @@ }`
}

type Filter struct {
	Labels     []*Value    `"label" "(" @@ [ "=" @@ ] ")"`
	Expression *Expression `| @@`
}

type Expression struct {
	Token    string   `@Ident`
	SubToken string   `[ "." @Ident ]`
	Operand  *Operand `@@`
}

type Operand struct {
	Numeric *NumericalOperand `@@`
	Literal *LiteralOperand   `| "=" @@`
}

type LiteralOperand struct {
	Equivalent *Value `@@`
	Regex      string `| @(String|RawString)`
}

type NumericalOperand struct {
	Operator string  `@(">"|">="|"<"|"<=")`
	Val      float64 `@(Float|Int)`
}

type Value struct {
	String string  `@Ident`
	Number float64 `| @(Float|Int)`
}

func main() {
	parser, err := participle.Build(&GGrammer{})
	if err != nil {
		panic(err)
	}
	ggrammer := &GGrammer{}
	err = parser.ParseString("table.column = val123 aNd simple = complex AND song = 'alelu[a|b|c]{2,}ia.+' and label(malicious) and label  ( status = malicious1 ) and score > 5", ggrammer)
	if err != nil {
		panic(err)
	}

	json.NewEncoder(os.Stdout).Encode(ggrammer)
}
