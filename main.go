package main

import (
	"errors"
	"github.com/hashicorp/go-multierror"
	"gopkg.in/ffmt.v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Redirect struct {
	Source string
	Dest string
}

type MapConfig struct {
	Url string
	Path string
}

var errorsList *multierror.Error

func errorHandler(newError error) {
	errorsList = multierror.Append(errorsList, newError)
	ffmt.MarkStack(1, newError) // uses runtime.Caller(2)
}

func main() {
	configs := loadConfigs(`config.yml`)

	if configs != nil {
		for _, config := range configs {
			redirects := loadRedirects(config.Path)
			checkRedirects(redirects, config.Url)
		}
	} else {
		errorHandler(
			errors.New("read config.yml: empty file"))
	}

	if errorsList != nil {
		ffmt.Mark(errorsList)
		os.Exit(1)
	}
}

func loadConfigs(path string) []MapConfig {
	var configs []MapConfig

	fileData, ioError := ioutil.ReadFile(path)

	if ioError != nil {
		errorHandler(ioError)
	}

	yaml.Unmarshal([]byte(fileData), &configs)

	return configs
}

func loadRedirects(path string) []Redirect {
	var redirects []Redirect

	ffmt.Printf(">>> Reading redirects from: %s \n", path)

	// Loads everything into memory, it may be good idea to optimize it
	content, ioError := ioutil.ReadFile(path)
	lines := strings.Split(string(content), "\n")

	if ioError == nil {
		for _, line := range lines{
			urls := strings.Fields(line)

			redirect := Redirect{
				Source: urls[0],
				Dest: urls[1],
			}

			redirects = append(redirects, redirect)

			if len(line) == 0 {
				break
			}
		}
	} else {
		errorHandler(ioError)
	}

	return redirects
}



func checkRedirects(redirects []Redirect, url string) {
	ffmt.Printf(">>> Checking redirects for: %s \n", url)

	// TODO Implement redirects checking
	for i, redirect := range redirects {
		response, httpError := http.Get(url)

		if httpError == nil {
			ffmt.Printf("---- [%d] ---- \n", i)
			ffmt.Printf("> Source: %s \n", redirect.Source)
			ffmt.Printf("> Dest: %s \n", redirect.Dest)
			ffmt.Printf("> Actual: %s \n", redirect.Dest)
			ffmt.Printf("> Status: %s \n", response.Status)
		} else {
			errorHandler(httpError)
		}
	}

}
	//if ioError != nil {
	//	errorHandler(
	//		errors.New(
	//			strings.Join(
	//				[]string{"read ", path, ": ", ioError.Error()}, "")))
	//}