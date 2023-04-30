package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

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

type context int

const (
	ctxMain     context = iota // Main = main screen with note lists. No input boxes, no notes opened.
	ctxNoteView                // NoteView = the screen, in tcell, for viewing full note contents. Not the text editor.
	ctxConfirm
)

var currentCtx = ctxMain

type action string

const (
	ActionQuit       action = "quit"
	ActionDeleteNote        = "delete this note"
	ActionDeleteList        = "delete this list"
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

func drawScreen(s tcell.Screen) {
	s.Clear() // Because of the background square, this might not be necessary.
	xmax, ymax := s.Size()
	drawBox(s, 0, 0, xmax-1, ymax-1, boxStyle, "", "") // Background
	kan.draw(s)

	// TODO: edit keybinding strings

	switch currentCtx {
	case ctxNoteView:
		var noteViewW, noteViewH = xmax - 8, ymax - 8
		left := (xmax-1)/2 - (noteViewW / 2)
		right := (xmax-1)/2 + (noteViewW / 2)
		top := (ymax-1)/2 - (noteViewH / 2)
		bottom := (ymax-1)/2 + (noteViewH / 2)
		windowTitle := " Note "
		drawBox(s, left, top, right, bottom, noteViewBoxStyle, windowTitle, "")
		// drawText(s, left+2, top, left+2+len(windowTitle), top, defStyle, windowTitle)
		drawText(s, left+2, top+2, right-2, bottom-2, defStyle, kan.currentNote().Text)
		drawText(s, 5, ymax-1, xmax-1, ymax-1, defStyle, keyBindingsStrings[3])
	case ctxConfirm:
		promptMsg := fmt.Sprintf(" Are you sure you want to %s? [y/N] ", attemptedAction)
		var confirmBoxW, confirmBoxH = len(promptMsg) + 2, 3
		left := (xmax-1)/2 - (confirmBoxW / 2)
		right := (xmax-1)/2 + (confirmBoxW / 2)
		top := (ymax-1)/2 - (confirmBoxH / 2)
		bottom := (ymax-1)/2 + (confirmBoxH / 2)
		drawBox(s, left, top, right, bottom, confirmPromptBoxStyle, "", "")
		drawText(s, left+1, top+1, right-1, bottom-1, defStyle, promptMsg)
	default:
		drawText(s, 5, ymax-1, xmax-1, ymax-1, defStyle, keyBindingsStrings[keyBindingsStringIndex])
	}

	if DEBUG_MODE {
		drawBox(s, 1, ymax-5, xmax-2, ymax-2, errBoxStyle, " DEBUG ", errMsg)
	}
}

func defErr() string {
	return ""
}

// fn is expected to be things like list.newNote, list.editNote etc
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
	// log.Printf("Using %s", editor)

	cmd := exec.Command(editor, file.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	// log.Printf("Running %s...", editor)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	b, err := os.ReadFile(file.Name())
	if err != nil {
		log.Fatalf("%+v", err)
	}

	return string(b)
}

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
