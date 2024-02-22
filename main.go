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
	listenPort  = "80"
)

func main() {
	fmt.Println("Zoomer Project v" + Version + " by ^[GS]^")

	flag.StringVar(&pathProject, "path", "", "project path")
	flag.StringVar(&listenPort, "port", "80", "port to listen")
	flag.Parse()

	if pathProject == "" {
		fmt.Println("Usage: zoomer --path <project path> [--port <port>]")
		return
	}

	// check if port is valid
	if !isValidPort(listenPort) {
		fmt.Println(listenPort + " is not a valid port")
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
