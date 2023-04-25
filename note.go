package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
)

type Note struct {
	text string
}

var selected = -1

type List struct {
	notes      []Note
	x, y, w, h int
}

var list *List = &List{x: 1, y: 1, w: 22, h: 0}

func (l *List) UpdateHeight() {
	a := 0
	if len(l.notes) > 0 {
		a = 1
	}
	l.h = (4*len(l.notes) + a)
}

func newNote(text string) {
	note := Note{text}
	list.notes = append(list.notes, note)
	list.UpdateHeight()
}

func deleteNote() {
	l := len(list.notes)
	if l == 0 {
		return
	}
	firstHalf := list.notes[:selected]

	if selected == l-1 {
		list.notes = firstHalf
	} else {
		list.notes = append(firstHalf, list.notes[selected+1:]...)
	}
	list.UpdateHeight()
	moveSelection("up")
}

func drawListBox(s tcell.Screen, style tcell.Style) {
	if list.h > 0 {
		drawBox(s, list.x, list.y, list.x+list.w, list.y+list.h, style, "")
	}
}

func drawNotes(s tcell.Screen, style tcell.Style) {
	x, y := list.x+1, list.y+1

	for i := 0; i < len(list.notes); i++ {
		n := list.notes[i]
		txt := n.text
		if i == selected {
			txt = "> " + txt
		}
		curY := y + (4 * i)
		drawBox(s, x, curY, x+list.w-2, curY+3, style, txt)
	}
}

// Move selection "up" or "down"
func moveSelection(dir string) {
	if len(list.notes) == 0 {
		return
	}

	switch dir {
	case "up":
		if selected > 0 {
			selected--
		}
	case "down":
		if selected < len(list.notes)-1 {
			selected++
		}
	default:
		log.Fatalf("func moveSelection given invalid input: %v+", dir)
	}
}
