package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

var tryingToQuit = false

var onConfirm = func() {}

func updateLoop(s tcell.Screen) {
	xmax, _ := s.Size()
	screenListCap = maxListsOnScreen(xmax)
	for {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
			xmax, _ = s.Size()
			screenListCap = maxListsOnScreen(xmax)
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
				if ctxConfirmHandler(s, ev) {
					if tryingToQuit {
						return
					}
					onConfirm()
					currentCtx = ctxMain
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
		tryingToQuit = true
		attemptedAction = ActionQuit
		currentCtx = ctxConfirm
	case 'u':
		s.Sync()
	case 'a':
		addNote(s)
	case 'A':
		kan.newList(fmt.Sprintf("List %d", len(kan.Lists)+1))
	case 'e':
		editNote(s)
	case 'r':
		renameList(s)
	case 'd':
		if kan.isNoteDeletable() {
			setConfirm(kan.deleteNote, ActionDeleteNote)
		}
	case 'D':
		if kan.isListDeletable() {
			setConfirm(kan.deleteList, ActionDeleteList)
		}
	case 's':
		saveToFile()
	case 'v':
		currentCtx = ctxNoteView
	case 'o':
		keyBindingsStringIndex = (keyBindingsStringIndex + 1) % 3
	default:
		errMsg = defErr()
		handleSelectionMovement(ev)
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
		if kan.isNoteDeletable() {
			setConfirm(kan.deleteNote, ActionDeleteNote)
		}
	case 's':
		saveToFile()
	default:
		handleSelectionMovement(ev)

		errMsg = defErr()

	}
}

func ctxConfirmHandler(s tcell.Screen, ev *tcell.EventKey) (quit bool) {
	r := ev.Rune()
	switch r {
	case 'y', 'Y':
		return true
	default:
		tryingToQuit = false
		currentCtx = ctxMain
		return false
	}
}

// fn is the function that will run after the confirmation prompt is given a y/Y by the user.
// a is the action enum that will be substituted into the confirmation prompt. e.g. "Are you sure you want to {a} [y/N]"
func setConfirm(fn func(), a action) {
	onConfirm = fn
	attemptedAction = a
	currentCtx = ctxConfirm
}

func addNote(s tcell.Screen) {
	openEditorStart(s, "", kan.newNote)
}

func editNote(s tcell.Screen) {
	if kan.currentList().length() > 0 {
		openEditorStart(s, kan.currentNote().Text, kan.editNote)
	}
}

func save(s tcell.Screen) {
	// TODO: maybe prompt for confirmation, but also show a small window for a few seconds that confirms whether or not the file was saved successfully.
}

func renameList(s tcell.Screen) {
	openEditorStart(s, kan.currentList().Name, kan.renameList)
}

func handleSelectionMovement(ev *tcell.EventKey) {
	mod := ev.Modifiers()
	shiftHeld := mod == tcell.ModShift
	ctrlHeld := mod == tcell.ModCtrl

	switch ev.Key() {
	case tcell.KeyDown:
		kan.moveSelection("down", shiftHeld, ctrlHeld)
	case tcell.KeyUp:
		kan.moveSelection("up", shiftHeld, ctrlHeld)
	case tcell.KeyLeft:
		kan.moveSelection("left", shiftHeld, ctrlHeld)
	case tcell.KeyRight:
		kan.moveSelection("right", shiftHeld, ctrlHeld)
	default:
		errMsg = "that key does nothing"
	}
	// errMsg = fmt.Sprintf("EventKey Modifiers: %d, noteIndex: %d, listIndex: %d", mod, kan.curNoteIdx, kan.curListIdx)
	errMsg += fmt.Sprintf("l-list = %d, r-list = %d", kan.l_list, kan.r_list)

}
