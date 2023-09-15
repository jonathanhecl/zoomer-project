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
	fmt.Fprintf(w, "<html><head><title>"+configProject.ProjectName+"</title></head><body>")
	fmt.Fprintf(w, "<link rel=\"stylesheet\" href=\"https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.8.0/styles/default.min.css\">\n<script src=\"https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.8.0/highlight.min.js\"></script>")
	fmt.Fprintf(w, "<h1>"+configProject.ProjectName+"</h1>")
	fmt.Fprintf(w, "<h2>Project path: "+pathProject+"</h2>")
	fmt.Fprintf(w, "<h3>Project files:</h3>")
	for _, filename := range projectFiles {
		fmt.Fprintf(w, "<h4>"+filename+"</h4>")
		fmt.Fprintf(w, "<pre><code>")
		fmt.Fprintf(w, parseEscapeHTML(fileData[filename]))
		fmt.Fprintf(w, "</code></pre>")
	}
	fmt.Fprintf(w, "<script>hljs.highlightAll();</script>")
	fmt.Fprintf(w, "</body></html>")
}
