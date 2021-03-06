package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/thundergolfer/simplegraphdb"
)

func initTestHexastore() (*simplegraphdb.HexastoreDB, error) {
	store, err := simplegraphdb.InitHexastoreFromJSON("./db.json")
	return store, err
}

func main() {
	fmt.Println("vim-go")
	var result string
	store, err := initTestHexastore()
	if err != nil {
		fmt.Println("Error: failed to initialise example store")
		os.Exit(1)
	}
	exampleQuery := "SELECT ?x WHERE { ?x 'likes' 'money' }"

	fmt.Println("Example query: ")
	fmt.Println(exampleQuery)

	fmt.Println("Result of example query: ")
	result, err = simplegraphdb.RunQuery(exampleQuery, store)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)

	fmt.Println("END DEMO\n---------------------------------")

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter simpleSPARQL style query: ")
		query, _ := reader.ReadString('\n')
		results, err := simplegraphdb.RunQuery(query, store)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Result: ")
		fmt.Println("-------------------------------------------------------------------------")
		fmt.Println(results)
	}
}
