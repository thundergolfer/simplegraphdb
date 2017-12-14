package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/thundergolfer/simplegraphdb"
)

const (
	defaultTwitterDb  string = "../../scripts/my_twitter_example_db.json"
	personalTwitterDb string = "../../scripts/your_twitter_example_db.json"
)

func main() {
	store, err := simplegraphdb.InitHexastoreFromJSONRows(personalTwitterDb)

	if err != nil {
		fmt.Println("INFO: Using default Twitter network DB")
		store, err = simplegraphdb.InitHexastoreFromJSONRows(defaultTwitterDb)

		if err != nil {
			fmt.Printf("Error: %v \n", err.Error())
			os.Exit(1)
		}
	}

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
