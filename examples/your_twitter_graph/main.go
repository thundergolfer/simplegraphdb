package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/thundergolfer/simplegraphdb"
)

func main() {
	store := simplegraphdb.InitHexastoreFromJson("../../scripts/your_twitter_example_db.json")
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter simpleSPARQL style query: ")
		query, _ := reader.ReadString('\n')
		results := simplegraphdb.RunQuery(query, store)
		fmt.Println("Result: ")
		fmt.Println("-------------------------------------------------------------------------")
		fmt.Println(results)
	}
}
