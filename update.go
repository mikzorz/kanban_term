package main

import (
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
			// case tcell.KeyEscape, tcell.KeyCtrlC:
			case tcell.KeyCtrlC:
				return
			case tcell.KeyCtrlL:
				// s.Clear() // With drawScreen() after this in loop, this might be useless.
				errMsg = ""
			}

			switch currentCtx {
			case ctxMain:
				ctxMainHandler(s, ev)
			case ctxNoteView:
				ctxNoteViewHandler(s, ev)
			case ctxConfirm:
				ctxConfirmHandler(s, ev)
				if ctxConfirmHandler(s, ev) {
					return
				}
			default:
				errMsg = "unimplemented context enum"
			}

			drawScreen(s)
			s.Show()
		}
	}
}

func ctxMainHandler(s tcell.Screen, ev *tcell.EventKey) {
	r := ev.Rune()
	switch r {
	case 'q':
		currentCtx = ctxConfirm
	case 'u':
		s.Sync()
	case 'a':
		addNote(s)
	case 'e':
		editNote(s)
	case 'd':
		list.deleteNote()
	case 's':
		saveToFile()
	case 'v':
		currentCtx = ctxNoteView
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
}

func ctxNoteViewHandler(s tcell.Screen, ev *tcell.EventKey) {
	r := ev.Rune()
	switch r {
	case 'q', 'v':
		currentCtx = ctxMain
	case 'u':
		s.Sync()
	case 'e':
		editNote(s)
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
}

func ctxConfirmHandler(s tcell.Screen, ev *tcell.EventKey) (quit bool) {
	r := ev.Rune()
	switch r {
	case 'y', 'Y':
		return true
	default:
		currentCtx = ctxMain
		return false
	}
}

func addNote(s tcell.Screen) {
	openEditorStart(s, "", list.newNote)
}

func editNote(s tcell.Screen) {
	openEditorStart(s, list.selected().Text, list.editNote)
}
