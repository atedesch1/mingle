package cmd

import (
	"log"

	"github.com/atedesch1/mingle/api"
	"github.com/atedesch1/mingle/db"
)

func StartServer(addr string) {
    storage, err := db.NewStorage()
    if err != nil {
        log.Fatal(err)
    }

	if err := api.NewHandler(storage, addr).Serve(); err != nil {
		log.Fatal(err)
	}
}
