package main

import (
	"fmt"
	"os"
)

const (
	// Version of the application
	Version = "0.0.2"
)

var (
	pathProject string
)

func main() {
	fmt.Println("Zoomer Project v" + Version + " by ^[GS]^")

	if len(os.Args) < 2 {
		fmt.Println("Usage: zoomer <project path>")
		os.Exit(1)
	}
	pathProject = os.Args[1]

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
