package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"
)

type fileData struct {
	Filename string
	Content  []string
	Methods  []int
}

var (
	projectFiles []string
	filesData    map[string]fileData
	userFields   []fieldsData
	lastChange   time.Time
	lastSave     time.Time
)

func (f fileData) getContentHTMLWithFields() string {
	var content string = ""
	var prevMethod int = 0
	for _, method := range f.Methods {
		content += `<pre>`
		if configProject.LangHighlight != "" {
			content += `<code class="` + configProject.LangHighlight + `">`
		} else {
			content += `<code>`
		}
		content += parseEscapeHTML(strings.Join(f.Content[prevMethod:method], "\n"))
		content += `</code></pre>`

		if len(configProject.UserFields) > 0 {
			content += `<div class="fields">`
			content += `<div class="method">` + parseEscapeHTML(f.Content[method]) + `</div><br>`
			for _, field := range configProject.UserFields {
				content += `<div class="field">`
				fieldNameEscaped := parseEscapeHTML(field.Name)
				fieldNameAttr := parseEscapeHTML(createFieldName(f.Filename, f.Content[method], field.Name))
				if field.Type == EnumBoolean {
					content += `<label><input type="checkbox" name="` + fieldNameAttr + `" value="` + fieldNameEscaped + `" `
					if getUserValue(f.Filename, f.Content[method], field.Name) == "1" {
						content += `checked`
					}
					content += ` onchange="saveChange(this)"> ` + fieldNameEscaped + `</label>`
				} else if field.Type == EnumTextBox {
					content += `<label>` + fieldNameEscaped + `<br/><textarea name="` + fieldNameAttr + `" onchange="saveChange(this)">`
					content += parseEscapeHTML(getUserValue(f.Filename, f.Content[method], field.Name))
					content += `</textarea></label>`
				}
				content += `</div>`
			}

			content += `</div>`
		}
		prevMethod = method
	}

	return content
}

func (f fileData) getContent() string {
	return strings.Join(f.Content, "\n")
}

func getFilename(filepath string) string {
	filename := strings.ReplaceAll(filepath, pathProject, "")
	filename = strings.ReplaceAll(filename, "\\", "/")
	return filename
}

func loadProject() bool {
	fmt.Println("Project path:", pathProject)

	if !loadConfig() {
		return false
	}

	lastChange = time.Now()
	lastSave = lastChange

	loadUserFields()

	projectFiles = make([]string, 0)
	filesData = make(map[string]fileData)

	var err error

	projectFiles, err = scanProject(pathProject, nil)
	if err != nil {
		fmt.Printf("Error scanning project: %v\n", err)
		return false
	}

	fmt.Printf("Project loaded: %d file(s) found\n", len(projectFiles))
	return true
}

func isExtFilter(filename string) bool {
	for _, ext := range configProject.ExtFilter {
		if filepath.Ext(filename) == ext {
			return true
		}
	}
	return false
}

func loadFileData(filename string) error {
	const maxFileSize = 10 * 1024 * 1024 // 10MB limit

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	if stat.Size() > maxFileSize {
		return fmt.Errorf("file too large: %s (size: %d, max: %d)", filename, stat.Size(), maxFileSize)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	fileString := string(data)

	charsetUFT8 := utf8.ValidString(fileString)
	charsetANSI := !charsetUFT8

	if charsetANSI {
		fileString = fromWindows1252(fileString)
	}

	content := []string{}
	methods := []int{}

	for i, line := range strings.Split(fileString, "\n") {
		content = append(content, line)
		for _, re := range methodFilterRegexes {
			if re.MatchString(line) {
				methods = append(methods, i)
				break // Solo necesitamos que coincida con un patrón
			}
		}
	}

	filesData[getFilename(filename)] = fileData{
		Filename: filename,
		Content:  content,
		Methods:  methods,
	}

	return nil
}

func scanProject(root string, list []string) ([]string, error) {
	var filesOut []string = list
	var filesTmp []string

	files, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if !f.IsDir() {
			if isExtFilter(f.Name()) {
				filesOut = append(filesOut, path.Join(root+f.Name()))
				if err := loadFileData(path.Join(root + f.Name())); err != nil {
					return nil, err
				}
			}
		}
	}

	for _, f := range files {
		if f.IsDir() {
			filesTmp, err = scanProject(path.Join(root, f.Name())+string(filepath.Separator), nil)
			if err != nil {
				return nil, err
			}
			filesOut = append(filesOut, filesTmp...)
		}
	}

	return filesOut, nil
}

func changedUserField(name string, value string) bool {
	filename, method, field := disassemblyFieldName(name)
	if filename == "" || method == "" || field == "" {
		fmt.Printf("Invalid field name format: %s\n", name)
		return false
	}

	setUserValue(filename, method, field, value)
	lastChange = time.Now()

	return true
}

func fromWindows1252(str string) string {
	var arr = []byte(str)
	var buf bytes.Buffer
	var r rune

	for _, b := range arr {
		switch b {
		case 0x80:
			r = 0x20AC
		case 0x82:
			r = 0x201A
		case 0x83:
			r = 0x0192
		case 0x84:
			r = 0x201E
		case 0x85:
			r = 0x2026
		case 0x86:
			r = 0x2020
		case 0x87:
			r = 0x2021
		case 0x88:
			r = 0x02C6
		case 0x89:
			r = 0x2030
		case 0x8A:
			r = 0x0160
		case 0x8B:
			r = 0x2039
		case 0x8C:
			r = 0x0152
		case 0x8E:
			r = 0x017D
		case 0x91:
			r = 0x2018
		case 0x92:
			r = 0x2019
		case 0x93:
			r = 0x201C
		case 0x94:
			r = 0x201D
		case 0x95:
			r = 0x2022
		case 0x96:
			r = 0x2013
		case 0x97:
			r = 0x2014
		case 0x98:
			r = 0x02DC
		case 0x99:
			r = 0x2122
		case 0x9A:
			r = 0x0161
		case 0x9B:
			r = 0x203A
		case 0x9C:
			r = 0x0153
		case 0x9E:
			r = 0x017E
		case 0x9F:
			r = 0x0178
		default:
			r = rune(b)
		}

		buf.WriteRune(r)
	}

	return string(buf.Bytes())
}
