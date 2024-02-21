package main

import (
	"flag"
	"fmt"
)

const (
	// Version of the application
	Version = "0.0.5"
)

var (
	pathProject string
)

func main() {
	fmt.Println("Zoomer Project v" + Version + " by ^[GS]^")

	flag.StringVar(&pathProject, "path", "", "project path")
	flag.Parse()

	if pathProject == "" {
		fmt.Println("Usage: zoomer --path <project path>")
		return
	}

	if !isValidPath(pathProject) {
		fmt.Println(pathProject + " is not a valid path")
		return
	}

	if !loadProject() {
		fmt.Println("Error loading project")
		return
	}

	go waitToSave()

	initServer()
}
