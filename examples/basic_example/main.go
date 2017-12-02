package main

import (
	"fmt"

	"github.com/alecthomas/repr"
	"github.com/thundergolfer/simplegraphdb"
)

func main() {
	fmt.Println("vim-go")
	var results *[][]string
	query := "SELECT ?x WHERE { ?x 'likes' 'money' }"

	fmt.Println("Trying to do runQuery()")
	store := simplegraphdb.InitTestHexastore()
	results = simplegraphdb.RunQuery(query, store)
	repr.Println(results, repr.Indent("  "), repr.OmitEmpty(true))
	// repr.Println(simplegraphdb.PresentableResults(results, store), repr.Indent("  "), repr.OmitEmpty(true))
}
