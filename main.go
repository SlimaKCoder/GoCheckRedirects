package main

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Redirect struct {
	Source string
	Dest string
}

type MapConfig struct {
	Name string
	Domain string
	Path string
}

func main() {
	var configs = loadConfigs(`config.yml`)

	for _, config := range configs {
		checkMap(config)
	}

}

func loadConfigs(path string) []MapConfig {
	var configs []MapConfig = nil // Support up to 100 maps

	fileData, _ := ioutil.ReadFile(path)
	yaml.Unmarshal([]byte(fileData), &configs)

	return configs
}

func loadRedirects(path string) []Redirect {
	fmt.Printf("Reading file from: %s \n", path)

	file, io_error := os.Open(path)

	if io_error != nil {
		fmt.Printf("Error: %s\n", io_error)
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	var redirects []Redirect

	for {
		line, _ := reader.ReadString('\n')

		fmt.Printf(" > Read %d characters\n", len(line))

		// Process the line here.
		//fmt.Println(" > > " + string(line))

		if file == nil {
			break
		}
	}

	return redirects
}

func checkMap(config MapConfig) {
	redirects := loadRedirects(config.Path)
	redirects = redirects
}