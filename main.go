package gogrammer

// GGrammer is the root of the grammer
type GGrammer struct {
	FetchAll bool   `@("*"|"all")`
	Or       []*And `| ( @@ { "or" @@ } )`
}

type And struct {
	Filters []*Filter `@@ { "and" @@ }`
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
	Simple *SimpleValue `@@`
	Regex  string       `| @(String|RawString)`
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
