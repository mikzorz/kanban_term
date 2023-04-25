package main

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
)

var defStyle = tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorDefault) // Provide option to change, later.
var boxStyle = tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorRed)
var errBoxStyle = tcell.StyleDefault.Background(tcell.ColorOrangeRed).Foreground(tcell.ColorBlack)
var errTextStyle = tcell.StyleDefault.Background(tcell.ColorOrangeRed).Foreground(tcell.ColorBlack)

type context int

const (
	ctxMain context = iota // Main = main screen with note lists. No input boxes, no notes opened.
	ctxInput
)

var currentCtx = ctxMain

var errMsg = ""

func drawScreen(s tcell.Screen) {
	s.Clear() // Because of the background square, this might not be necessary.
	xmax, ymax := s.Size()
	drawBox(s, 0, 0, xmax-1, ymax-1, boxStyle, "") // Background
	drawListBox(s, boxStyle)
	drawNotes(s, boxStyle)
	errMsg = fmt.Sprintf("DEBUG: selected == %d", selected)
	if errMsg != "" {
		drawBox(s, 1, ymax-5, xmax-2, ymax-2, errBoxStyle, errMsg)
	}
	s.Show()
}

func updateLoop(s tcell.Screen) {
	for {

		ev := s.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return
			case tcell.KeyCtrlL:
				// s.Clear() // With drawScreen() at after this in loop, this might be useless.
				errMsg = ""
			}

			r := ev.Rune()
			switch currentCtx {
			case ctxMain:
				if r == 'q' {
					return
				} else if r == 'u' {
					s.Sync()
				} else if r == 'a' {
					newNote(fmt.Sprintf("Note %d", len(list.notes)+1))
				} else if r == 'd' {
					deleteNote()
				} else if ev.Key() == tcell.KeyDown {
					moveSelection("down")
				} else if ev.Key() == tcell.KeyUp {
					moveSelection("up")
				}
			case ctxInput:
				fmt.Print("")
			default:
				errMsg = "unimplemented context enum"
			}
		}
		drawScreen(s)
	}
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
