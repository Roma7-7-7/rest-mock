package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", MappingHandler)
	log.Printf("Starting http server on port [%v]\n", Port)
	http.ListenAndServe(fmt.Sprintf(":%v", Port), nil)
}
