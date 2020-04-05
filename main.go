package main

import (
	"log"
	"net/http"

	"github.com/nocturnalBadger/memecan/connectors"
	"github.com/nocturnalBadger/memecan/app"
)

func main() {
	log.Println("Starting memecan server")

	router := app.Routes()

	connectors.InitMinio()

	connectors.InitES()

	http.ListenAndServe(":3000", router)
}
