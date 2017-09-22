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
	Method   string
	Path     string
	Headers  map[string][]string
	Params   map[string][]string
	Response ResponseData
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
	Add(m RequestMapping) error
	Get(r *http.Request) (*ResponseData, error)
}

//DictMapper :
type DictMapper struct {
	mappings map[string][]RequestMapping
}

//Add :
func (d DictMapper) Add(m RequestMapping) error {
	lowerPath := strings.ToLower(m.Path)
	if mappings, ok := d.mappings[lowerPath]; !ok {
		d.mappings[lowerPath] = []RequestMapping{m}
	} else {
		for _, existed := range mappings {
			if m.matchesMapping(existed) {
				return errors.New("Mapping already exist in mapper")
			}
		}

		d.mappings[lowerPath] = append(d.mappings[lowerPath], m)
	}
	return nil
}

//Get :
func (d DictMapper) Get(r *http.Request) (*ResponseData, error) {
	if mappings, ok := d.mappings[strings.ToLower(r.URL.Path)]; ok {
		for _, mapping := range mappings {
			if mapping.matchesRequest(r) {
				return &mapping.Response, nil
			}
		}
	}

	return nil, errors.New("Endpoint not found")
}

//DefaultMapper :
var DefaultMapper Mapper

//MappingHandler :
func MappingHandler(w http.ResponseWriter, r *http.Request) {
	if rd, err := DefaultMapper.Get(r); rd == nil || err != nil {
		log.Printf("Failed to find request mapping for endpoint [%v]\n", r.URL)
		w.WriteHeader(http.StatusNotFound)
	} else {
		mapStatus(w, rd)
		mapHeaders(w, rd)
		w.Write(rd.Data)
	}
}

func mapStatus(w http.ResponseWriter, r *ResponseData) {
	if r.Status == 0 {
		return
	}

	w.WriteHeader(r.Status)
}

func mapHeaders(w http.ResponseWriter, r *ResponseData) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	for key, vavlues := range r.Headers {
		for _, value := range vavlues {
			w.Header().Set(key, value)
		}
	}
}

func init() {
	DefaultMapper = DictMapper{
		mappings: make(map[string][]RequestMapping, 0),
	}

	DefaultMapper.Add(RequestMapping{
		Method:  "GET",
		Path:    "/",
		Headers: make(map[string][]string),
		Params:  make(map[string][]string),
		Response: ResponseData{
			Status:  201,
			Headers: make(map[string][]string),
			Data:    []byte("Success"),
		},
	})

	headers1 := make(map[string][]string)
	headers1["Content-Type"] = []string{"application/json"}

	DefaultMapper.Add(RequestMapping{
		Method:  "POST",
		Path:    "/test2",
		Headers: headers1,
		Params:  make(map[string][]string),
		Response: ResponseData{
			Status:  201,
			Headers: make(map[string][]string),
			Data:    []byte("Success 2"),
		},
	})
}
