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
		case *tcell.EventKey:
			loopCount++
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return
			case tcell.KeyCtrlL:
				// s.Clear() // With drawScreen() after this in loop, this might be useless.
				errMsg = ""
			}

			r := ev.Rune()
			switch currentCtx {
			case ctxMain:
				switch r {
				case 'q':
					return
				case 'u':
					s.Sync()
				case 'a':
					newNote(fmt.Sprintf("Note %d", len(list.Notes)+1))
				case 'e':
					// Suspend and Resume are needed to stop text editor from bugging out. Took me too long to figure this out.
					err := s.Suspend()
					if err != nil {
						log.Fatalf("%+v", err)
					}

					newText := openTextPrompt(list.Notes[selected].Text)
					editNote(&list.Notes[selected], newText)

					err = s.Resume()
					if err != nil {
						log.Fatalf("%+v", err)
					}
				case 'd':
					deleteNote()
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
			case ctxNoteView:
				// TODO
			default:
				errMsg = "unimplemented context enum"
			}

			drawScreen(s)
		}
		s.Show()

	}
}
