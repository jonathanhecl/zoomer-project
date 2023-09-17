package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type fieldsData struct {
	Field string
	Value string
}

type fileData struct {
	Content    []string
	Methods    []int
	UserFields map[string]fieldsData
}

var (
	projectFiles []string
	filesData    map[string]fileData
)

func (f fieldsData) getValue(field string) string {
	if f.Field == field {
		return f.Value
	}
	return ""
}

func (f fileData) getContent() string {
	return strings.Join(f.Content, "\n")
}

func (f fileData) getMethods() []string {
	methods := []string{}
	for _, method := range f.Methods {
		methods = append(methods, f.Content[method])
	}
	return methods
}

func loadProject() bool {
	fmt.Println("Project path:", pathProject)

	if !loadConfig() {
		return false
	}

	projectFiles = make([]string, 0)
	filesData = make(map[string]fileData)

	var err error

	projectFiles, err = scanProject(pathProject, nil)
	if err != nil {
		fmt.Println(err)
		return false
	}

	//fmt.Println(projectFiles)
	//fmt.Println(fileData)

	return true
}

func isExtFilter(filename string) bool {
	for _, ext := range configProject.ExtFilter {
		if filepath.Ext(filename) == ext {
			return true
		}
	}
	return false
}

func loadFileData(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	content := []string{}
	methods := []int{}

	for i, line := range strings.Split(string(data), "\n") {
		content = append(content, line)
		for _, method := range configProject.MethodFilter {
			re, _ := regexp.Compile(method)
			if re.MatchString(line) {
				methods = append(methods, i)
			}
		}
	}

	filesData[filename] = fileData{
		Content: content,
		Methods: methods,
	}

	return nil
}

func scanProject(root string, list []string) ([]string, error) {

	var filesOut []string = list
	var filesTmp []string

	files, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if !f.IsDir() {
			if isExtFilter(f.Name()) {
				filesOut = append(filesOut, path.Join(root+f.Name()))
				if err := loadFileData(path.Join(root + f.Name())); err != nil {
					return nil, err
				}
			}
		}
	}

	for _, f := range files {
		if f.IsDir() {
			filesTmp, err = scanProject(path.Join(root, f.Name())+string(filepath.Separator), filesOut)
			if err != nil {
				return nil, err
			}
			filesOut = filesTmp
		}
	}

	return filesOut, nil

}
