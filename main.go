package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	initFlags()
	mapResponses()
	log.Println("Starting http server")

	http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}
