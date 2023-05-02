package main

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
)

type Kanban struct {
	Lists          []*List `json:"lists"`
	curListIdx     int
	curNoteIdx     int
	l_list, r_list int // left and rightmost list indices
}

// If starting from a blank save, create three lists.
func (k *Kanban) newKanban() {
	k.Lists = make([]*List, 0)
	k.newList("To do")
	k.newList("Doing")
	k.newList("Done!")
	for i := 0; i < 5; i++ {
		k.newList(fmt.Sprintf("List %d", i))
	}
	k.curListIdx = 0
	kan.newNote("Example Note 1")
	kan.newNote("Example Note 2")
	k.curNoteIdx = 0
	k.l_list = 0
	k.r_list = 0
}

func (k *Kanban) newList(name string) {
	l := &List{Name: name, Notes: make([]*Note, 0)}
	l.UpdateHeight()
	k.Lists = append(k.Lists, l)
	k.l_list = len(k.Lists) - screenListCap
	k.curListIdx = len(k.Lists) - 1
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
	if text == "" {
		return
	}
	k.currentList().newNote(text)
	k.curNoteIdx = k.currentList().length() - 1
}

func (k *Kanban) editNote(newText string) {
	k.currentList().editNote(k.curNoteIdx, newText)
}

func (k *Kanban) renameList(newName string) {
	k.currentList().Name = newName
}

func (k *Kanban) isNoteDeletable() bool {
	if k.currentList().length() == 0 {
		return false
	}
	return true
}

func (k *Kanban) deleteNote() {
	k.currentList().deleteNote(k.curNoteIdx)
}

func (k *Kanban) isListDeletable() bool {
	if len(k.Lists) == 0 {
		return false
	}
	return true
}

func (k *Kanban) deleteList() {
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
	k.boundSelection()
	k.boundLeftRightListIndices()
}

// Move note from current list to target list
func (k *Kanban) moveNote(target int) {
	k.Lists[target].Notes = append(k.Lists[target].Notes, k.currentNote())
	k.deleteNote()
	k.curListIdx = target
	k.currentList().UpdateHeight()
	k.curNoteIdx = k.currentList().length() - 1
}

func (k *Kanban) UpdateAllListHeights() {
	for _, l := range k.Lists {
		l.UpdateHeight()
	}
}

// Swap positions of two lists.
func (k *Kanban) swap(i, j int) {
	k.Lists[i], k.Lists[j] = k.Lists[j], k.Lists[i]
}

func (k *Kanban) draw(s tcell.Screen) {
	for i, l := range k.listsOnScreen() {
		l.draw(s, i, i == k.curListIdx, k.curNoteIdx)
	}
}

// Restrict cursor position to the bounds of the kanban (horizontally) and current list (vertically).
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

// Move cursor "up" & "down" through a list. Move "left" & "right" between lists.
// Hold "Shift" to move note.
// Hold "Control" to move list.
func (k *Kanban) moveSelection(dir string, shiftHeld, ctrlHeld bool) {

	switch dir {
	case "up":
		if k.curNoteIdx > 0 {
			k.moveVertical(k.curNoteIdx-1, shiftHeld)
		}
	case "down":
		if k.curNoteIdx < k.currentList().length()-1 {
			k.moveVertical(k.curNoteIdx+1, shiftHeld)
		}
	case "left":
		if k.curListIdx > 0 {
			if k.curListIdx-1 < k.l_list {
				k.l_list--
				k.r_list--
			}
			k.moveHorizontal(k.curListIdx-1, shiftHeld, ctrlHeld)
		}
	case "right":
		if k.curListIdx < len(k.Lists)-1 {
			if k.curListIdx+1 > k.r_list {
				k.l_list++
				k.r_list++
			}
			k.moveHorizontal(k.curListIdx+1, shiftHeld, ctrlHeld)
		}
	default:
		log.Fatalf("method Kanban.moveSelection given invalid input: %v+", dir)
	}
}

// TODO implement scrolling vertical through long lists
func (k *Kanban) moveVertical(targetIndex int, shiftHeld bool) {
	if shiftHeld {
		k.currentList().swap(k.curNoteIdx, targetIndex)
	}
	k.curNoteIdx = targetIndex
}

func (k *Kanban) moveHorizontal(targetIndex int, shiftHeld, ctrlHeld bool) {
	if ctrlHeld {
		k.swap(k.curListIdx, targetIndex)
		k.curListIdx = targetIndex
	} else if shiftHeld {
		if k.currentList().length() > 0 {
			k.moveNote(targetIndex)
		}
	} else {
		k.curListIdx = targetIndex
		k.curNoteIdx = max(0, min(k.curNoteIdx, k.currentList().length()-1))
	}
}

// Return only the Lists that are on screen on this moment.
func (k *Kanban) listsOnScreen() []*List {
	k.boundLeftRightListIndices()

	// TODO May need to indicate to user if there are lists offscreen.
	if len(k.Lists) == 0 { //
		return []*List{}
	}
	return k.Lists[k.l_list : k.r_list+1]
}

func (k *Kanban) boundLeftRightListIndices() {
	if k.l_list >= len(k.Lists) {
		k.l_list = len(k.Lists) - 1
	}

	rlist := k.l_list + screenListCap - 1
	k.r_list = min(rlist, len(k.Lists)-1)

}
