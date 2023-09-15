package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

var (
	projectFiles []string
	fileData     map[string]string
)

func loadProject() bool {
	fmt.Println("Project path:", pathProject)

	if !loadConfig() {
		return false
	}

	projectFiles = make([]string, 0)
	fileData = make(map[string]string, 0)

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

	fileData[filename] = string(data)

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
