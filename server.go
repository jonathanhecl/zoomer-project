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
	headerHtml(w)
	fmt.Fprintf(w, "<h3>Project files:</h3>")
	for _, filename := range projectFiles {
		showSourceHtml(w, filename)
		break
	}
	footerHtml(w)
}

func headerHtml(w http.ResponseWriter) {
	fmt.Fprintf(w, `<html><head><title>`+configProject.ProjectName+`</title></head><body>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.8.0/styles/default.min.css">
		<script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.8.0/highlight.min.js"></script>
		<h1>`+configProject.ProjectName+`</h1>
		<span>Project path: `+pathProject+`</span>`)
}

func footerHtml(w http.ResponseWriter) {
	fmt.Fprintf(w, `<script>
		hljs.highlightAll();
		</script>
		</body></html>`)
}

func getFilename(filename string) string {
	return strings.ReplaceAll(filename, pathProject, "")
}

func showSourceHtml(w http.ResponseWriter, filename string) {
	fmt.Fprintf(w, `<h4>`+getFilename(filename)+`</h4><pre>`)
	if configProject.LangHighlight != "" {
		fmt.Fprintf(w, `<code class="`+configProject.LangHighlight+`">`)
	} else {
		fmt.Fprintf(w, `<code>`)
	}
	fmt.Fprintf(w, parseEscapeHTML(fileData[filename]))
	fmt.Fprintf(w, `</code></pre>`)
}
