package main

import (
	"flag"
	"fmt"
)

const (
	// Version of the application
	Version = "0.0.7"
)

var (
	pathProject string
	listenPort  = "80"
)

func main() {
	fmt.Println("Zoomer Project v" + Version + " by ^[GS]^")

	flag.StringVar(&pathProject, "path", "", "project path")
	flag.StringVar(&listenPort, "port", "80", "port to listen")
	flag.Parse()

	if pathProject == "" || !isValidPath(pathProject) || !isValidPort(listenPort) {
		fmt.Println("Usage: zoomer --path <project path> [--port <port>]")
		return
	}

	if !loadProject() {
		fmt.Println("Error loading project")
		return
	}

	go waitToSave()

	initServer()
}
