package simplegraphdb

import (
	"fmt"
	"strings"

	"github.com/thundergolfer/simplegraphdb/simplesparql"
)

// RunQuery takes `simplesparql` valid string query and a Hexastore instance
// and returns a formatted table of query results
func RunQuery(query string, hexastore Hexastore) (string, error) {
	resultsGrid, err := runQuery(query, hexastore)
	if err != nil {
		return "", err
	}
	return PresentResultGrid(resultsGrid), nil
}

func runQuery(query string, hexastore Hexastore) ([][]string, error) {
	queryModel, err := simplesparql.Parse(query)
	if err != nil {
		return [][]string{}, err
	}

	err = validateQuery(queryModel)
	if err != nil {
		return [][]string{}, err
	}

	returnVars := extractReturnVariables(queryModel)
	rawResults := retreiveQueryResults(queryModel, hexastore)
	mappedResults := mapTriplePartsToVars(hexastore, queryModel, rawResults)
	resultsGrid := buildResultsGrid(returnVars, mappedResults, len(*rawResults))

	return resultsGrid, nil
}

func buildResultsGrid(returnVars []string, mappedResults map[string][]string, numResults int) [][]string {
	stringResults := make([][]string, numResults+1)
	stringResults[0] = returnVars // add header

	for i := 0; i < numResults; i++ {
		stringResults[i+1] = make([]string, len(returnVars))
		for j, rVar := range returnVars {
			stringResults[i+1][j] = mappedResults[rVar][i]
		}
	}

	return stringResults
}

func retreiveQueryResults(queryModel *simplesparql.Select, hexastore Hexastore) *[]Triple {
	first, second, third := extractTripleExpressionElements(queryModel)

	if isSparqlVariable(first) { // X??
		if isSparqlVariable(second) { // XX?
			if isSparqlVariable(third) { // XXX
				return hexastore.QueryXXX()
			}
			objID, _ := hexastore.GetEntityKey(third)
			return hexastore.QueryXXO(objID) // XXO
		} else if isSparqlVariable(third) { // XPX
			propID, _ := hexastore.GetPropKey(second)
			return hexastore.QueryXPX(propID)
		} // XPO

		propID, _ := hexastore.GetPropKey(second)
		objID, _ := hexastore.GetEntityKey(third)
		return hexastore.QueryXPO(propID, objID)
	} else if isSparqlVariable(second) { // SX?
		subjID, _ := hexastore.GetEntityKey(first)
		if isSparqlVariable(third) { // SXX
			return hexastore.QuerySXX(subjID)
		} // SXO
		objID, _ := hexastore.GetEntityKey(third)
		return hexastore.QuerySXO(subjID, objID)
	} else if isSparqlVariable(third) { // SPX
		subjID, _ := hexastore.GetEntityKey(first)
		propID, _ := hexastore.GetPropKey(second)
		return hexastore.QuerySPX(subjID, propID)
	} // SPO

	subjID, _ := hexastore.GetEntityKey(first)
	propID, _ := hexastore.GetPropKey(second)
	objID, _ := hexastore.GetEntityKey(third)
	return hexastore.QuerySPO(subjID, propID, objID)
}

func mapTriplePartsToVars(store Hexastore, queryModel *simplesparql.Select, results *[]Triple) (mappedParts map[string][]string) {
	var firstVarName, scndVarName, thirdVarName string
	numResults := len(*results)
	mappedParts = make(map[string][]string)
	first, second, third := extractTripleExpressionElements(queryModel)

	if isSparqlVariable(first) {
		firstVarName = first
		mappedParts[firstVarName] = make([]string, numResults)

	}
	if isSparqlVariable(second) {
		scndVarName = second
		mappedParts[scndVarName] = make([]string, numResults)
	}
	if isSparqlVariable(third) {
		thirdVarName = third
		mappedParts[thirdVarName] = make([]string, numResults)
	}

	for i, triple := range *results {
		if firstVarName != "" {
			mappedParts[firstVarName][i] = store.ResolveEntity(triple.Subject)
		}
		if scndVarName != "" {
			mappedParts[scndVarName][i] = store.ResolveProp(triple.Prop)
		}
		if thirdVarName != "" {
			mappedParts[thirdVarName][i] = store.ResolveEntity(triple.Object)
		}
	}

	return
}

func validateQuery(queryModel *(simplesparql.Select)) error {
	var ok bool

	returnVars := extractReturnVariables(queryModel)
	first, second, third := extractTripleExpressionElements(queryModel)
	tripleExprElems := []string{first, second, third}

	ok = validateNoDuplicateVariables(returnVars)
	if !ok {
		return fmt.Errorf("Duplicate variable name in SELECT variables")
	}

	ok = validateNoDuplicateVariables(tripleExprElems)
	if !ok {
		return fmt.Errorf("Duplicate variable name in WHERE expression variables")
	}

	whereVars := getVariablesFromStrings(tripleExprElems...)

	ok = validateVariablesBalance(returnVars, whereVars)
	if !ok {
		return fmt.Errorf("Cant fulfil SELECT expression with variables from WHERE expression")
	}

	return nil
}

func extractReturnVariables(queryModel *(simplesparql.Select)) (returnVars []string) {
	selectExpressions := queryModel.Expression.Expressions

	for _, expr := range selectExpressions {
		variable := getVariableFromExpression(expr)
		returnVars = append(returnVars, variable)
	}
	return
}

func getVariableFromExpression(expr *(simplesparql.Expression)) string {
	return expr.And.Or[0].Operand.Operand.Summand[0].LHS.LHS.Var
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
		if _, ok := checker[elem]; !ok {
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
