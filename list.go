package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type List struct {
	Name             string  `json:"name"`
	Notes            []*Note `json:"notes"`
	h                int     // height, in cells
	topNote, botNote int     // top and bottom-most note indices
}

// Add a new note to end of List.
func (l *List) newNote(text string) {
	note := &Note{text}
	l.Notes = append(l.Notes, note)
	l.UpdateHeight()
	l.boundTopBottomNoteIndices()
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
	l.boundTopBottomNoteIndices()
}

// Swap notes at indices i and j.
func (l *List) swap(i, j int) {
	l.Notes[j], l.Notes[i] = l.Notes[i], l.Notes[j]
}

func (l *List) draw(s tcell.Screen, i int, isListFocused bool, curSelected int) {
	x := listMarginX + (listWidth+listMarginX)*i
	l.drawBox(s, x)
	l.drawNotes(s, x, noteBoxStyle, isListFocused, curSelected)
}

func (l *List) drawBox(s tcell.Screen, x int) {
	h := listMarginY + l.h
	style := focusedListStyle
	if l != kan.currentList() {
		style = unfocusedListStyle
	}
	name := fmt.Sprintf(" %s ", strings.TrimSpace(l.Name))
	drawBox(s, x, listMarginY, x+listWidth, h, style, name, "")
}

func (l *List) drawNotes(s tcell.Screen, x int, style tcell.Style, isListFocused bool, selectedNote int) {
	left, topOfFirstNote := x+noteMargin, listMarginY+noteMargin

	for i, n := range l.notesOnScreen() {
		txt := n.Text
		if isListFocused && i == selectedNote-l.topNote {
			txt = "> " + txt
		}
		topOfCurrentNote := topOfFirstNote + ((noteHeight) * i) + ((noteMargin - 1) * i)
		right := left + listWidth - 2*noteMargin
		bottomOfCurrentNote := topOfCurrentNote + noteHeight - 1
		drawBox(s, left, topOfCurrentNote, right, bottomOfCurrentNote, style, "", txt)
	}
}

// Return only the Notes that are on screen at this moment.
func (l *List) notesOnScreen() []*Note {
	// l.boundTopBottomNoteIndices()

	// // TODO May need to indicate to user if there are notes offscreen.
	if l.length() == 0 {
		return []*Note{}
	}

	return l.Notes[l.topNote : l.botNote+1]
}

func (l *List) boundTopBottomNoteIndices() {
	if l.topNote > l.length()-screenNoteCap {
		l.topNote = l.length() - screenNoteCap
	}

	if l.topNote < 0 {
		l.topNote = 0
	}

	botNote := l.topNote + screenNoteCap - 1
	l.botNote = min(botNote, l.length()-1)
}

// Sets height of List in accordance with the amount of notes it has.
func (l *List) UpdateHeight() {
	a := 0
	if len(l.Notes) == 0 {
		a = noteMargin
	}
	maxNoteCount := min(l.length(), screenNoteCap)
	l.h = ((noteHeight-1)*maxNoteCount + ((maxNoteCount + 1) * noteMargin) + a)
}

// len() wrapper. Used to be more useful, I think.
func (l *List) length() int {
	return len(l.Notes)
}
