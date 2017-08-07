package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type responseItem struct {
	Path string `yaml:"path"`
	File string `yaml:"file"`
}

func mapFile(path string, item responseItem) {
	responseFilePath := filepath.Join(path, item.File)

	if _, err := os.Stat(responseFilePath); os.IsNotExist(err) {
		log.Printf("Failed to find response file [%v]", responseFilePath)
		return
	}

	var h = func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile(responseFilePath)

		if err != nil {
			log.Printf("Failed to read file [%v]\n%v\n", responseFilePath, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(data)
	}

	http.HandleFunc(item.Path, h)
}

func mapResponse(path string) {
	apiFile := filepath.Join(path, "api.yml")
	data, err := ioutil.ReadFile(apiFile)
	item := responseItem{}

	if err != nil {
		log.Printf("Failed to read file [%v]\n%v\n", apiFile, err)
		return
	}

	err = yaml.Unmarshal(data, &item)

	if err != nil {
		log.Printf("Failed to read file [%v]\n%v\n", apiFile, err)
		return
	}

	if len(item.File) > 0 {
		mapFile(path, item)
	}
}

func mapResponses() {
	log.Println("Mapping request to responses")

	files, err := ioutil.ReadDir(responseFolder)

	if err != nil {
		log.Println(err)
		log.Fatalf("Failed to locate responses directory [%v]", responseFolder)
	}

	for _, f := range files {
		fullPath := filepath.Join(responseFolder, f.Name())

		if !f.IsDir() {
			log.Printf("[%v] is not a directory\n", f.Name())
			continue
		}

		mapResponse(fullPath)
	}
}
