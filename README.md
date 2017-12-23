# Simple Graph Database - A Hexastore based Triple Store  [![Build Status](https://travis-ci.com/thundergolfer/simplegraphdb.svg?token=yHGWQ42iK2BPk1FjaUMc&branch=master)](https://travis-ci.com/thundergolfer/simplegraphdb)
------

Developed to help learn the Golang language, this project explores a
simple form of graph database, the *Triple Store*, by implementing the
*Hexastore* architecture from the paper [*"Hexastore: Sextuple Indexing
for Semantic Web Data
Management"*](http://www.vldb.org/pvldb/1/1453965.pdf). The project also
explores the [SPARQL](https://en.wikipedia.org/wiki/SPARQL) graph query language, by implementing a simplified
version of it which I dubbed `simplesparql`.

### Introduction - Usage

As a simple demo, we can scrape your Twitter network and query it with
this library. To do so:

1. `git clone` this library to your designated Go workspace
2. Run `./scripts/create_twitter_followers_graph_example_db.sh`. **Note:** If
   you don't have a Twitter account, skip this step and a default
Twitter network (mine) will be used.
3. Go into the `/examples/your_twitter_graph/` folder and run `go build`
4. A binary is created, which you can run with `./your_twitter_graph`
5. Query the network with `simplesparql` syntax (See the `README.md`
   section on it below).

##### Example starting query

Get all the people that you follow with: 

`SELECT ?x WHERE { '<YOUR SCREEN NAME' 'follows' ?x }` 

----------

### simplesparql

`simplesparql` is an implementation of a subset of the SPARQL query
language. It supports the asking of basic questions about a graph,
questions like "Who likes 'Manchester United'?", but not anything
involving aggregations, grouping, conditionals. With time, these things
may be implemented.

The basics of it are the you preface variables with `?`, with a variable
in the `WHERE` clause acting like a wildcard (`*`). For example:

`SELECT ?screen_name WHERE { 'jonobelotti_IO' 'follows' ?screen_name }`

will return all people that `jonobelotti_IO` follows (ie. all vertices
on the graph connected by a 'follows' edge from the vertice
`'jonobelotti_IO'`.

You can use up to a maximum of 3 variables in a `WHERE` clause. All
`WHERE` clauses must have exactly three components. For example: 

`SELECT ?screen_name, ?property WHERE { 'jonobelotti_IO' ?property
?screen_name }`

will essentially return all triples in the graph involving the
`'jonobelotti_IO'` vertice.
