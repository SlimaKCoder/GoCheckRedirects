package main

import (
	"errors"
	"github.com/hashicorp/go-multierror"
	"github.com/remeh/sizedwaitgroup"
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

	yamlError := yaml.Unmarshal([]byte(fileData), &configs)

	if yamlError != nil {
		errorHandler(yamlError)
	}

	return configs
}

func loadRedirects(path string) []Redirect {
	var redirects []Redirect

	ffmt.Printf(">>> Reading redirects from: %s \n", path)

	// Loads everything into memory, it may be good idea to optimize it
	content, ioError := ioutil.ReadFile(path)
	lines := strings.Split(string(content), "\n")

	if ioError == nil {
		for _, line := range lines {
			if len(line) == 0 {
				break
			}

			urls := strings.Fields(line)

			redirect := Redirect{Source: urls[0], Dest: urls[1]}

			redirects = append(redirects, redirect)
		}
	} else {
		errorHandler(ioError)
	}

	return redirects
}

func checkRedirects(redirects []Redirect, url string) {
	var swg = sizedwaitgroup.New(4)

	ffmt.Printf(">>> Checking redirects for: %s \n", url)

	for index, redirect := range redirects {
		go checkRedirectAsync(redirect, url, index, &swg)
		swg.Add()
	}

	swg.Wait()
}

func checkRedirectAsync(redirect Redirect, url string, index int, swg *sizedwaitgroup.SizedWaitGroup) {
	defer swg.Done()
	checkRedirect(redirect, url, index)
}

func checkRedirect(redirect Redirect, url string, optionalArgs ...int) {
	var index int

	if optionalArgs != nil {
		index = optionalArgs[0]
	}

	response, httpError := http.Get(url)

	// TODO Implement redirects checking
	if httpError == nil {
		ffmt.Printf(
			"---- [%d] ---- \n" +
				"> Source: %s \n" +
				"> Dest: %s \n" +
				"> Actual: %s \n" +
				"> Status: %d \n",
			index,
			redirect.Source,
			redirect.Dest,
			redirect.Dest,
			response.StatusCode)
	} else {
		errorHandler(httpError)
	}
}