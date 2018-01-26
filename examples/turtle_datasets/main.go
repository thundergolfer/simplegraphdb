package main

import (
	"fmt"
	"log"

	"github.com/thundergolfer/simplegraphdb"
)

const (
	exampleTurtleDb        string = "./example.turtle"
	countriesTurtleDb      string = "./countries.ttl"
	countriesBriefTurtleDb string = "./countries_brief.ttl"
)

func main() {
	store, err := simplegraphdb.InitHexastoreFromTurtle(countriesBriefTurtleDb)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println(store)
}
