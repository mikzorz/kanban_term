package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

type List struct {
	Name  string  `json:"name"`
	Notes []*Note `json:"notes"`
	Box
}

func (l *List) SetDimensions(listIndex int) {
	l.x, l.y, l.w, l.h = 2+(l.w+2)*listIndex, 1, 22, 0
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
	note := &Note{text}
	l.Notes = append(l.Notes, note)
	l.UpdateHeight()
}

func (l *List) editNote(i int, newText string) {
	l.Notes[i].Text = newText
}

func (l *List) deleteNote(i int) {
	le := l.length()
	if le == 0 {
		return
	}
	firstPart := l.Notes[:i]

	if i == le-1 {
		l.Notes = firstPart
	} else {
		l.Notes = append(firstPart, l.Notes[i+1:]...)
	}
	l.UpdateHeight()
	kan.moveSelection("up")
}

func (l *List) draw(s tcell.Screen, isListFocused bool, curSelected int) {
	l.drawBox(s)
	l.drawNotes(s, noteBoxStyle, isListFocused, curSelected)
}

func (l *List) drawBox(s tcell.Screen) {
	h := l.y + l.h
	if l.length() == 0 {
		h = l.y + 3
	}
	style := focusedListStyle
	if l != kan.currentList() {
		style = unfocusedListStyle
	}
	drawBox(s, l.x, l.y, l.x+l.w, h, style, "")
	name := fmt.Sprintf(" %s ", l.Name)
	ox := 2
	drawText(s, l.x+ox, l.y, l.x+ox+len(name), l.y, defStyle, name)
}

func (l *List) drawNotes(s tcell.Screen, style tcell.Style, isListFocused bool, selectedNote int) {
	x, y := l.x+1, l.y+1

	for i, n := range l.Notes {
		txt := n.Text
		if isListFocused && i == selectedNote {
			txt = "> " + txt
		}
		curY := y + (4 * i)
		drawBox(s, x, curY, x+l.w-2, curY+3, style, txt)
	}
}
