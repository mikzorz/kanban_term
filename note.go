package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
)

type Note struct {
	Text string `json:"text"`
}

var selected = -1

// TODO add json
type List struct {
	Notes []Note `json:"notes"`
	Box
}

type Box struct {
	x, y, w, h int
}

var list *List = &List{}

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

func newNote(text string) {
	note := Note{text}
	list.Notes = append(list.Notes, note)
	list.UpdateHeight()
}

func editNote(n *Note, newText string) {
	n.Text = newText
}

func deleteNote() {
	l := len(list.Notes)
	if l == 0 {
		return
	}
	firstHalf := list.Notes[:selected]

	if selected == l-1 {
		list.Notes = firstHalf
	} else {
		list.Notes = append(firstHalf, list.Notes[selected+1:]...)
	}
	list.UpdateHeight()
	moveSelection("up")
}

func drawListBox(s tcell.Screen, style tcell.Style) {
	if len(list.Notes) > 0 {
		drawBox(s, list.x, list.y, list.x+list.w, list.y+list.h, style, "")
		name := " List "
		ox := 2
		drawText(s, list.x+ox, list.y, list.x+ox+len(name), list.y, defStyle, name)
	}
}

func drawNotes(s tcell.Screen, style tcell.Style) {
	x, y := list.x+1, list.y+1

	for i := 0; i < len(list.Notes); i++ {
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
	if len(list.Notes) == 0 {
		return
	}

	switch dir {
	case "up":
		if selected > 0 {
			selected--
		}
	case "down":
		if selected < len(list.Notes)-1 {
			selected++
		}
	default:
		log.Fatalf("func moveSelection given invalid input: %v+", dir)
	}
}
