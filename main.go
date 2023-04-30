package main

import (
	"log"
	"os"
)

const DEBUG_MODE = true

var dir = os.TempDir() + "/kanban_term"

var saveFileName = "kanban.json"
var saveFile *os.File

// var list *List
var kan Kanban
var curList *List

type Box struct {
	x, y, w, h int
}

func main() {
	kan.init()

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
