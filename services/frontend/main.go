package main

import (
	"log"

	"github.com/joshgav/vanpool-manager/services/frontend/actions"
)

func main() {
	app := actions.App()
	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}
