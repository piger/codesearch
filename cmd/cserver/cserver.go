package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/piger/codesearch/index"
	"github.com/piger/codesearch/server"
)

var (
	indexPath = flag.String("index", "$HOME/.csearchindex", "Path to the search index")
	address   = flag.String("addr", "127.0.0.1:40123", "Listen address")
)

func main() {
	flag.Parse()

	indexFilename := os.ExpandEnv(*indexPath)
	fmt.Printf("Using index %s\n", indexFilename)

	ix := index.Open(indexFilename)
	server.RunServer(*address, ix)
}
