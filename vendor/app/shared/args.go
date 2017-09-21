package shared

import (
	"flag"
	"log"
)

//Port :
var Port = 8888

//ResponseFolder :
var ResponseFolder = "responses"

func init() {
	log.Println("Initialising application flags")

	flag.IntVar(&Port, "p", 8888, "Port that will be used to run application")
	flag.StringVar(&ResponseFolder, "rf", "responses", "Path to response folder")

	flag.Parse()
}
