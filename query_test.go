package simplegraphdb

import (
	"fmt"
	"testing"

	"github.com/go-test/deep"
)

func createTestHexastore() *HexastoreDB {
	h := newHexastore()

	h.Add("Apple", "Likes", "Cow", "jonob")
	h.Add("Apple", "Likes", "Banana", "jonob")
	h.Add("Apple", "Likes", "Apple", "jonob")
	h.Add("Cow", "Dislikes", "Banana", "jonob")
	h.Add("Banana", "Dislikes", "Cow", "jonob")
	h.Add("Cow", "Likes", "Apple", "jonob")

	return h
}

// returned results grid are non-deterministic, so until
// deterministic output can be introduced to hexastore implementation
// this function is used to check equality
func checkResultsEquality(r1, r2 [][]string) bool {
	// check header
	if diff := deep.Equal(r1[0], r2[0]); diff != nil {
		return false
	}

	for _, row1 := range r1 {
		foundMatch := false
		// don't check if row has already been matched
		// there's no duplicates in a results grid
		for _, row2 := range r2 {
			if diff := deep.Equal(row1, row2); diff == nil {
				foundMatch = true
				break
			}
		}

		if !foundMatch {
			return false
		}
	}

	return true
}

// runQuery is tested instead of RunQuery because testing long table-form string equality is awkward
func Test_runQueryOnValidQueries(t *testing.T) {
	hexastore := createTestHexastore()
	cases := []struct {
		comment  string
		query    string
		expected [][]string
	}{
		{
			comment: "single variable return",
			query:   "SELECT ?x WHERE { ?x 'Likes' 'Banana' }",
			expected: [][]string{
				[]string{"?x"},
				[]string{"Apple"},
			},
		},
		{
			comment: "multiple variable return",
			query:   "SELECT ?x, ?y WHERE { ?x 'Dislikes' ?y }",
			expected: [][]string{
				[]string{"?x", "?y"},
				[]string{"Cow", "Banana"},
				[]string{"Banana", "Cow"},
			},
		},
		{
			comment: "out of order multiple variable return",
			query:   "SELECT ?y, ?x WHERE { ?x 'Likes' ?y }",
			expected: [][]string{
				[]string{"?y", "?x"},
				[]string{"Cow", "Apple"},
				[]string{"Banana", "Apple"},
				[]string{"Apple", "Apple"},
				[]string{"Apple", "Cow"},
			},
		},
	}

	for _, c := range cases {
		actual, err := runQuery(c.query, hexastore)

		if err != nil {
			t.Errorf("Error in test '%s': expected no error but got %s", c.comment, err.Error())
		}
		if len(c.expected) != len(actual) {
			t.Errorf("Error in test '%s': expected %d rows of data, got %d", c.comment, len(c.expected), len(actual))
		}
		if !checkResultsEquality(c.expected, actual) {
			diff := deep.Equal(c.expected, actual)
			t.Error(diff)
		}
	}
}

func Test_runQueryWithWellFormedButInvalidQueries(t *testing.T) {
	hexastore := createTestHexastore()
	cases := []struct {
		query    string
		expected error
	}{
		{
			query:    "SELECT ?x, ?x WHERE { ?x 'Likes' 'Banana' }",
			expected: fmt.Errorf("Duplicate variable name in SELECT variables"),
		},
		{
			query:    "SELECT ?x, ?y WHERE { ?x ?x ?y }",
			expected: fmt.Errorf("Duplicate variable name in WHERE expression variables"),
		},
		{
			query:    "SELECT ?y, ?x WHERE { ?x 'Likes' 'Stuff' }",
			expected: fmt.Errorf("Cant fulfil SELECT expression with variables from WHERE expression"),
		},
	}

	for _, c := range cases {
		_, err := runQuery(c.query, hexastore)
		if err == nil {
			t.Errorf("FAIL: expected error '%s' but got no error", c.expected)
		} else if err.Error() != c.expected.Error() {
			t.Errorf("FAIL: expected error '%s', got '%s'", c.expected, err)
		}
	}
}

func Test_runQueryWithMalformedQueries(t *testing.T) {
	hexastore := createTestHexastore()
	malformedQuery := "SELECT WHERE { ?x 'Likes' 'Banana' }"
	_, err := runQuery(malformedQuery, hexastore)
	if err == nil {
		t.Errorf("FAIL: expected error from malformed simplesparql query but got no error")
	}
}
