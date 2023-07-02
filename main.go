package main

import (
	"flag"
	"fmt"

	"github.com/atedesch1/mingle/cmd"
)

func main() {
	addr := flag.String("listenaddr", ":3000", "the api address")
	flag.Parse()
    fmt.Printf("Listening on localhost%s\n", *addr)
    cmd.StartServer()
}
