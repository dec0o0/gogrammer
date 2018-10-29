package gogrammer

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
)

// ValueMapper is cool
type ValueMapper interface {
	Stringify(tokenName string, operator string, value string) string
}

type colMetadata struct {
	tableName *string
	mapper    *ValueMapper
}

// ParseToSQL transforms a input string to SQL
// returns a list of where clauses
func ParseToSQL(input string, tablesColumnsMapper map[string]map[string]*ValueMapper) (string, []error) {
	var result string
	var errs []error
	parsed, err := Parse(input)
	if err != nil {
		errs = []error{err}
	} else if parsed == nil {
		errs = []error{errors.New("Parsed input is actually null")}
	} else if parsed.FetchAll {

	} else {
		traverse(parsed, adapt(tablesColumnsMapper))
	}
	return result, errs
}

func (e *Expression) getFullToken() (result string) {
	if e.SubToken == "" {
		result = e.Token
	} else {
		result = e.Token + "." + e.SubToken
	}
	return
}

func (e *Expression) getValue() (result string) {
	if op := e.Operand; op != nil {
		if num := op.Numeric; num != nil {
			result = strconv.FormatFloat(num.Val, 'f', 2, 64)
		} else if lit := op.Literal; lit != nil {
			if regex := lit.Value.Regex; regex != "" { // TODO: check out what do with regex
				result = regex
			} else if simple := lit.Value.Simple; simple != nil {
				if nr := simple.Number; nr != nil {
					result = strconv.FormatFloat(num.Val, 'f', 2, 64)
				} else {
					result = *simple.String
				}
			}
		}
	}
	return
}

func (e *Expression) getOperand() (result string) {
	if op := e.Operand; op != nil {
		if num := op.Numeric; num != nil {
			result = num.Operator
		} else if lit := op.Literal; lit != nil {
			if lit.IsNegation {
				result = "!="
			} else {
				result = "="
			}
		}
	}
	return
}

func traverse(asTree *GGrammer, columnInfoMap map[string]colMetadata) ([]string, []string, []error) {
	var from = make([]string, 0)
	var where = make([]string, 0)
	var errs = make([]error, 0)
	for _, and := range asTree.Or {
		for _, filter := range and.Filters {
			if expr := filter.Expression; expr != nil {
				fullToken := expr.getFullToken()
				if metadata, found := columnInfoMap[fullToken]; found {
					from = append(from, *metadata.tableName)
					where = append(where, (*(metadata.mapper)).Stringify(fullToken, expr.getOperand(), expr.getValue()))
				} else {
					errs = append(errs, errors.New("Column "+fullToken+" is invalid"))
				}
			} else if label := filter.Labeled; label != nil {

			}
		}
	}
	return from, where, errs
}

func adapt(tablesColumnsMapper map[string]map[string]*ValueMapper) (columnInfoMap map[string]colMetadata) {
	var totalSize = 0
	for _, cols := range tablesColumnsMapper {
		totalSize += len(cols)
	}
	columnInfoMap = make(map[string]colMetadata, totalSize)
	for table, cols := range tablesColumnsMapper {
		for colName, mapper := range cols {
			columnInfoMap[colName] = colMetadata{&table, mapper}
		}
	}
	return
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
