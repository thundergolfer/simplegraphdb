package main

import (
	"fmt"

	"github.com/thundergolfer/simplegraphdb"
)

func main() {
	fmt.Println("vim-go")
	var result string
	query := "SELECT ?x WHERE { ?x 'likes' 'money' }"

	fmt.Println("Trying to do runQuery()")
	store := simplegraphdb.InitTestHexastore()
	result = simplegraphdb.RunQuery(query, store)
	fmt.Println(result)
}
