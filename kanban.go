package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
)

type Kanban struct {
	Lists      []*List `json:"lists"`
	curListIdx int
	curNoteIdx int
}

func (k *Kanban) init() {
	k.Lists = make([]*List, 0)
	k.newList("List 1")
	k.newList("List 2")
	k.newList("List 3")
	k.SetListDimensions()
	k.curListIdx = 0
	k.curNoteIdx = 0
}

func (k *Kanban) newList(name string) {
	k.Lists = append(k.Lists, &List{Name: name, Notes: make([]*Note, 0)})
	k.SetListDimensions()
}

func (k *Kanban) currentList() *List {
	if len(k.Lists) == 0 {
		return &List{}
	}
	return k.Lists[k.curListIdx]
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

func (k *Kanban) renameList(newName string) {
	k.currentList().Name = newName
}

func (k *Kanban) deleteNote() {
	// TODO: confirm prompt
	k.currentList().deleteNote(k.curNoteIdx)
}

func (k *Kanban) deleteList() {
	// TODO: confirm prompt
	le := len(k.Lists)
	if le == 0 {
		return
	}
	i := k.curListIdx
	firstPart := k.Lists[:i]

	if i == le-1 {
		k.Lists = firstPart
	} else {
		k.Lists = append(firstPart, k.Lists[i+1:]...)
	}
	k.SetListDimensions()
	k.boundSelection()
}

// Move note from current list to target list
func (k *Kanban) moveNote(target int) {
	k.Lists[target].Notes = append(k.Lists[target].Notes, k.currentNote())
	k.deleteNote()
	k.SetListDimensions()
	k.curListIdx = target
	k.curNoteIdx = k.currentList().length() - 1
}

func (k *Kanban) SetListDimensions() {
	for i, l := range k.Lists {
		l.SetDimensions(i)
	}
}

func (k *Kanban) draw(s tcell.Screen) {
	// TODO Should probably determine list positions here, not in l.draw()
	for i, l := range k.Lists {
		l.draw(s, i == k.curListIdx, k.curNoteIdx)
	}
}

func (k *Kanban) boundSelection() {
	if k.curListIdx < 0 || len(k.Lists) == 0 {
		k.curListIdx = 0
	} else if k.curListIdx >= len(k.Lists) {
		k.curListIdx = len(k.Lists) - 1
	}

	if k.curNoteIdx < 0 || k.currentList().length() == 0 {
		k.curNoteIdx = 0
	} else if k.curNoteIdx >= k.currentList().length() {
		k.curNoteIdx = k.currentList().length() - 1
	}
}

// Move selection "up" or "down"
func (k *Kanban) moveSelection(dir string, shiftHeld bool) {
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
			target := k.curNoteIdx - 1
			if shiftHeld {
				k.currentList().swap(k.curNoteIdx, target)
			}
			k.curNoteIdx = target
		}
	case "down":
		if k.curNoteIdx < k.currentList().length()-1 {
			target := k.curNoteIdx + 1
			if shiftHeld {
				k.currentList().swap(k.curNoteIdx, target)
			}
			k.curNoteIdx = target
		}
	case "left":
		if k.curListIdx > 0 {
			target := k.curListIdx - 1
			if shiftHeld {
				if k.currentList().length() > 0 {
					k.moveNote(target)
				}
			} else {
				k.curListIdx = target
				k.curNoteIdx = max(0, min(k.curNoteIdx, k.currentList().length()-1))
			}

		}
	case "right":
		if k.curListIdx < len(k.Lists)-1 {
			target := k.curListIdx + 1
			if shiftHeld {
				if k.currentList().length() > 0 {
					k.moveNote(target)
				}
			} else {
				k.curListIdx = target
				k.curNoteIdx = max(0, min(k.curNoteIdx, k.currentList().length()-1))

			}
		}
	default:
		log.Fatalf("method Kanban.moveSelection given invalid input: %v+", dir)
	}
}
