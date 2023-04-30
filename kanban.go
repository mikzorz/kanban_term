package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
)

type Kanban struct {
	lists      []*List
	curListIdx int
	curNoteIdx int
}

func (k *Kanban) init() {
	k.lists = make([]*List, 0)
	k.newList("List 1")
	k.newList("List 2")
	k.newList("List 3")
	k.SetListDimensions()
	k.curListIdx = 0
	k.curNoteIdx = 0
}

func (k *Kanban) newList(name string) {
	k.lists = append(k.lists, &List{Name: name, Notes: make([]*Note, 0)})
}

func (k *Kanban) currentList() *List {
	return k.lists[k.curListIdx]
}

func (k *Kanban) currentNote() *Note {
	return k.currentList().Notes[k.curNoteIdx]
}
func (k *Kanban) newNote(text string) {
	k.currentList().newNote(text)
	k.curNoteIdx = k.currentList().length() - 1
}

func (k *Kanban) editNote(newText string) {
	k.currentList().editNote(k.curNoteIdx, newText)
}

func (k *Kanban) deleteNote() {
	// TODO: confirm prompt
	k.currentList().deleteNote(k.curNoteIdx)
}

func (k *Kanban) SetListDimensions() {
	for i, l := range k.lists {
		l.SetDimensions(i)
	}
}

func (k *Kanban) draw(s tcell.Screen) {
	// TODO Should probably determine list positions here, not in l.draw()
	for i, l := range k.lists {
		l.draw(s, i == k.curListIdx, k.curNoteIdx)
	}
}

// Move selection "up" or "down"
func (k *Kanban) moveSelection(dir string) {
	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	max := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	switch dir {
	case "up":
		if k.curNoteIdx > 0 {
			k.curNoteIdx--
		}
	case "down":
		if k.curNoteIdx < k.currentList().length()-1 {
			k.curNoteIdx++
		}
	case "left":
		if k.curListIdx > 0 {
			k.curListIdx--
			k.curNoteIdx = max(0, min(k.curNoteIdx, k.currentList().length()-1))
		}
	case "right":
		if k.curListIdx < len(k.lists)-1 {
			k.curListIdx++
			k.curNoteIdx = max(0, min(k.curNoteIdx, k.currentList().length()-1))
		}
	default:
		log.Fatalf("method Kanban.moveSelection given invalid input: %v+", dir)
	}
}
