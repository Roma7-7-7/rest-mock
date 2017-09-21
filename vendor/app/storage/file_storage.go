package storage

import (
	"app/api"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

//FileStorage :
type FileStorage struct {
	endpoints []api.Endpoint
}

//GetAll :
func (s FileStorage) GetAll() ([]api.Endpoint, error) {
	return s.endpoints, nil
}

//BuildFileStorage :
func BuildFileStorage(responsesPath string) *FileStorage {

	var endpoints []api.Endpoint

	for _, f := range getFiles(responsesPath) {

		if !f.IsDir() {
			log.Printf("[%v] is not a directory\n", f.Name())
			continue
		}

		fullPath := filepath.Join(responsesPath, f.Name())
		if endpoint := buildEndpoint(fullPath); endpoint != nil {
			endpoint.Name = f.Name()
			endpoints = append(endpoints, *endpoint)
		}
	}

	return &FileStorage{
		endpoints: endpoints,
	}
}

func getFiles(responsesPath string) []os.FileInfo {
	if f, err := os.Stat(responsesPath); err != nil || !f.Mode().IsDir() {
		if err != nil {
			log.Print(err)
		}

		log.Fatalf("Failed to locate response files path [%v]\n", responsesPath)
	}

	files, err := ioutil.ReadDir(responsesPath)

	if err != nil {
		log.Println(err)
		log.Fatalf("Failed to locate responses directory [%v]", responsesPath)
	}

	return files
}

func buildEndpoint(path string) *api.Endpoint {
	apiFile := filepath.Join(path, "api.yml")
	data, err := ioutil.ReadFile(apiFile)
	result := api.Endpoint{}

	if err != nil {
		log.Printf("Failed to read file [%v]\n%v\n", apiFile, err)
		return nil
	}

	err = yaml.Unmarshal(data, &result)

	if err != nil {
		log.Printf("Failed to read file [%v]\n%v\n", apiFile, err)
		return nil
	}

	return &result
}
