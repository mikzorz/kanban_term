package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/gdamore/tcell/v2"
)

var defStyle = tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorDefault) // Provide option to change, later.
var boxStyle = tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorRed)
var noteBoxStyle = tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorWhite)
var focusedListStyle = tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorDarkOrange)
var unfocusedListStyle = tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorMediumPurple)
var noteViewBoxStyle = tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorGreen)
var confirmPromptBoxStyle = tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorYellow)
var errBoxStyle = tcell.StyleDefault.Background(tcell.ColorOrangeRed).Foreground(tcell.ColorBlack)
var errTextStyle = tcell.StyleDefault.Background(tcell.ColorOrangeRed).Foreground(tcell.ColorBlack)

// Enum for controlling state
type context int

const (
	ctxMain     context = iota // Main = main screen with note lists. No input boxes, no notes opened.
	ctxNoteView                // NoteView = the screen, in tcell, for viewing full note contents. Not the text editor.
	ctxConfirm
)

var currentCtx = ctxMain

// Enum for confirmation prompt string substitution
type action string

const (
	ActionQuit       action = "quit"
	ActionDeleteNote        = "delete this note"
	ActionDeleteList        = "delete this list"
	ActionSave              = "save"
)

var attemptedAction = action("")

var errMsg = ""

// TODO: There's currently no guarantee that these strings will fit the window. Build them upon window resize.
var keyBindingsStrings = []string{
	" q: Quit, s: Save, Arrows: Change selection, Shift+Arrows: Move selection, u: Refresh, o: OtherCmds ",
	" Notes:: a: Add, e: Edit, d: Delete, v: View, o: OtherCmds ",
	" Lists:: A: Add, r: Rename, D: Delete, o: OtherCmds ",
	" q,v: Back, s: Save, e: Edit, d: Delete, arrows: Change selection, Shift+arrows: Move selection, u: refresh ",
}
var keyBindingsStringIndex = 0

var screenListCap int
var screenNoteCap int

func drawScreen(s tcell.Screen) {
	s.Clear() // Because of the background square, this might not be necessary.

	xmax, ymax := s.Size()
	drawBox(s, 0, 0, xmax-1, ymax-1, boxStyle, "", "") // Background

	kan.draw(s)

	midX, midY := (xmax-1)/2, (ymax-1)/2

	switch currentCtx {
	case ctxNoteView:
		var noteViewW, noteViewH = xmax - 8, ymax - 8
		b := BoundaryBox{midX - (noteViewW / 2), midY - (noteViewH / 2), midX + (noteViewW / 2), midY + (noteViewH / 2)}
		x1, y1, x2, y2 := b.x1, b.y1, b.x2, b.y2
		windowTitle := " Note "
		drawBox(s, x1, y1, x2, y2, noteViewBoxStyle, windowTitle, "")
		drawText(s, x1+2, y1+2, x2-2, y2-2, defStyle, kan.currentNote().Text)
		drawText(s, 5, ymax-1, xmax-1, ymax-1, defStyle, keyBindingsStrings[3])
	case ctxConfirm:
		promptMsg := fmt.Sprintf(" Are you sure you want to %s? [y/N] ", attemptedAction)
		var confirmBoxW, confirmBoxH = len(promptMsg) + 2, 3
		b := BoundaryBox{midX - (confirmBoxW / 2), midY - (confirmBoxH / 2), midX + (confirmBoxW / 2), midY + (confirmBoxH / 2)}
		x1, y1, x2, y2 := b.x1, b.y1, b.x2, b.y2
		drawBox(s, x1, y1, x2, y2, confirmPromptBoxStyle, "", "")
		drawText(s, x1+1, y1+1, x2-1, y2-1, defStyle, promptMsg)
	default:
		drawText(s, 5, ymax-1, xmax-1, ymax-1, defStyle, keyBindingsStrings[keyBindingsStringIndex])
	}

	// Show selected list and note numbers at top right of screen
	currentNoteIndex := kan.curNoteIdx + 1
	if kan.currentList().length() == 0 {
		currentNoteIndex = 0
	}

	noteNumber := fmt.Sprintf(" Note %d/%d ", currentNoteIndex, kan.currentList().length())
	noteNumberXPos := xmax - 2 - len(noteNumber)
	drawText(s, noteNumberXPos, 0, xmax-1, 0, defStyle, noteNumber)

	currentListIndex := kan.curListIdx + 1
	if len(kan.Lists) == 0 {
		currentListIndex = 0
	}
	listNumber := fmt.Sprintf(" List %d/%d ", currentListIndex, len(kan.Lists))
	drawText(s, noteNumberXPos-len(listNumber), 0, xmax-10, 0, defStyle, listNumber)

	if DEBUG_MODE {
		drawBox(s, 1, ymax-5, xmax-2, ymax-2, errBoxStyle, " DEBUG ", errMsg)
	}

	if infoBox.visible {
		infoBox.draw(s)
	}

	s.Show()
}

func defErr() string {
	return ""
}

// fn is expected to be something like list.newNote, list.editNote etc
func openEditorStart(s tcell.Screen, defText string, fn func(string)) {
	// Suspend and Resume are needed to stop text editor from bugging out. Took me too long to figure this out.
	err := s.Suspend()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	newText := openEditor(defText)
	fn(newText)

	err = s.Resume()
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

// Create a temporary file, write string s to the file.
// Open a text editor which either uses the user's EDITOR environment variable or vi.
// After the file is closed by the user, read it and return its contents.
func openEditor(s string) string {
	// I don't know how to hook directly into a text editor, I tried, didn't work, so I won't. Using tempfiles instead.

	file, err := ioutil.TempFile(dir, "note*.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write([]byte(s))
	if err != nil {
		log.Fatalf("%+v", err)
	}
	file.Close() // Closing before editing manually seems like a good idea, right?

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	cmd := exec.Command(editor, file.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	b, err := os.ReadFile(file.Name())
	if err != nil {
		log.Fatalf("%+v", err)
	}

	return string(b)
}

// Initialize tcell.screen
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

// Draw a box with a title on the left-side of the top border
func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, title, contents string) {
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

	titleOffsetX := 2
	drawText(s, x1+titleOffsetX, y1, x1+titleOffsetX+len(title), y1, defStyle, title)
	drawText(s, x1+1, y1+1, x2-1, y2-1, textStyle, contents)
}

func drawErrBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, title, text string) {
	drawBox(s, x1, y1, x2, y2, style, title, text)
}

func maxListsOnScreen(screenWidth int) int {
	listSpace := screenWidth - 2 - listMarginX
	return listSpace / (listWidth + listMarginX)
}

func maxNotesOnScreen(screenHeight int) int {
	// 4 == window border + list border
	noteSpace := screenHeight - 4 - 2*(listMarginY-1) - (noteMargin - 1)
	return noteSpace / (noteHeight + (noteMargin - 1))
}

type InfoBox struct {
	msg     string
	visible bool
}

var infoBox = InfoBox{visible: false}

func showInfoBox(dur time.Duration, msg string) {
	hideInfoBox := func() {
		infoBox.visible = false
		drawScreen(s)
	}
	infoBox.msg = msg
	infoBox.visible = true
	time.AfterFunc(dur, hideInfoBox)
}

func (i *InfoBox) draw(s tcell.Screen) {
	xmax, ymax := s.Size()
	midX, midY := (xmax-1)/2, (ymax-1)/2
	infoMsg := fmt.Sprintf(" %s ", i.msg)
	var infoBoxW, infoBoxH = len(infoMsg) + 2, 3
	x1, y1, x2, y2 := midX-(infoBoxW/2), midY-(infoBoxH/2), midX+(infoBoxW/2), midY+(infoBoxH/2)
	drawBox(s, x1, y1, x2, y2, errBoxStyle, "", "")
	drawText(s, x1+1, y1+1, x2-1, y2-1, errTextStyle, infoMsg)
}
