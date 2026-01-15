package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sync"
	"time"
)

type fieldsData struct {
	Filename string
	Method   string
	Field    string
	Value    string
}

var (
	userFieldsMutex sync.Mutex
)

func waitToSave() {
	c := time.Tick(time.Second * 30)
	for range c {
		if time.Since(lastSave) > time.Second*30 {
			saveFileUserFields()
		}
	}
}

func setUserValue(filename string, method string, field string, value string) {
	userFieldsMutex.Lock()
	defer userFieldsMutex.Unlock()

	found := false

	for i := range userFields {
		if userFields[i].Filename == filename &&
			userFields[i].Method == method &&
			userFields[i].Field == field {
			userFields[i].Value = value
			found = true
			break
		}
	}

	if !found {
		userFields = append(userFields, fieldsData{
			Filename: filename,
			Method:   method,
			Field:    field,
			Value:    value,
		})
	}
}

func getUserValue(filename string, method string, field string) string {
	userFieldsMutex.Lock()
	defer userFieldsMutex.Unlock()

	for _, userField := range userFields {
		if userField.Filename == filename &&
			userField.Method == method &&
			userField.Field == field {
			return userField.Value
		}
	}
	return ""
}

func loadUserFields() bool {
	if !isValidFile(path.Join(pathProject, userFieldsFilename)) {
		return false
	}

	userFieldsData, err := os.Open(path.Join(pathProject, userFieldsFilename))
	if err != nil {
		fmt.Println(err)
		return false
	}

	userFields = make([]fieldsData, 0)

	decoder := json.NewDecoder(userFieldsData)
	err = decoder.Decode(&userFields)
	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println("User fields loaded! (", len(userFields), ")")

	return true
}

func saveFileUserFields() {
	userFieldsMutex.Lock()
	defer userFieldsMutex.Unlock()

	if time.Since(lastSave) <= time.Since(lastChange) {
		fmt.Println("No changes to save")
		return
	}
	fmt.Println("Saving changes")

	if _, err := os.Stat(path.Join(pathProject, userFieldsFilename)); err == nil {
		os.Remove(path.Join(pathProject, userFieldsFilename))
	}

	f, err := os.Create(path.Join(pathProject, userFieldsFilename))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	err = encoder.Encode(userFields)
	if err != nil {
		fmt.Println(err)
		return
	}

	lastSave = time.Now()
}
