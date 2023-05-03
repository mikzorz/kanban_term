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
		loadFromBackup()
	} else {
		err = loadFromSave(f)
		if err != nil {
			loadFromBackup()
		} else {
			backupSave()
		}
		kan.UpdateAllLists()
	}

	return f
}

func newSaveFile(filename string) *os.File {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("can't create new save file, %+v", err)
	}

	return f
}

// Overwrite kanban.json with current in-mem kanban.
func saveToFile() {
	j, err := json.Marshal(&kan)
	if err != nil {
		errMsg = fmt.Sprintf("Error: failed to marshal JSON")
	}

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
	backupSave()
}

func loadFromSave(f *os.File) error {
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	var data []byte = make([]byte, fi.Size())
	if _, err = f.Read(data); err != nil {
		f = newSaveFile(saveFileName)
		return err
	}

	if err = json.Unmarshal(data, &kan); err != nil {
		f = newSaveFile(saveFileName)
	}
	return err
}

func backupSave() {
	data, err := os.ReadFile(saveFileName)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	backupFile, err := os.OpenFile(backupFileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer backupFile.Close()

	if err = backupFile.Truncate(0); err != nil {
		// log.Fatalf("%+v", err)
	}
	if _, err = backupFile.Seek(0, 0); err != nil {
		// log.Fatalf("%+v", err)
	}

	_, err = backupFile.Write(data)
	// err = os.WriteFile(backupFileName, data, 0755)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func loadFromBackup() *os.File {
	// TODO I'm returning the main save file, but setting the global backup file directly. Should probably return both or set both.
	f := newSaveFile(saveFileName)
	bck, err := os.OpenFile(backupFileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		backupFile = newSaveFile(backupFileName)
		kan.newKanban()
	} else {
		err = loadFromSave(bck)
		if err != nil {
			backupFile = newSaveFile(backupFileName)
			kan.newKanban()
		}
	}
	return f
}
