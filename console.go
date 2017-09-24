package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

//Port :
var Port = 8888

//AddMapping :
var AddMapping = false

var reader = bufio.NewReader(os.Stdin)

type consoleInput struct {
	Key   string
	Value string
}

var consoleInputs = make([]consoleInput, 0)

//CheckTrue :
func CheckTrue(header string) bool {
	for {
		fmt.Printf("%v [y/n]\n", header)

		if line := strings.ToLower(readString()); line != "y" && line != "n" {
			fmt.Printf("Invalid value [%v]\n", line)
		} else {
			return line == "y"
		}
	}
}

//AddRequestMapping :
func AddRequestMapping() {
	for {
		method := choose("Select method:", consoleInputs)
		fmt.Printf("Selected method: %v\n", method)
		path := readPath()
		fmt.Printf("Selected path: %v\n", path)
		headers := getHeaders()
		fmt.Println(headers)

		result := RequestMapping{
			Method:  method,
			Path:    path,
			Headers: headers,
			Params:  make(map[string][]string),
		}

		if err := DefaultMapper.Add(result, ResponseData{Status: 201, Data: []byte("Boltdb")}); err != nil {
			log.Println(err)
			log.Println("Failed to add request mapping")
		}

		if !CheckTrue("Do you want to add additional mapping?") {
			return
		}
	}
}

func readString() string {
	line, _ := reader.ReadString('\n')
	return strings.Replace(line, "\n", "", 1)
}

func choose(header string, variants []consoleInput) string {
	for {
		fmt.Println(header)
		for _, i := range variants {
			fmt.Printf("%v. %v\n", i.Key, i.Value)
		}

		line := readString()
		for _, v := range variants {
			if v.Key == line {
				return v.Value
			}
		}

		fmt.Printf("Invalud value [%v]\n", line)
	}
}

func readPath() string {
	for {
		fmt.Println("Please enter path to endpoint")

		if line := readString(); strings.HasPrefix(line, "/") {
			return line
		}

		fmt.Println("Path should start with '/'")

	}
}

func getHeaders() map[string][]string {
	result := make(map[string][]string)

	for {
		if !CheckTrue("Do you want to add header?") {
			if len(result) == 0 {
				fmt.Println("Skip headers")
			}
			return result
		}

		fmt.Println("Enter key")
		key := readString()
		fmt.Println("Enter value")
		value := readString()

		if CheckTrue(fmt.Sprintf("Is it correct header?\n%v: %v", key, value)) {
			lowerKey := strings.ToLower(key)
			if values, ok := result[lowerKey]; ok {
				result[lowerKey] = append(values, value)
			} else {
				result[lowerKey] = []string{value}
			}
		}
	}
}

func init() {

	log.Println("Initialising application flags")

	flag.IntVar(&Port, "p", 8888, "Port that will be used to run application")
	flag.BoolVar(&AddMapping, "add", false, "Add new endpoint")

	flag.Parse()

	for i, value := range SupportedMethods {
		consoleInputs = append(consoleInputs, consoleInput{
			Key:   strconv.Itoa(i + 1),
			Value: value,
		})
	}
}
