package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/alecthomas/repr"
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
		repr.Println(simplegraphdb.PresentableResults(results, store), repr.Indent("  "), repr.OmitEmpty(true))
	}
}
