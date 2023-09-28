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
	http.HandleFunc("/save", saveHandler)

	fmt.Println("Server is listening on port", listenPort)
	srv.ListenAndServe()
}

func handler(w http.ResponseWriter, r *http.Request) {
	headerHtml(w)
	fmt.Fprintf(w, "<h3>Project files:</h3>")
	for _, filepath := range projectFiles {
		//showFilelistHtml(w, filepath)
		showSourceHtml(w, filepath)
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
			h4 {
				font-size: xx-large;
				text-align: center;
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
			.collumns {
				display: flex; 
				flex-direction: row;	
			}
			.collumns > .codes {
				min-width: 400px;
			}
			.float-right {
				position: fixed;
				bottom: 10px;
				right: 20px;
			}
			.fields > .method {
				font-size: large;
				color: lime;
    			//text-align: center;
				//padding: 2em 0 0;
			}
			.field > label {
				border: 1px solid #ccc;
				padding: 0.6em;
				margin: 0;
				display: block;
				color: #ccc;
			}
			.field > label:hover {
				 background:#333;
				 cursor:pointer;
			}
			.field textarea {
				background-color: #333;
				color: #ddd;
				display: block;
    			width: 100%%;
			}
			</style>
		<body>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.8.0/styles/github-dark.min.css">
		<script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.8.0/highlight.min.js"></script>
		<h1 id="top">`+configProject.ProjectName+`</h1>
		<span>Project path: `+pathProject+`</span>
		<div class="float-right">
			`+getFilelistDropdownHtml()+`
			<a href="#top">Go Top</a>
		</div>`)
}

func footerHtml(w http.ResponseWriter) {
	fmt.Fprintf(w, `<script>
		hljs.highlightAll();

		function saveChange(obj) {
			var name = obj.name;
			var value = "";
			if (obj.type == "checkbox") {
				value = obj.checked ? 1 : 0;
			} else if (obj.type == "textarea") {
				value = obj.value;
			}
			var xhttp = new XMLHttpRequest();
			xhttp.open("POST", "/save", true);
			xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
			xhttp.send("name="+name+"&value="+value);	
		}

		</script>
		</body></html>`)
}

func getFilelistDropdownHtml() string {
	var html string = `<select onchange="location = this.value;">`
	for _, filepath := range projectFiles {
		filename := getFilename(filepath)
		html += `<option value="#` + getFileID(filename) + `">` + filename + `</option>`
	}
	html += `</select>`
	return html
}

func showFilelistHtml(w http.ResponseWriter, filepath string) {
	filename := getFilename(filepath)
	fmt.Fprintf(w, `<a href="#`+filename+`">`+filename+`</a><br>`)
}

func getFileID(filename string) string {
	return strings.ReplaceAll(filename, "/", ".")
}

func parseEscapeHTML(data string) string {
	data = strings.ReplaceAll(data, "<", "&lt;")
	data = strings.ReplaceAll(data, ">", "&gt;")

	return data
}

func showSourceHtml(w http.ResponseWriter, filepath string) {
	filename := getFilename(filepath)
	fmt.Fprintf(w, `<div id="`+getFileID(filename)+`" class="mark"></div>
						<h4>`+filename+`</h4>`)

	fmt.Fprintf(w, `<div class="collumns">`)
	fmt.Fprintf(w, `<div class="codes">`)
	fmt.Fprintf(w, filesData[filename].getContentHTMLWithFields())
	fmt.Fprintf(w, `</div></div>`)
}

func createFieldName(filename string, method string, field string) string {
	return filename + `<>` + method + `<>` + field
}

func disassemblyFieldName(fieldName string) (string, string, string) {
	fields := strings.Split(fieldName, `<>`)
	return fields[0], fields[1], fields[2]
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !changedUserField(r.Form.Get("name"), r.Form.Get("value")) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
