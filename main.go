package main

import (
	"bufio"
	"errors"
	"github.com/hashicorp/go-multierror"
	"gopkg.in/ffmt.v1"
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
	ffmt.MarkStack(1, newError) // uses runtime.Caller(2)
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
		ffmt.Mark(errorsList)
		os.Exit(1)
	}
}

func loadConfigs(path string) []MapConfig {
	var configs []MapConfig = nil // Support up to 100 maps

	fileData, ioError := ioutil.ReadFile(path)

	if ioError != nil {
		errorHandler(ioError)
	}

	yaml.Unmarshal([]byte(fileData), &configs)

	return configs
}

func loadRedirects(path string) []Redirect {
	var redirects []Redirect

	ffmt.Printf("> Reading redirects from: %s \n", path)

	file, ioError := os.Open(path)
	defer file.Close()

	// TODO Finish implementation
	if ioError == nil {
		reader := bufio.NewReader(file)

		for {
			line, ioError := reader.ReadString('\n')

			if ioError != nil {
				errorHandler(ioError)
			}

			// Process the line here.
			//fmt.Println(" > > " + string(line))

			if len(line) == 0 {
				break
			}
		}
	} else {
		errorHandler(ioError)
	}

	return redirects
}

func checkMap(config MapConfig) {
	redirects := loadRedirects(config.Path)
	redirects = redirects

	// TODO Implement redirects checking
}