package main

import (
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"
)

func initServer() {
	srv := http.Server{
		Addr:         fmt.Sprint(":", listenPort),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	http.HandleFunc("/", handler)
	http.HandleFunc("/save", saveHandler)

	fmt.Println("Server is listening on port", listenPort)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Server error: %v\n", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	headerHtml(w)
	fmt.Fprintf(w, `<div class="content">`)
	fmt.Fprintf(w, "<h3>📄 Project Files</h3>")
	for _, filepath := range projectFiles {
		showSourceHtml(w, filepath)
	}
	fmt.Fprintf(w, `</div></div class="container">`)
	footerHtml(w)
}

func headerHtml(w http.ResponseWriter) {
	fmt.Fprintf(w, `
	<!DOCTYPE html>
		<html data-theme="dark">
			<head>
			<meta charset="UTF-8">
			<title>`+configProject.ProjectName+`</title>
			</head>
			<style>
			* {
				box-sizing: border-box;
			}
			
			body {
				background: linear-gradient(135deg, #1a1a2e 0%%, #16213e 50%%, #0f3460 100%%);
				color: #e4e4e4;
				font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', sans-serif;
				line-height: 1.6;
				margin: 0;
				padding: 20px;
				min-height: 100vh;
			}
			
			.container {
				max-width: 1400px;
				margin: 0 auto;
			}
			
			header {
				background: rgba(255, 255, 255, 0.05);
				backdrop-filter: blur(10px);
				border-radius: 12px;
				padding: 24px 32px;
				margin-bottom: 30px;
				box-shadow: 0 4px 6px rgba(0, 0, 0, 0.3);
				border: 1px solid rgba(255, 255, 255, 0.1);
			}
			
			h1 {
				font-size: 2.5em;
				font-weight: 700;
				margin: 0 0 10px 0;
				background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
				-webkit-background-clip: text;
				-webkit-text-fill-color: transparent;
				background-clip: text;
			}
			
			header span {
				color: #a0a0a0;
				font-size: 0.95em;
				font-family: 'Courier New', monospace;
			}
			
			h3 {
				font-size: 1.8em;
				margin: 30px 0 20px 0;
				color: #fff;
				font-weight: 600;
			}
			
			h4 {
				font-size: 1.8em;
				text-align: center;
				margin: 40px 0 25px 0;
				padding: 15px;
				background: rgba(255, 255, 255, 0.05);
				border-radius: 8px;
				border-left: 4px solid #667eea;
				color: #fff;
				font-weight: 600;
			}
			
			.mark {
				scroll-margin-top: 100px;
			}
			
			.file-section {
				background: rgba(255, 255, 255, 0.03);
				border-radius: 12px;
				padding: 25px;
				margin-bottom: 40px;
				box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
				border: 1px solid rgba(255, 255, 255, 0.08);
				transition: all 0.3s ease;
			}
			
			.file-section:hover {
				background: rgba(255, 255, 255, 0.05);
				box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
			}
			
			.collumns {
				display: flex;
				flex-direction: column;
				gap: 20px;
			}
			
			.codes {
				width: 100%%;
			}
			
			pre {
				background: #1e1e1e;
				padding: 20px;
				border-radius: 8px;
				overflow-x: auto;
				border: 1px solid rgba(255, 255, 255, 0.1);
				box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.3);
				margin: 0 0 20px 0;
			}
			
			code {
				font-family: 'Fira Code', 'Consolas', 'Monaco', 'Courier New', monospace;
				font-size: 0.95em;
				line-height: 1.6;
			}
			
			.fields {
				background: rgba(102, 126, 234, 0.1);
				border-radius: 8px;
				padding: 20px;
				margin-top: 20px;
				border: 1px solid rgba(102, 126, 234, 0.3);
			}
			
			.fields > .method {
				font-size: 1.1em;
				color: #4ade80;
				font-weight: 600;
				margin-bottom: 15px;
				padding: 10px;
				background: rgba(74, 222, 128, 0.1);
				border-radius: 6px;
				border-left: 4px solid #4ade80;
				font-family: 'Fira Code', 'Consolas', monospace;
				word-break: break-all;
			}
			
			.field {
				margin-bottom: 15px;
			}
			
			.field > label {
				display: block;
				margin-bottom: 8px;
				color: #d4d4d4;
				font-weight: 500;
				font-size: 0.95em;
			}
			
			.field input[type="checkbox"] {
				margin-right: 8px;
				width: 18px;
				height: 18px;
				cursor: pointer;
				accent-color: #667eea;
			}
			
			.field > label:has(input[type="checkbox"]) {
				padding: 12px 16px;
				background: rgba(255, 255, 255, 0.05);
				border: 1px solid rgba(255, 255, 255, 0.1);
				border-radius: 6px;
				cursor: pointer;
				transition: all 0.2s ease;
				display: flex;
				align-items: center;
			}
			
			.field > label:has(input[type="checkbox"]):hover {
				background: rgba(102, 126, 234, 0.15);
				border-color: rgba(102, 126, 234, 0.4);
				transform: translateX(3px);
			}
			
			.field textarea {
				background-color: rgba(30, 30, 30, 0.8);
				color: #e4e4e4;
				display: block;
				width: 100%%;
				padding: 12px;
				border: 1px solid rgba(255, 255, 255, 0.2);
				border-radius: 6px;
				font-family: inherit;
				font-size: 0.95em;
				resize: vertical;
				min-height: 80px;
				transition: all 0.2s ease;
			}
			
			.field textarea:focus {
				outline: none;
				border-color: #667eea;
				box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.2);
				background-color: rgba(30, 30, 30, 0.95);
			}
			
			.float-right {
				position: fixed;
				bottom: 30px;
				right: 30px;
				display: flex;
				flex-direction: column;
				gap: 12px;
				z-index: 1000;
			}
			
			.float-right select {
				background: rgba(102, 126, 234, 0.9);
				color: white;
				border: none;
				padding: 12px 16px;
				border-radius: 8px;
				font-size: 0.95em;
				cursor: pointer;
				box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
				transition: all 0.2s ease;
				min-width: 200px;
				font-weight: 500;
			}
			
			.float-right select:hover {
				background: rgba(102, 126, 234, 1);
				box-shadow: 0 6px 16px rgba(102, 126, 234, 0.4);
				transform: translateY(-2px);
			}
			
			.float-right select:focus {
				outline: none;
				box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.3);
			}
			
			.go-top-btn {
				background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
				color: white;
				text-decoration: none;
				padding: 12px 20px;
				border-radius: 8px;
				font-weight: 600;
				box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
				transition: all 0.2s ease;
				text-align: center;
				display: inline-block;
			}
			
			.go-top-btn:hover {
				transform: translateY(-3px);
				box-shadow: 0 6px 16px rgba(102, 126, 234, 0.4);
				text-decoration: none;
			}
			
			a {
				color: #667eea;
				text-decoration: none;
				transition: color 0.2s ease;
			}
			
			a:hover {
				color: #764ba2;
				text-decoration: underline;
			}
			
			/* Scrollbar personalizado */
			::-webkit-scrollbar {
				width: 12px;
				height: 12px;
			}
			
			::-webkit-scrollbar-track {
				background: rgba(0, 0, 0, 0.2);
			}
			
			::-webkit-scrollbar-thumb {
				background: rgba(102, 126, 234, 0.5);
				border-radius: 6px;
			}
			
			::-webkit-scrollbar-thumb:hover {
				background: rgba(102, 126, 234, 0.7);
			}
			
			/* Responsive */
			@media (max-width: 768px) {
				body {
					padding: 15px;
				}
				
				h1 {
					font-size: 2em;
				}
				
				h4 {
					font-size: 1.4em;
				}
				
				.float-right {
					bottom: 15px;
					right: 15px;
					left: 15px;
					flex-direction: row;
				}
				
				.float-right select {
					flex: 1;
				}
				
				.file-section {
					padding: 15px;
				}
			}
			</style>
		<body>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.8.0/styles/github-dark.min.css">
		<link rel="preconnect" href="https://fonts.googleapis.com">
		<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
		<link href="https://fonts.googleapis.com/css2?family=Fira+Code:wght@400;500;600&display=swap" rel="stylesheet">
		<script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.8.0/highlight.min.js"></script>
		<div class="container">
			<header>
				<h1 id="top">`+parseEscapeHTML(configProject.ProjectName)+`</h1>
				<span>📁 `+parseEscapeHTML(pathProject)+`</span>
			</header>
			<div class="float-right">
				`+getFilelistDropdownHtml()+`
				<a href="#top" class="go-top-btn">⬆️ Go Top</a>
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
		html += `<option value="#` + getFileID(filename) + `">` + parseEscapeHTML(filename) + `</option>`
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
	return html.EscapeString(data)
}

func showSourceHtml(w http.ResponseWriter, filepath string) {
	filename := getFilename(filepath)
	fmt.Fprintf(w, `<div id="`+getFileID(filename)+`" class="mark"></div>
						<h4>`+parseEscapeHTML(filename)+`</h4>`)

	fmt.Fprintf(w, `<div class="collumns">`)
	fmt.Fprintf(w, `<div class="codes">`)
	fmt.Fprintf(w, filesData[filename].getContentHTMLWithFields())
	fmt.Fprintf(w, `</div></div>`)
}

func createFieldName(filename string, method string, field string) string {
	// Los valores no deben ser escapados aquí porque se usan para identificación
	// El escape se hace donde se muestran en el HTML
	return filename + `<>` + method + `<>` + field
}

func disassemblyFieldName(fieldName string) (string, string, string) {
	fields := strings.Split(fieldName, `<>`)
	if len(fields) != 3 {
		return "", "", ""
	}
	return fields[0], fields[1], fields[2]
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	name := r.Form.Get("name")
	value := r.Form.Get("value")

	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !changedUserField(name, value) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
