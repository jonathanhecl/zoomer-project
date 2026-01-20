package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"regexp"
)

const (
	configFilename     = "zoomer-config.json"
	userFieldsFilename = "zoomer-userfields.json"
)

var (
	configProject        config
	methodFilterRegexes []*regexp.Regexp
)

type EnumFieldType string

const (
	EnumTextBox EnumFieldType = "textbox"
	EnumBoolean EnumFieldType = "boolean"
)

type UserField struct {
	Name string
	Type EnumFieldType
}

type config struct {
	ProjectName   string      `json:"project_name"`
	LangHighlight string      `json:"lang_highlight"`
	ExtFilter     []string    `json:"ext_filter"`
	MethodFilter  []string    `json:"method_filter"`
	UserFields    []UserField `json:"user_fields"`
}

func createConfig() bool {
	fmt.Println("Creating config file:", path.Join(pathProject, configFilename))
	newConfig := config{
		ProjectName:   "New Project",
		LangHighlight: "go",
		ExtFilter:     []string{".go"},
		MethodFilter:  []string{"func (\\(.*\\))?(.*)\\(.*?\\).*{"},
		UserFields:    []UserField{{"Checked", EnumBoolean}},
	}

	configFile, err := os.Create(path.Join(pathProject, configFilename))
	if err != nil {
		fmt.Printf("Error creating config file: %v\n", err)
		return false
	}
	defer configFile.Close()

	encoder := json.NewEncoder(configFile)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(newConfig)
	if err != nil {
		fmt.Printf("Error encoding config: %v\n", err)
		return false
	}

	fmt.Println("Config file created successfully")
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

	// Precompilar expresiones regulares para mejor rendimiento
	methodFilterRegexes = make([]*regexp.Regexp, 0, len(configProject.MethodFilter))
	for _, pattern := range configProject.MethodFilter {
		re, err := regexp.Compile(pattern)
		if err != nil {
			fmt.Printf("Warning: invalid regex pattern '%s': %v\n", pattern, err)
			continue
		}
		methodFilterRegexes = append(methodFilterRegexes, re)
	}

	fmt.Println("Config loaded! (", configProject.ProjectName, ")")
	return true
}
