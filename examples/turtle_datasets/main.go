package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/thundergolfer/simplegraphdb"
)

const (
	exampleTurtleDb        string = "./example.turtle"
	countriesTurtleDb      string = "./countries.ttl" // broken: currently the Turtle parsing library used can't handle the 'a' predicate
	countriesBriefTurtleDb string = "./countries_brief.ttl"
)

func main() {
	store, err := simplegraphdb.InitHexastoreFromTurtle(countriesBriefTurtleDb)
	if err != nil {
		log.Fatalf(err.Error())
	}

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
