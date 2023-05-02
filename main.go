package main

import (
	"log"
	"os"
)

const DEBUG_MODE = false

// const DEBUG_MODE = true

var dir = os.TempDir() + "/kanban_term"

const saveFileName = "kanban_term.json" // TODO make sure that the save file is in either the PWD, or is supplied as an arg.

var saveFile *os.File

var kan Kanban
var curList *List

const noteHeight = 4
const noteMargin = 1 // margin is actually an offset. The gap between 2 lines == margin - 1 (the line being offset).

const listWidth = 22
const listMarginX = 2
const listMarginY = 1

func main() {

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
