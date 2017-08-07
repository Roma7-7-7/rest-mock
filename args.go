package main

import (
	"flag"
	"log"
)

var port = 8888
var responseFolder = "responses"

func initFlags() {
	log.Println("Initialising application flags")

	flag.IntVar(&port, "p", 8888, "Port that will be used to run application")
	flag.StringVar(&responseFolder, "rf", "responses", "Path to response folder")

	flag.Parse()
}
