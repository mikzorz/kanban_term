package main

import (
	"log"
	"os"
)

var dir = os.TempDir() + "/kanban_term" //

func main() {
	var err error
	// dir, err = ioutil.TempDir("", "")
	err = os.MkdirAll(dir, 0700)
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

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

	newNote("Note 1")
	newNote("Note 2")
	newNote("Note 3 Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.")
	selected = 0

	updateLoop(s)
}
