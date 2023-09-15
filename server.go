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

func handler(w http.ResponseWriter, r *http.Request) {
	headerHtml(w)
	fmt.Fprintf(w, "<h3>Project files:</h3>")
	for _, filename := range projectFiles {
		//showFilelistHtml(w, filename)
		showSourceHtml(w, filename)
	}
	footerHtml(w)
}

func headerHtml(w http.ResponseWriter) {
	fmt.Fprintf(w, `
		<html data-theme="dark">
			<head>
			<title>`+configProject.ProjectName+`</title>
			</head>
			<style>
			body {
				background-color: #1e1e1e;
				color: #d4d4d4;
				font-family: monospace;
			}
			a {
				color: #d4d4d4;
			}
			a:hover {
				color: #d4d4d4;
				text-decoration: underline;
			}
			pre {
				background-color: #2d2d2d;
				padding: 10px;
				border-radius: 5px;
			}
			code {
				font-family: monospace;
			}
			.float-right {
				position: fixed;
				bottom: 10px;
				right: 20px;
			}
			</style>
		<body>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.8.0/styles/github-dark.min.css">
		<script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.8.0/highlight.min.js"></script>
		<h1 id="top">`+configProject.ProjectName+`</h1>
		<span>Project path: `+pathProject+`</span>
		<a href="#top" class="float-right">Go Top</a>`)
}

func footerHtml(w http.ResponseWriter) {
	fmt.Fprintf(w, `<script>
		hljs.highlightAll();
		</script>
		</body></html>`)
}

func showFilelistHtml(w http.ResponseWriter, filename string) {
	fmt.Fprintf(w, `<a href="/`+filename+`">`+getFilename(filename)+`</a><br>`)
}

func getFilename(filename string) string {
	return strings.ReplaceAll(filename, pathProject, "")
}

func parseEscapeHTML(data string) string {
	data = strings.ReplaceAll(data, "<", "&lt;")
	data = strings.ReplaceAll(data, ">", "&gt;")

	return data
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
