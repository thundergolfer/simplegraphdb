package simplegraphdb

import "strings"

func PresentResultGrid(resultsGrid *[][]string) (presentable string) {
	max_cell_size := 30
	presentable = ""
	for i, row := range *resultsGrid {
		for _, col := range row {
			curr_len := len(col)
			if curr_len >= max_cell_size {
				presentable += col[:max_cell_size-2] + ".." + " | "
			} else {
				presentable += rightPad2Len(col, " ", max_cell_size) + " | "
			}
		}
		presentable += "\n"
		if i == 0 {
			presentable += strings.Repeat("-", (max_cell_size+2)*len(row)) + "\n"
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
