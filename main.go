package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rayou/go-theperfectgiftcard"
	"github.com/rayou/theperfectgiftcard-api/handler"
)

var port = fmt.Sprintf(":%s", os.Getenv("PORT"))

func main() {
	client, err := theperfectgiftcard.NewClient()
	if err != nil {
		panic(err)
	}

	log.Printf("Now listening on %s...\n", port)
	err = http.ListenAndServe(port, handler.NewHandler(client))
	if err != nil {
		panic(err)
	}
}
