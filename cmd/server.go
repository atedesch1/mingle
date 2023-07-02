package cmd

import (
	"log"

	"github.com/atedesch1/mingle/db"
)

func StartServer() {
    _, err := db.NewStorage()
    if err != nil {
        log.Fatal(err)
    }
}
