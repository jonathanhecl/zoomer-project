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
	userFieldsPath := path.Join(pathProject, userFieldsFilename)
	if !isValidFile(userFieldsPath) {
		fmt.Println("User fields file not found, starting with empty fields")
		userFields = make([]fieldsData, 0)
		return true
	}

	userFieldsData, err := os.Open(userFieldsPath)
	if err != nil {
		fmt.Printf("Error opening user fields file: %v\n", err)
		userFields = make([]fieldsData, 0)
		return true // Continuar sin campos de usuario
	}
	defer userFieldsData.Close()

	userFields = make([]fieldsData, 0)

	decoder := json.NewDecoder(userFieldsData)
	err = decoder.Decode(&userFields)
	if err != nil {
		fmt.Printf("Error decoding user fields file: %v\n", err)
		userFields = make([]fieldsData, 0)
		return true // Continuar sin campos de usuario
	}

	fmt.Printf("User fields loaded: %d field(s)\n", len(userFields))
	return true
}

func saveFileUserFields() {
	userFieldsMutex.Lock()
	defer userFieldsMutex.Unlock()

	if !lastChange.After(lastSave) {
		return // Sin cambios, no hay nada que guardar
	}

	userFieldsPath := path.Join(pathProject, userFieldsFilename)
	
	// Eliminar archivo existente si existe
	if _, err := os.Stat(userFieldsPath); err == nil {
		if err := os.Remove(userFieldsPath); err != nil {
			fmt.Printf("Warning: error removing old user fields file: %v\n", err)
		}
	}

	f, err := os.Create(userFieldsPath)
	if err != nil {
		fmt.Printf("Error creating user fields file: %v\n", err)
		return
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(userFields)
	if err != nil {
		fmt.Printf("Error encoding user fields: %v\n", err)
		return
	}

	lastSave = time.Now()
	fmt.Printf("User fields saved: %d field(s)\n", len(userFields))
}
