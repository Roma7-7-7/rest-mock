package main

//SupportedMethods :
var SupportedMethods = []string{
	"GET",
	"POST",
	"PUT",
	"PATCH",
	"HEAD",
	"DELETE",
}

//Header :
type Header struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

//Param :
type Param struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}
