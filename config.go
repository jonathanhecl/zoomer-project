package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

const (
	configFilename = "zoomer-config.json"
)

var (
	configProject config
)

type config struct {
	ProjectName  string   `json:"project_name"`
	ExtFilter    []string `json:"ext_filter"`
	MethodFilter []string `json:"method_filter"`
}

func createConfig() bool {
	fmt.Println("Creating config file")
	newConfig := config{
		ProjectName:  "New Project",
		ExtFilter:    []string{".go"},
		MethodFilter: []string{"func (\\(.*\\))?(.*)\\(.*?\\).*{"},
	}

	configFile, err := os.Create(path.Join(pathProject, configFilename))
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer configFile.Close()

	encoder := json.NewEncoder(configFile)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(newConfig)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func loadConfig() bool {
	if !isValidFile(path.Join(pathProject, configFilename)) {
		if !createConfig() {
			fmt.Println("Failed to create config file")
			return false
		}
		return false
	}

	configFile, err := os.Open(path.Join(pathProject, configFilename))
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&configProject)
	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println("Config loaded! (", configProject.ProjectName, ")")
	return true
}
