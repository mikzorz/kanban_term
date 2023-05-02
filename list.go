package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type List struct {
	Name  string  `json:"name"`
	Notes []*Note `json:"notes"`
	Box
}

// Sets size and position of List based on its index within the kanban and the amount of notes the List has.
func (l *List) SetDimensions(listIndex int) {
	l.w = listWidth
	l.x, l.y, l.h = listMargin+(l.w+listMargin)*listIndex, 1, 0
	l.UpdateHeight()
}

// Sets height of List in accordance with the amount of notes it has.
func (l *List) UpdateHeight() {
	a := 0
	if len(l.Notes) == 0 {
		a = noteMargin
	}
	l.h = ((noteHeight-1)*len(l.Notes) + ((len(l.Notes) + 1) * noteMargin) + a)
}

// len() wrapper. Used to be more useful, I think.
func (l *List) length() int {
	return len(l.Notes)
}

// Add a new note to end of List.
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
	kan.boundSelection()
}

// Swap notes at indices i and j.
func (l *List) swap(i, j int) {
	l.Notes[j], l.Notes[i] = l.Notes[i], l.Notes[j]
}

func (l *List) draw(s tcell.Screen, isListFocused bool, curSelected int) {
	l.drawBox(s)
	l.drawNotes(s, noteBoxStyle, isListFocused, curSelected)
}

func (l *List) drawBox(s tcell.Screen) {
	h := l.y + l.h
	style := focusedListStyle
	if l != kan.currentList() {
		style = unfocusedListStyle
	}
	name := fmt.Sprintf(" %s ", strings.TrimSpace(l.Name))
	drawBox(s, l.x, l.y, l.x+l.w, h, style, name, "")
}

func (l *List) drawNotes(s tcell.Screen, style tcell.Style, isListFocused bool, selectedNote int) {
	left, topOfFirstNote := l.x+noteMargin, l.y+noteMargin

	for i, n := range l.Notes {
		txt := n.Text
		if isListFocused && i == selectedNote {
			txt = "> " + txt
		}
		topOfCurrentNote := topOfFirstNote + ((noteHeight) * i) + ((noteMargin - 1) * i)
		right := left + l.w - 2*noteMargin
		bottomOfCurrentNote := topOfCurrentNote + noteHeight - 1
		drawBox(s, left, topOfCurrentNote, right, bottomOfCurrentNote, style, "", txt)
	}
}
