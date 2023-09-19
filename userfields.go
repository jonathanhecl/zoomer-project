package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"
)

func waitToSave() {
	for {
		if time.Since(lastSave) > time.Second*30 {
			saveFileUserFields()
		}
		time.Sleep(time.Second)
	}
}

func saveFileUserFields() {
	if time.Since(lastSave).Seconds() > time.Since(lastChange).Seconds() {
		return
	}

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
	err = encoder.Encode(filesData)
	if err != nil {
		fmt.Println(err)
		return
	}

	lastSave = time.Now()
}
