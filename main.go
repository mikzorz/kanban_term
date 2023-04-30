package main

import (
	"log"
	"os"
)

// TODO: implement scrolling. Too many lists, lists go off screen. Too many notes, notes go offscreen.
// Selection cursor will "push" the camera if there is more to see offscreen. Camera is not tightly locked to cursor.

const DEBUG_MODE = true

var dir = os.TempDir() + "/kanban_term"

var saveFileName = "kanban.json"
var saveFile *os.File

// var list *List
var kan Kanban
var curList *List

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
