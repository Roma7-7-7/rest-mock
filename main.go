package main

import (
	"fmt"
	"log"
	"net/http"

	"app/shared"
)

func main() {
	mapResponses()

	log.Printf("Starting http server on port [%v]\n", shared.Port)
	http.ListenAndServe(fmt.Sprintf(":%v", shared.Port), nil)
}
