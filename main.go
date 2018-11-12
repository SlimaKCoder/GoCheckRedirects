package main

import (
	"bufio"
	"errors"
	"github.com/hashicorp/go-multierror"
	fmt "gopkg.in/ffmt.v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Redirect struct {
	Source string
	Dest string
}

type MapConfig struct {
	Domain string
	Path string
}

var errorsList *multierror.Error

func errorHandler(newError error) {
	errorsList = multierror.Append(errorsList, newError)
	//_, _, line, _ := runtime.Caller(1)

	fmt.MarkStack(1, newError) // good, position of errorHandler func call
	fmt.Mark(newError)              // bad, current position
}

func main() {
	var configs = loadConfigs(`config.yml`)

	if configs != nil {
		for _, config := range configs {
			checkMap(config)
		}
	} else {
		errorHandler(errors.New("read config.yml: empty file"))
	}

	if errorsList != nil {
		fmt.Mark(errorsList)
		os.Exit(1)
	}
}

func loadConfigs(path string) []MapConfig {
	var configs []MapConfig = nil // Support up to 100 maps

	fileData, ioError := ioutil.ReadFile(path)

	if ioError != nil {
		errorsList = multierror.Append(errorsList, ioError)
		fmt.Mark(ioError)
	}

	yaml.Unmarshal([]byte(fileData), &configs)

	return configs
}

func loadRedirects(path string) []Redirect {
	var redirects []Redirect

	fmt.Printf("Reading redirects from: %s \n", path)

	file, ioError := os.Open(path)

	if ioError != nil {
		errorsList = multierror.Append(errorsList, ioError)
		fmt.Mark(ioError)
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, ioError := reader.ReadString('\n')

		if ioError != nil {
			errorsList = multierror.Append(errorsList, ioError)
			fmt.Mark(ioError)
		}

		// Process the line here.
		//fmt.Println(" > > " + string(line))

		if len(line) == 0 {
			break
		}
	}

	return redirects
}

func checkMap(config MapConfig) {
	redirects := loadRedirects(config.Path)
	redirects = redirects
}