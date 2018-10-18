package main

import (
	"encoding/json"
	"os"

	"github.com/alecthomas/participle"
)

type GGrammer struct {
	All         bool          `@"*"`
	Expressions []*Expression `| @@ { ("AND"|"and") @@ }`
}

type Expression struct {
	LabelEquivalence *Equivalence `"label" "(" @@ ")"`
	Equivalence      *Equivalence `| @@`
}

type Equivalence struct {
	SearchToken    string     `@Ident`
	SubSearchToken string     `[ "." @Ident ]`
	Operation      *Operation `@@`
}

type Operation struct {
	Numerical   *NumericalOperand `@@`
	Equivalence *VariableValue             `| "=" @@`
}

type NumericalOperand struct {
	Operator string  `@(">"|">="|"<"|"<=")`
	Val      float64 `@(Float|Int)`
}

type VariableValue struct {
	PrefixStar  bool   `[@"*"]`
	Value       *Value `@@`
	SuffixStart bool   `[@"*"]`
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
	err = parser.ParseString("id.ca = c123* AND blalb = *luia and label( status= malicious1 ) and score > 5 and label(c11c.v=44)", ggrammer)
	if err != nil {
		panic(err)
	}

	json.NewEncoder(os.Stdout).Encode(ggrammer)
}
