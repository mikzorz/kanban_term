package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Opens file if it exists, creates new if it doesn't. Tries to read saved kanban.
func initSaveFile() *os.File {
	f, err := os.OpenFile(saveFileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		f = newSaveFile()
	} else {
		fi, err := f.Stat()
		if err != nil {
			log.Fatalf("%+v", err)
		}

		var data []byte = make([]byte, fi.Size())
		if _, err = f.Read(data); err != nil {
			log.Fatalf("file exists but can't read data, %+v", err)
		}

		if err = json.Unmarshal(data, &list); err != nil {
			errMsg = fmt.Sprintf("Error: file can be read but can't parse json, %s", err.Error())
			f = newSaveFile()
		}
		list.SetDimensions()
	}

	return f
}

func newSaveFile() *os.File {
	f, err := os.Create(saveFileName)
	if err != nil {
		log.Fatalf("can't create new save file, %+v", err)
	}

	newNote("Example Note 1")
	newNote("Example Note 2")

	return f
}

func saveToFile() {
	j, err := json.Marshal(&list)
	if err != nil {
		errMsg = fmt.Sprintf("Error: failed to marshal JSON")
	}

	// Erase old contents. TODO: should backup first. Or just use a db like sqlite.
	if err = saveFile.Truncate(0); err != nil {
		// log.Fatalf("%+v", err)
	}
	if _, err = saveFile.Seek(0, 0); err != nil {
		// log.Fatalf("%+v", err)
	}

	if _, err = saveFile.Write(j); err != nil {
		errMsg = fmt.Sprintf("Error: failed to write JSON to file")
		log.Fatalf("%+v", err)
	}
}

func backupSave() {
	// create copy of saveFile, restore when necessary
}
