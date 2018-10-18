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
	Equivalence *GenericValue     `| "=" @@`
}

type NumericalOperand struct {
	Operator string  `@(">"|">="|"<"|"<=")`
	Val      float64 `@(Float|Int)`
}

type GenericValue struct {
	String string  `@Ident`
	Number float64 `| @(Float|Int)`
	Star   bool    `[@"*"]`
}

func main() {
	parser, err := participle.Build(&GGrammer{})
	if err != nil {
		panic(err)
	}
	ggrammer := &GGrammer{}
	err = parser.ParseString("id.ca = 123* AND blalb = aleluia and label( status= malicious1 ) and score > 5", ggrammer)
	if err != nil {
		panic(err)
	}

	json.NewEncoder(os.Stdout).Encode(ggrammer)
}
