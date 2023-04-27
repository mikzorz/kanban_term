package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/gdamore/tcell/v2"
)

var defStyle = tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorDefault) // Provide option to change, later.
var boxStyle = tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorRed)
var errBoxStyle = tcell.StyleDefault.Background(tcell.ColorOrangeRed).Foreground(tcell.ColorBlack)
var errTextStyle = tcell.StyleDefault.Background(tcell.ColorOrangeRed).Foreground(tcell.ColorBlack)

type context int

const (
	ctxMain     context = iota // Main = main screen with note lists. No input boxes, no notes opened.
	ctxNoteView                // NoteView = the screen, in tcell, for viewing full note contents. Not the text editor.
)

var currentCtx = ctxMain

var errMsg = ""

// For debugging purposes only.
var loopCount = 0

var inputBoxW, inputBoxH = 60, 8

func drawScreen(s tcell.Screen) {
	s.Clear() // Because of the background square, this might not be necessary.
	xmax, ymax := s.Size()
	drawBox(s, 0, 0, xmax-1, ymax-1, boxStyle, "") // Background
	drawListBox(s, boxStyle)
	drawNotes(s, boxStyle)

	if currentCtx == ctxNoteView {
		left := (xmax-1)/2 - (inputBoxW / 2)
		right := (xmax-1)/2 + (inputBoxW / 2)
		top := (ymax-1)/2 - (inputBoxH / 2)
		bottom := (ymax-1)/2 + (inputBoxH / 2)
		drawBox(s, left, top, right, bottom, boxStyle, "")
		promptMsg := " Note "
		drawText(s, left+2, top, left+2+len(promptMsg), top, defStyle, promptMsg)
		drawText(s, left+2, top+2, right-2, bottom-2, defStyle, list.notes[selected].text)
		// TODO separate text upto cursor, then invert text under cursor, then continue with normal text.
	} else {
		drawText(s, 5, ymax-1, xmax-1, ymax-1, defStyle, " q: Quit, a: Add, e: Edit, d: Delete, up/down arrows: Change selection, u: refresh ")
	}

	if errMsg != "" {
		drawBox(s, 1, ymax-5, xmax-2, ymax-2, errBoxStyle, errMsg)
	}
}

func defErr() string {
	return fmt.Sprintf("DEBUG: len(list.notes)=%d selected=%d loopCount=%d", len(list.notes), selected, loopCount)
}

func updateLoop(s tcell.Screen) {
	errMsg = defErr()
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
					newNote(fmt.Sprintf("Note %d", len(list.notes)+1))
				case 'e':
					// Suspend and Resume are needed to stop text editor from bugging out. Took me too long to figure this out.
					err := s.Suspend()
					if err != nil {
						log.Fatalf("%+v", err)
					}

					newText := openTextPrompt(list.notes[selected].text)
					editNote(&list.notes[selected], newText)

					err = s.Resume()
					if err != nil {
						log.Fatalf("%+v", err)
					}
				case 'd':
					deleteNote()
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

func openTextPrompt(s string) string {
	// I don't know how to hook directly into a text editor, I tried, didn't work, so I won't. Using tempfiles instead.

	file, err := ioutil.TempFile(dir, "note*.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write([]byte(s))
	if err != nil {
		log.Fatalf("%+v", err)
	}
	file.Close() // Closing before editing manually seems like a good idea, right?

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	log.Printf("Using %s", editor)

	cmd := exec.Command(editor, file.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	log.Printf("Running %s...", editor)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	b, err := os.ReadFile(file.Name())
	if err != nil {
		log.Fatalf("%+v", err)
	}

	return string(b)
}

func newScreen() tcell.Screen {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	return s
}

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}
func drawErrBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	drawBox(s, x1, y1, x2, y2, style, text)
}

func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			s.SetContent(col, row, ' ', nil, style)
		}
	}

	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, style)
		s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1; row <= y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}

	if y1 != y2 && x1 != x2 {
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}

	textStyle := defStyle
	if style == errBoxStyle {
		textStyle = errTextStyle
	}
	drawText(s, x1+1, y1+1, x2-1, y2-1, textStyle, text)
}
