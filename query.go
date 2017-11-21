package main

import (
	"fmt"

	"github.com/thundergolfer/simple-graph-database/simplesparql"
)

func main() {
	fmt.Println("Parsing query: ")
	simplesparql.Parse("SELECT ?x WHERE { ?x 'likes' 'eminem' }")
}
