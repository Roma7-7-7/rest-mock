package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	if AddMapping {
		AddRequestMapping()
		if !CheckTrue("Do you want to run server?") {
			return
		}
	}

	http.HandleFunc("/", MappingHandler)
	log.Printf("Starting http server on port [%v]\n", Port)
	http.ListenAndServe(fmt.Sprintf(":%v", Port), nil)
}
