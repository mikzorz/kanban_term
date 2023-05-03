package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Opens file if it exists, creates new if it doesn't. Tries to read saved kanban.
func initSaveFile() *os.File {
	f, err := os.OpenFile(saveFileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		f, err = nil, fmt.Errorf("unable to load or create a save file")
	} else {
		err = loadFromSave(f)
		if err != nil {
			f, err = loadFromBackup()
		} else {
			err = backupSave()
		}
		kan.UpdateAllLists()
	}
	if err != nil {
		showInfoBox(time.Second*5, err.Error())
	}
	return f
}

func newSaveFile(filename string) (*os.File, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// Overwrite kanban.json with current in-mem kanban.
func saveToFile() error {
	j, err := json.Marshal(&kan)
	if err != nil {
		return err
	}

	if err = saveFile.Truncate(0); err != nil {
		// log.Fatalf("%+v", err)
	}
	if _, err = saveFile.Seek(0, 0); err != nil {
		// log.Fatalf("%+v", err)
	}

	if _, err = saveFile.Write(j); err != nil {
		return err
	}
	err = backupSave()
	if err != nil {
		return err
	}
	return nil
}

func loadFromSave(f *os.File) error {
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	var data []byte = make([]byte, fi.Size())
	if _, err = f.Read(data); err != nil {
		return err
	}

	if err = json.Unmarshal(data, &kan); err != nil {
	}
	return err
}

func backupSave() error {
	data, err := os.ReadFile(saveFileName)
	if err != nil {
		return err
	}

	backupFile, err := os.OpenFile(backupFileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
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
		return err
	}
	return nil
}

func loadFromBackup() (*os.File, error) {
	// TODO I'm returning the main save file, but setting the global backup file directly. Should probably return both or set both.

	f, err := newSaveFile(saveFileName)
	if err != nil {
		return f, err
	}
	bck, err := os.OpenFile(backupFileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("unable to load or create backup save file")
	} else {
		err = loadFromSave(bck)
		if err != nil {
			backupFile, err = newSaveFile(backupFileName)
			kan.newKanban()
		}
	}
	return f, err
}
