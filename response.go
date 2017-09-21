package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"app/api"
	"app/shared"
	"app/storage"
)

func toLowercase(keys map[string][]string) map[string][]string {
	result := make(map[string][]string)

	for k, v := range keys {
		result[strings.ToLower(k)] = v
	}

	return result
}

func checkHeaders(r *http.Request, e api.Endpoint) bool {
	headers := toLowercase(r.Header)

	for _, h := range e.Expects.Headers {
		keys, ok := headers[strings.ToLower(h.Key)]

		if !ok || len(keys) != 1 {
			log.Printf("Request do not contain expected header [%v] or more than one value for that key", h.Key)
			return false
		}
		//TODO: add support values
		if keys[0] != h.Value {
			log.Printf("Request do not contain expected header [%v]:[%v]", h.Key, h.Value)
			return false
		}
	}

	return true
}

func mapStatus(w http.ResponseWriter, e api.Endpoint) {
	if e.Response.Status == 0 {
		return
	}

	w.WriteHeader(e.Response.Status)
}

func mapHeaders(w http.ResponseWriter, endpoint api.Endpoint) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	for _, header := range endpoint.Response.Headers {
		w.Header().Set(header.Key, header.Value)
	}
}

func mapFile(w http.ResponseWriter, e api.Endpoint) {
	if len(e.Response.ResponseFilePath) == 0 {
		w.Write(e.Response.Data)
		return
	}

	responseFilePath := filepath.Join(shared.ResponseFolder, e.Name, e.Response.ResponseFilePath)

	if _, err := os.Stat(responseFilePath); os.IsNotExist(err) {
		log.Printf("Failed to find response file [%v]", responseFilePath)
		return
	}

	data, err := ioutil.ReadFile(responseFilePath)

	if err != nil {
		log.Printf("Failed to read file [%v]\n%v\n", responseFilePath, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func mapResponse(e api.Endpoint) {
	var h = func(w http.ResponseWriter, r *http.Request) {
		if !checkHeaders(r, e) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		mapHeaders(w, e)
		mapStatus(w, e)
		mapFile(w, e)
	}

	http.HandleFunc(e.Path, h)
}

func mapResponses() {
	log.Println("Mapping request to responses")

	_, err := ioutil.ReadDir(shared.ResponseFolder)

	if err != nil {
		log.Println(err)
		log.Fatalf("Failed to locate responses directory [%v]", shared.ResponseFolder)
	}

	endpoints, err := storage.Storage.GetAll()

	if err != nil {
		log.Println(err)
		log.Fatal("Failed to load endpoints")
	}

	for _, e := range endpoints {
		mapResponse(e)
	}
}
