package main

import (
	"errors"
	"log"
	"net/http"
	"reflect"
	"strings"
)

//Expects :
type Expects struct {
	Params  []Param  `yaml:"params"`
	Headers []Header `yaml:"headers"`
}

//ResponseData :
type ResponseData struct {
	Status  int                 `yaml:"status"`
	Headers map[string][]string `yaml:"headers"`
	Data    []byte              `yaml:"data"`
}

//RequestMapping :
type RequestMapping struct {
	Method  string
	Path    string
	Headers map[string][]string
	Params  map[string][]string
}

func toLowerHeaders(r *http.Request) map[string][]string {
	result := make(map[string][]string)

	for key, value := range r.Header {
		result[strings.ToLower(key)] = value
	}

	return result
}

func toLowerList(l []string) []string {
	var result []string

	for _, val := range l {
		result = append(result, strings.ToLower(val))
	}

	return result
}

func toLowerMap(m map[string][]string) map[string][]string {
	result := make(map[string][]string)

	for key, values := range m {
		result[strings.ToLower(key)] = toLowerList(values)
	}

	return result
}

func inSlice(slice []string, value string) bool {
	for _, val := range slice {
		if val == value {
			return true
		}
	}

	return false
}

func (m RequestMapping) matchesRequest(r *http.Request) bool {
	if m.Method != r.Method || strings.ToLower(m.Path) != strings.ToLower(r.URL.Path) {
		return false
	}

	lowerHeaders := toLowerHeaders(r)

	for key, values := range m.Headers {
		if reqVals, ok := lowerHeaders[strings.ToLower(key)]; ok {
			for _, val := range values {
				if !inSlice(reqVals, val) {
					return false
				}
			}
		} else {
			return false
		}
	}

	return true
}

func (m RequestMapping) matchesMapping(other RequestMapping) bool {
	if m.Method != other.Method || strings.ToLower(m.Path) != strings.ToLower(other.Path) {
		return false
	}

	return reflect.DeepEqual(toLowerMap(m.Headers), toLowerMap(other.Headers))
}

//Mapper :
type Mapper interface {
	Add(m RequestMapping, d ResponseData) error
	GetByRequest(r *http.Request) (ResponseData, error)
}

//DefaultMapper :
var DefaultMapper Mapper

type boltDbMapping struct {
	ID      int
	Request RequestMapping
}

type pathMapping map[string][]boltDbMapping

//BoltDBMapper :
type BoltDBMapper struct {
	mappings pathMapping
}

//Add :
func (mapper BoltDBMapper) Add(m RequestMapping, r ResponseData) error {
	// if mapper.Get(m) != nil
	lowerPath := strings.ToLower(m.Path)
	if _, ok := mapper.mappings[lowerPath]; !ok {
		mapper.mappings[lowerPath] = make([]boltDbMapping, 0)
	}

	for _, existed := range mapper.mappings[lowerPath] {
		if m.matchesMapping(existed.Request) {
			return errors.New("Mapping already exist in mapper")
		}
	}

	if id, err := AddResponseData(r); err == nil {
		mapper.mappings[lowerPath] = append(mapper.mappings[lowerPath], boltDbMapping{
			ID:      id,
			Request: m,
		})
	} else {
		return err
	}

	return nil
}

//GetByRequest :
func (mapper BoltDBMapper) GetByRequest(r *http.Request) (ResponseData, error) {
	if mappings, ok := mapper.mappings[strings.ToLower(r.URL.Path)]; ok {
		for _, mapping := range mappings {
			if !mapping.Request.matchesRequest(r) {
				continue
			}

			return GetResponseData(mapping.ID)
		}
	}

	return ResponseData{}, errors.New("Endpoint not found")
}

//MappingHandler :
func MappingHandler(w http.ResponseWriter, r *http.Request) {
	if rd, err := DefaultMapper.GetByRequest(r); err != nil {
		log.Printf("Failed to find request mapping for endpoint [%v]\n", r.URL)
		w.WriteHeader(http.StatusNotFound)
	} else {
		if rd.Status != 0 {
			w.WriteHeader(rd.Status)
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		for key, vavlues := range rd.Headers {
			for _, value := range vavlues {
				w.Header().Set(key, value)
			}
		}

		w.Write(rd.Data)
	}
}

func init() {
	DefaultMapper = BoltDBMapper{
		mappings: make(pathMapping, 0),
	}
}
