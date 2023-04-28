package main

import (
	"log"
	"os"
)

var dir = os.TempDir() + "/kanban_term"

var saveFileName = "kanban.json"
var saveFile *os.File

var list *List
var selected = 0

func main() {
	list = &List{}

	errMsg = defErr()

	var err error

	if err = os.MkdirAll(dir, 0700); err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	saveFile = initSaveFile()
	defer saveFile.Close()

	s := newScreen()

	s.SetStyle(defStyle)

	quit := func() {
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	updateLoop(s)
}
