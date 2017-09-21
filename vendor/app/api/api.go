package api

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

//Expects :
type Expects struct {
	Params  []Param  `yaml:"params"`
	Headers []Header `yaml:"headers"`
}

//ResponseData :
type ResponseData struct {
	Status           int      `yaml:"status"`
	Headers          []Header `yaml:"headers"`
	Data             []byte   `yaml:"data"`
	ResponseFilePath string   `yaml:"file"`
}

//Endpoint :
type Endpoint struct {
	Name     string
	Path     string       `yaml:"path"`
	Expects  Expects      `yaml:"expects"`
	Response ResponseData `yaml:"response"`
}
