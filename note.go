package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
)

type Note struct {
	Text string `json:"text"`
}

// TODO add json
type List struct {
	Notes []Note `json:"notes"`
	Box
}

type Box struct {
	x, y, w, h int
}

func (l *List) SetDimensions() {
	l.x, l.y, l.w, l.h = 1, 1, 22, 0
	l.UpdateHeight()
}

func (l *List) UpdateHeight() {
	a := 0
	if len(l.Notes) > 0 {
		a = 1
	}
	l.h = (4*len(l.Notes) + a)
}

func (l *List) length() int {
	return len(l.Notes)
}

func (l *List) newNote(text string) {
	note := Note{text}
	l.Notes = append(l.Notes, note)
	l.UpdateHeight()
}

func (l *List) selected() *Note {
	return &l.Notes[selected]
}

func (l *List) editNote(newText string) {
	l.selected().Text = newText
}

func (l *List) deleteNote() {
	le := l.length()
	if le == 0 {
		return
	}
	firstPart := l.Notes[:selected]

	if selected == le-1 {
		l.Notes = firstPart
	} else {
		l.Notes = append(firstPart, l.Notes[selected+1:]...)
	}
	l.UpdateHeight()
	moveSelection("up")
}

func (l *List) draw(s tcell.Screen) {
	l.drawBox(s, boxStyle)
	l.drawNotes(s, boxStyle)
}

func (l *List) drawBox(s tcell.Screen, style tcell.Style) {
	if l.length() > 0 {
		drawBox(s, list.x, list.y, list.x+list.w, list.y+list.h, style, "")
		name := " List "
		ox := 2
		drawText(s, list.x+ox, list.y, list.x+ox+len(name), list.y, defStyle, name)
	}
}

func (l *List) drawNotes(s tcell.Screen, style tcell.Style) {
	x, y := list.x+1, list.y+1

	for i := 0; i < list.length(); i++ {
		n := list.Notes[i]
		txt := n.Text
		if i == selected {
			txt = "> " + txt
		}
		curY := y + (4 * i)
		drawBox(s, x, curY, x+list.w-2, curY+3, style, txt)
	}
}

// Move selection "up" or "down"
func moveSelection(dir string) {
	if list.length() == 0 {
		return
	}

	switch dir {
	case "up":
		if selected > 0 {
			selected--
		}
	case "down":
		if selected < list.length()-1 {
			selected++
		}
	default:
		log.Fatalf("func moveSelection given invalid input: %v+", dir)
	}
}
