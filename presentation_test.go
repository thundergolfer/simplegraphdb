package simplegraphdb

import "testing"

func TestPresentResultGrid(t *testing.T) {
	cases := []struct {
		resultsGrid [][]string
		expected    string
	}{
		{
			resultsGrid: [][]string{
				{"A", "B", "C"},
				{"D", "E", "F"},
				{"G", "H", "I"},
			},
			expected: ("A                              | B                              | C                              | \n" +
				"------------------------------------------------------------------------------------------------\n" +
				"D                              | E                              | F                              | \n" +
				"G                              | H                              | I                              | \n"),
		},
	}

	for _, c := range cases {
		actual := PresentResultGrid(&c.resultsGrid)
		if c.expected != actual {
			t.Errorf("Error: expected:\n%s\ngot:\n%s", c.expected, actual)
		}
	}
}
