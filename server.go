package main

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	listenPort = 80
)

func initServer() {
	srv := http.Server{
		Addr: fmt.Sprint(":", listenPort),
	}

	http.HandleFunc("/", handler)

	fmt.Println("Server is listening on port", listenPort)
	srv.ListenAndServe()
}

func parseEscapeHTML(data string) string {
	data = strings.ReplaceAll(data, "<", "&lt;")
	data = strings.ReplaceAll(data, ">", "&gt;")
	return data
}

func handler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "<h1>"+configProject.ProjectName+"</h1>")
	fmt.Fprintf(w, "<h2>Project path: "+pathProject+"</h2>")
	fmt.Fprintf(w, "<h3>Project files:</h3>")
	for filename, data := range fileData {
		fmt.Fprintf(w, "<h4>"+filename+"</h4>")
		fmt.Fprintf(w, "<pre>")
		fmt.Fprintf(w, parseEscapeHTML(data))
		fmt.Fprintf(w, "</pre>")
	}
}
