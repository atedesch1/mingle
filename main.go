package main

import (
	"flag"
	"fmt"
)

func main() {
	addr := flag.String("listenaddr", ":3000", "the api address")
	flag.Parse()
    fmt.Printf("Listening on localhost%s\n", *addr)
}
