package simplegraphdb

import "strings"

// PresentResultGrid converts the 2D grid of string elements into a
// single string representation of a table for presentation via stdout
func PresentResultGrid(resultsGrid [][]string) (presentable string) {
	maxCellSize := 30
	presentable = ""
	for i, row := range resultsGrid {
		for _, col := range row {
			currLen := len(col)
			if currLen >= maxCellSize {
				presentable += col[:maxCellSize-2] + ".." + " | "
			} else {
				presentable += rightPad2Len(col, " ", maxCellSize) + " | "
			}
		}
		presentable += "\n"
		if i == 0 {
			presentable += strings.Repeat("-", (maxCellSize+2)*len(row)) + "\n"
		}
	}
	return
}

func rightPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = s + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}
