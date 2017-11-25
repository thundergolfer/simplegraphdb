# Simple Graph Database - A Hexastore based Triple Store  [![Build Status](https://travis-ci.com/thundergolfer/simplegraphdb.svg?token=yHGWQ42iK2BPk1FjaUMc&branch=master)](https://travis-ci.com/thundergolfer/simplegraphdb)
------

Developed to help learn the Golang language, this project explores a
simple form of graph database, the *Triple Store*, by implementing the
*Hexastore* architecture from the paper [*"Hexastore: Sextuple Indexing
for Semantic Web Data
Management"*](http://www.vldb.org/pvldb/1/1453965.pdf). The project also
explores the [SPARQL](https://en.wikipedia.org/wiki/SPARQL) graph query language, by implementing a simplified
version of it which I dubbed `simplesparql`.

### Usage

As a simple demo, we can scrape your Twitter network and query it with
this library. To do so:

1. `git clone` this library to your designated Go workspace
2. Run `./create_twitter_followers_graph_example_db.sh`. **Note:** If
   you don't have a Twitter account, skip this step and a default
Twitter network (mine) will be used.
3. Go into the `/examples/your_twitter_graph/` folder and run `go build`
4. A binary is created, which you can run with `./your_twitter_graph`
5. Query the network with `simplesparql` syntax (See the `README.md`
   section on it below).
