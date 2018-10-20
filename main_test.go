package gogrammer

import (
	"testing"
)

func testFetchAllWorksUsing(t *testing.T, input string) {
	res, err := Parse(input)
	if err != nil {
		t.Error(err)
	}
	if !res.FetchAll {
		t.Error("Fetch all failed")
	}
	if res.Filters != nil {
		t.Error("Fetch all failed because it also parsed subsequent filters")
	}
}

func TestStar(t *testing.T) {
	testFetchAllWorksUsing(t, "*")
}

func TestLiteralAll(t *testing.T) {
	testFetchAllWorksUsing(t, "all")
	testFetchAllWorksUsing(t, "ALL")
	testFetchAllWorksUsing(t, "AlL")
	testFetchAllWorksUsing(t, "ALl")
}

func TestSimpleLabel(t *testing.T) {
	res, err := Parse("label(a=b")
	if err != nil {
		t.Error("Label parsing failed")
	}
	if *res.Filters[0].Labels[0].String != "a" && res.Filters[0].Labels[0].Number != nil {
		t.Error("Label left hand side operand parsing failed ")
	}
	if *res.Filters[0].Labels[1].String != "b" && res.Filters[0].Labels[0].Number != nil {
		t.Error("Label right hand side operand parsing failed ")
	}
}

func TestComplexFilter(t *testing.T) {
	input := "table.column = val123 aNd simple = complex AND song = 'alelu[a|b|c]{2,}ia.+' and label(malicious) and label  ( status = malicious1 ) and score > 5"
	res, err := ParseToJSON(input)
	if err != nil || res == "" {
		t.Error("Failed to parse complex query")
	}
}
