package main

import (
	"flag"
	"fmt"
	"os"
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
		fmt.Println("Invalid path")
		os.Exit(1)
	}

	if !loadProject() {
		fmt.Println("Failed to load project")
		os.Exit(1)
	}

	go waitToSave()

	initServer()
}
