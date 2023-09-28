package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
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
			fmt.Println("Filename:", f.Filename+" Method:", f.Content[method])
			//content += `<div class="fields">`
			//for _, mtd := range filesData[f.Filename].getMethods() {
			//	content += `<div class="method">` + mtd + `</div><br>`
			//	for _, field := range configProject.UserFields {
			//		content += `<div class="field">`
			//		if field.Type == EnumBoolean {
			//			content += `<label><input type="checkbox" name="` + createFieldName(f.Filename, mtd, field.Name) + `" value="` + field.Name + `" `
			//			if getUserValue(f.Filename, mtd, field.Name) == "1" {
			//				content += `checked`
			//			}
			//			content += ` onchange="saveChange(this)"> ` + field.Name + `</label>`
			//		} else if field.Type == EnumTextBox {
			//			content += `<label>` + field.Name + `<br/><textarea name="` + createFieldName(f.Filename, mtd, field.Name) + `" onchange="saveChange(this)">`
			//			content += getUserValue(f.Filename, mtd, field.Name)
			//			content += `</textarea></label>`
			//		}
			//		content += `</div>`
			//	}
			//}
			//content += `</div>`
		}
		prevMethod = method
	}

	return content
}

func (f fileData) getContent() string {
	return strings.Join(f.Content, "\n")
}

func (f fileData) getMethods() []string {
	methods := []string{}
	for _, method := range f.Methods {
		mtd := f.Content[method]
		mtd = strings.ReplaceAll(mtd, "\n", "")
		mtd = strings.ReplaceAll(mtd, "\r", "")
		methods = append(methods, mtd)
	}
	return methods
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
		fmt.Println(err)
		return false
	}

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
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	content := []string{}
	methods := []int{}

	for i, line := range strings.Split(string(data), "\n") {
		content = append(content, line)
		for _, method := range configProject.MethodFilter {
			re, _ := regexp.Compile(method)
			if re.MatchString(line) {
				methods = append(methods, i)
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
			filesTmp, err = scanProject(path.Join(root, f.Name())+string(filepath.Separator), filesOut)
			if err != nil {
				return nil, err
			}
			filesOut = filesTmp
		}
	}

	return filesOut, nil
}

func changedUserField(name string, value string) bool {
	fmt.Println("Changed:", name, value)

	filename, method, field := disassemblyFieldName(name)

	setUserValue(filename, method, field, value)
	lastChange = time.Now()

	return true
}
