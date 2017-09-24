package main

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

//RequestMapping :
type RequestMapping struct {
	Method  string
	Path    string
	Headers map[string][]string
	Params  map[string][]string
}

//ResponseData :
type ResponseData struct {
	Status  int                 `yaml:"status"`
	Headers map[string][]string `yaml:"headers"`
	Data    []byte              `yaml:"data"`
}

func (m RequestMapping) matchesRequest(r *http.Request) bool {
	if m.Method != r.Method || strings.ToLower(m.Path) != strings.ToLower(r.URL.Path) {
		return false
	}

	return comapreHeaders(m.Headers, r.Header)
}

func (m RequestMapping) matchesMapping(other RequestMapping) bool {
	if m.Method != other.Method || strings.ToLower(m.Path) != strings.ToLower(other.Path) {
		return false
	}

	return comapreHeaders(m.Headers, other.Headers)
}

//BoltDbItem :
type BoltDbItem struct {
	Request  RequestMapping
	Response ResponseData
}

//Mapper :
type Mapper interface {
	Add(m RequestMapping, d ResponseData) error
	GetByRequest(r *http.Request) (*ResponseData, error)
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
	log.Printf("Adding new mapping for method %v and path %v\n", m.Method, m.Path)

	lowerPath := strings.ToLower(m.Path)
	if _, ok := mapper.mappings[lowerPath]; !ok {
		mapper.mappings[lowerPath] = make([]boltDbMapping, 0)
	}

	for _, existed := range mapper.mappings[lowerPath] {
		if m.matchesMapping(existed.Request) {
			return errors.New("Mapping already exist in mapper")
		}
	}

	log.Println("Adding boltdb item")
	item := BoltDbItem{Request: m, Response: r}
	if id, err := AddResponseData(item); err == nil {
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
func (mapper BoltDBMapper) GetByRequest(r *http.Request) (*ResponseData, error) {
	log.Printf("Looking for mapping for method %v and path %v\n", r.Method, r.URL.Path)
	if mappings, ok := mapper.mappings[strings.ToLower(r.URL.Path)]; ok {
		for _, mapping := range mappings {
			if !mapping.Request.matchesRequest(r) {
				continue
			}

			log.Printf("Looking for mapping with id %v\n", mapping.ID)
			result, err := GetBoltDbItem(mapping.ID)

			if err != nil {
				return nil, err
			}
			return &result.Response, nil
		}
	}

	return nil, errors.New("Endpoint not found")
}

func newBoltDBMapper(items map[int]*BoltDbItem) BoltDBMapper {
	log.Println("Creating new instance of Bolt DB Mapper")
	result := BoltDBMapper{mappings: make(pathMapping, 0)}

	for k, v := range items {
		log.Printf("Adding mapping with method %v and path %v to default mapper\n", v.Request.Method, v.Request.Path)
		lowerPath := strings.ToLower(v.Request.Path)
		value := boltDbMapping{
			ID:      k,
			Request: v.Request,
		}

		if _, ok := result.mappings[lowerPath]; !ok {
			result.mappings[lowerPath] = []boltDbMapping{value}
		} else {
			result.mappings[lowerPath] = append(result.mappings[lowerPath], value)
		}

		log.Printf("Added mapping with id %v method %v and path %v to default mapper\n", value.ID, value.Request.Method, value.Request.Path)
	}

	return result
}

//MappingHandler :
func MappingHandler(w http.ResponseWriter, r *http.Request) {
	if rd, err := DefaultMapper.GetByRequest(r); err != nil {
		log.Printf("Failed to find request mapping for endpoint [%v]\n", r.URL)
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
	} else {
		if rd.Status != 0 {
			w.WriteHeader(rd.Status)
		}

		w.Write(rd.Data)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		for key, vavlues := range rd.Headers {
			for _, value := range vavlues {
				w.Header().Set(key, value)
			}
		}
	}
}

func comapreHeaders(h1 map[string][]string, h2 map[string][]string) bool {
	h1, h2 = toLowerHeaders(h1), toLowerHeaders(h2)

	for key, values1 := range h1 {
		if values2, ok := h2[key]; ok {
			for _, v := range values1 {
				if !inSlice(values2, v) {
					return false
				}
			}
		} else {
			return false
		}
	}

	return true
}

func toLowerHeaders(h map[string][]string) map[string][]string {
	result := make(map[string][]string)

	for key, value := range h {
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

func init() {
	if items, err := GetAll(); err != nil {
		log.Fatalf("Failed to initialise mapper\n%v", err)
	} else {
		DefaultMapper = newBoltDBMapper(items)
	}
}
