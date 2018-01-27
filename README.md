# Simple Graph Database - A Hexastore based Triple Store  [![Build Status](https://travis-ci.com/thundergolfer/simplegraphdb.svg?token=yHGWQ42iK2BPk1FjaUMc&branch=master)](https://travis-ci.com/thundergolfer/simplegraphdb)
------

Developed to help learn the Golang language, this project explores a
simple form of graph database, the *Triple Store*, by implementing the
*Hexastore* architecture from the paper [*"Hexastore: Sextuple Indexing
for Semantic Web Data
Management"*](http://www.vldb.org/pvldb/1/1453965.pdf). The project also
explores the [SPARQL](https://en.wikipedia.org/wiki/SPARQL) graph query language, by implementing a simplified
version of it which I dubbed `simplesparql`.

## Introduction - Usage

As a simple demo, we can scrape your Twitter network and query it with
this library. To do so:

1. `git clone` this library to your designated Go workspace
2. Go into the `scripts/` directory
3. Install script dependencies with `pip install -r requirements.txt` (consider using `virtualenv`).
4. Run `./create_twitter_followers_graph_example_db.sh`. **Note:** If you don't have a Twitter account, skip this step and a default Twitter network (mine) will be used.
5. Go into the `/examples/your_twitter_graph/` folder and run `go build`
6. A binary is created, which you can run with `./your_twitter_graph`
7. Query the n etwork with `simplesparql` syntax (See the `README.md`
   section on it below).

##### Example starting query

Get all the people that you follow with:

`SELECT ?x WHERE { '<YOUR SCREEN NAME>' 'follows' ?x }`

## Usage - Golang Package

1 .Add to project:

`go get github.com/thundergolfer/simplegraphdb` or `glide get github.com/thundergolfer/simplegraphdb`

2. Add to imports:

```golang
import (
  "github.com/thundergolfer/simplegraphdb"
)
```

#### Package Interface

##### `InitHexastoreFromJSONRows(filename string) (*HexastoreDB, error)`

You can setup a Hexastore by passing a filepath to a `.json` file with the following format:

```
{"subject": <STRING>, "prop": <STRING>, "object": <STRING>}
{"subject": <STRING>, "prop": <STRING>, "object": <STRING>}
{"subject": <STRING>, "prop": <STRING>, "object": <STRING>}
...
...
```

##### `InitHexastoreFromTurtle(dbFilePath string) (Hexastore, error)`

> BETA: Turtle files are currently parsed with [d4l3k/turtle](https://github.com/d4l3k/turtle), which is quite limited.

[*Turtle* (Terse RDF Triple Language)](https://www.w3.org/TeamSubmission/turtle/) is a syntax for describing RDF semantic web graphs. You can load an RDF graph specified in turtle syntax with this function.

##### `RunQuery(query string, store Hexastore) (string, error)`

Run a well-formed `simplesparql` query (see more below) against a Hexastore instance. Just returns a printable table of results like:

```
?target                | ?prop                  |
------------------------------------------------
Cow                    | follows                |
Apple                  | follows                |
```

##### `SaveToJSONRows(filename string, store Hexastore) error`

> Not yet implemented, but coming soon


----------

## simplesparql

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
