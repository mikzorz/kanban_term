package main

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
)

func updateLoop(s tcell.Screen) {
	for {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
			drawScreen(s)
			s.Show()
		case *tcell.EventKey:

			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return
			case tcell.KeyCtrlL:
				// s.Clear() // With drawScreen() after this in loop, this might be useless.
				errMsg = ""
			}

			switch currentCtx {
			case ctxMain:
				if ctxMainHandler(s, ev) {
					return
				}
			case ctxNoteView:
				// TODO
			default:
				errMsg = "unimplemented context enum"
			}

			drawScreen(s)
			s.Show()
		}
	}
}

func ctxMainHandler(s tcell.Screen, ev *tcell.EventKey) (quit bool) {
	r := ev.Rune()
	switch r {
	case 'q':
		return true
	case 'u':
		s.Sync()
	case 'a':
		list.newNote(fmt.Sprintf("Note %d", list.length()+1))
	case 'e':
		// Suspend and Resume are needed to stop text editor from bugging out. Took me too long to figure this out.
		err := s.Suspend()
		if err != nil {
			log.Fatalf("%+v", err)
		}

		newText := openTextPrompt(list.selected().Text)
		list.editNote(newText)

		err = s.Resume()
		if err != nil {
			log.Fatalf("%+v", err)
		}
	case 'd':
		list.deleteNote()
	case 's':
		saveToFile()
	default:
		if ev.Key() == tcell.KeyDown {
			moveSelection("down")
		} else if ev.Key() == tcell.KeyUp {
			moveSelection("up")
		} else {
			errMsg = "that key does nothing"
		}

		errMsg = defErr()
	}
	return false
}
