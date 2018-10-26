package gogrammer

import (
	"testing"
)

func testFetchAllWorksUsing(t *testing.T, input string) {
	res, err := Parse(input)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if !res.FetchAll {
		t.Error("Fetch all failed")
		panic("Boolean false")
	}
	if len(res.Or) != 0 {
		t.Error("Fetch all failed because it also parsed subsequent filters")
		panic("noo")
	}
}

func TestStar(t *testing.T) {
	testFetchAllWorksUsing(t, "*")
}

func TestLiteralAll(t *testing.T) {
	testFetchAllWorksUsing(t, "all")
	testFetchAllWorksUsing(t, "AlL")
}

func TestSimpleLabel(t *testing.T) {
	res, err := Parse("label(a=b)")
	if err != nil {
		t.Error("Label parsing failed", err)
		return
	}
	if *res.Or[0].Filters[0].Labeled.Operand.Value.String != "a" && res.Or[0].Filters[0].Labeled.Operand.Value.Number != nil {
		t.Error("Label left hand side operand parsing failed ")
		return
	}
	if *res.Or[0].Filters[0].Labeled.Operand.Value.String != "b" && res.Or[0].Filters[0].Labeled.Operand.Value.Number != nil {
		t.Error("Label right hand side operand parsing failed ")
		return
	}
	if res.Or[0].Filters[0].Labeled.IsNegation {
		t.Error("Label is negated when it should not")
		return
	}
}

func TestNumericalExpression(t *testing.T) {
	res, err := Parse("jjjj >= 55")
	if err != nil {
		t.Error("Numeric expression parsing failed", err)
		return
	}
	if res.Or[0].Filters[0].Expression.Token != "jjjj" {
		t.Error("Token parse failed")
		return
	}
	if res.Or[0].Filters[0].Expression.Operand.Numeric.Operator != ">=" {
		t.Error("Numerical operand failed", res.Or[0].Filters[0].Expression.Operand.Numeric.Operator)
		return
	}
	if res.Or[0].Filters[0].Expression.Operand.Numeric.Val != 55 {
		t.Error("Numerical value failed", res.Or[0].Filters[0].Expression.Operand.Numeric.Val)
		return
	}
}

func TestNegatedLabel(t *testing.T) {
	res, err := Parse("label(caa!='yyy1')")
	if err != nil {
		t.Error("Negated label parsing failed", err)
		return
	}
	if !res.Or[0].Filters[0].Labeled.IsNegation {
		t.Error("Label is not negated")
		return
	}
}

func TestComplexFilter(t *testing.T) {
	res, err := parseComplexFilter()
	if err != nil || res == "" {
		t.Error("Failed to parse complex query", err, res)
	}
}

func parseComplexFilter() (string, error) {
	input := "table.column = val123 aNd simple = complex AND song != 'alelu[a|b|c]{2,}ia.+' and label(malicious) or label  ( status = malicious1 ) and score > 5"
	return ParseToJSON(input)
}

func BenchmarkComplexParsingParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			parseComplexFilter()
		}
	})
}

func TestOrExpressions(t *testing.T) {
	res, err := Parse("ij != 44 and label(malicious) or ij = 44 and label(status!=cleared)")
	if err != nil {
		t.Error("Parsed failed", err)
		return
	}
	if len(res.Or) != 2 && res.Or[1].Filters[1].Labeled.IsNegation {
		t.Error("Negation failed")
		return
	}
}
