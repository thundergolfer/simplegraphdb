package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	// "github.com/alecthomas/repr"
	"github.com/thundergolfer/simple-graph-database/simplesparql"
)

func main() {
	fmt.Println("Parsing query: ")
	queryModel := simplesparql.Parse("SELECT ?x WHERE { ?x 'likes' 'eminem' }")
	returnVars := extractReturnVariables(queryModel)
	first, second, third := extractTripleExpressionElements(queryModel)
	fmt.Println(first)
	fmt.Println(second)
	fmt.Println(third)

	_ = returnVars // TODO
}

func runQuery(query string, hexastore *Hexastore) *[]Triple {
	queryModel := simplesparql.Parse(query)

	_, err := validateQuery(queryModel, hexastore)
	if err != nil {
		log.Fatal(err)
	}

	extractReturnVariables(queryModel)
	first, second, third := extractTripleExpressionElements(queryModel)

	if isSparqlVariable(first) { // X??
		if isSparqlVariable(second) { // XX?
			if isSparqlVariable(third) { // XXX
				return hexastore.QueryXXX()
			}
			objId, _ := hexastore.entities.GetKey(third)
			return hexastore.QueryXXO(objId) // XXO
		} else if isSparqlVariable(third) { // XPX
			propId, _ := hexastore.props.GetKey(second)
			return hexastore.QueryXPX(propId)
		} // XPO

		propId, _ := hexastore.props.GetKey(second)
		objId, _ := hexastore.entities.GetKey(third)
		return hexastore.QueryXPO(propId, objId)
	} else if isSparqlVariable(second) { // SX?
		subjId, _ := hexastore.entities.GetKey(first)
		if isSparqlVariable(third) { // SXX
			return hexastore.QuerySXX(subjId)
		} // SXO
		objId, _ := hexastore.entities.GetKey(third)
		return hexastore.QuerySXO(subjId, objId)
	} else if isSparqlVariable(third) { // SPX
		subjId, _ := hexastore.entities.GetKey(first)
		propId, _ := hexastore.props.GetKey(second)
		return hexastore.QuerySPX(subjId, propId)
	} // SPO

	subjId, _ := hexastore.entities.GetKey(first)
	propId, _ := hexastore.props.GetKey(second)
	objId, _ := hexastore.entities.GetKey(third)
	return hexastore.QuerySPO(subjId, propId, objId)
}

func validateQuery(queryModel *(simplesparql.Select), hexastore *Hexastore) (bool, error) {
	var ok bool
	_ = hexastore // TODO do validations against hexastore

	returnVars := extractReturnVariables(queryModel)
	first, second, third := extractTripleExpressionElements(queryModel)
	tripleExprElems := []string{first, second, third}

	ok = validateNoDuplicateVariables(returnVars)
	if !ok {
		return false, errors.New("Duplicate variable name in SELECT variables")
	}

	ok = validateNoDuplicateVariables(tripleExprElems)
	if !ok {
		return false, errors.New("Duplicate variable name in WHERE expression variables")
	}

	whereVars := getVariablesFromStrings(tripleExprElems...)

	ok = validateVariablesBalance(returnVars, whereVars)
	if !ok {
		return false, errors.New("Cant fulfil SELECT expression with variables from WHERE expression")
	}

	return true, nil
}

func extractReturnVariables(queryModel *(simplesparql.Select)) (returnVars []string) {
	selectExpressions := queryModel.Expression.Expressions

	for _, expr := range selectExpressions {
		variable := getVariableFromExpression(expr)
		returnVars = append(returnVars, variable)
	}

	return
}

func getVariableFromExpression(expr *(simplesparql.Expression)) (variable string) {
	variable = expr.And.Or[0].Operand.Operand.Summand[0].LHS.LHS.Var
	return
}

func isSparqlVariable(val string) bool {
	return strings.HasPrefix(val, "?")
}

func getVariablesFromStrings(strings ...string) (variables []string) {
	variables = []string{}

	for _, val := range strings {
		if isSparqlVariable(val) {
			variables = append(variables, val)
		}
	}
	return
}

func validateNoDuplicateVariables(vars []string) bool {
	checker := map[string]bool{}

	for _, elem := range vars {
		if _, ok := checker[elem]; ok {
			return false
		}

		checker[elem] = true
	}

	return true
}

func validateVariablesBalance(selectExprVars, whereExprVars []string) bool {
	checker := map[string]bool{}

	for _, elem := range whereExprVars {
		checker[elem] = true
	}

	for _, elem := range selectExprVars {
		if _, ok := checker[elem]; ok {
			return false
		}
	}

	return true
}

func extractTripleExpressionElements(queryModel *(simplesparql.Select)) (first, second, third string) {
	where := queryModel.Where

	if where.Expression.First.Value != nil {
		first = *where.Expression.First.Value.String
	} else {
		first = where.Expression.First.Var
	}
	//, where.Expression.Second, where.Expression.Third.Value
	if where.Expression.Second.Value != nil {
		second = *where.Expression.Second.Value.String
	} else {
		second = where.Expression.Second.Var
	}

	if where.Expression.Third.Value != nil {
		third = *where.Expression.Third.Value.String
	} else {
		third = where.Expression.Third.Var
	}

	return
}
