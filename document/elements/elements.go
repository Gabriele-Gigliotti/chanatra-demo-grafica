package elements

import (
	"drawino/lib/rmm"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Selectable interface {
	ApplySelection()
	Select() Selectable
	Deselect() Selectable
	Start()
	Draw()
	Redraw()
}

type element struct {
	Row     int
	Col     int
	Width   int
	Height  int
	BoxType int
}

type Message struct {
	Author  string
	Message string
}

type MessageArea struct {
	element
	ScrollPercent float32
	Messages      []Message
}

func MakeMessage(author, message string) Message {
	return Message{
		Author:  author,
		Message: message,
	}
}

func NewMessageArea(row, col, width, height int) *MessageArea {
	ma := &MessageArea{
		element: element{
			Row:    row,
			Col:    col,
			Width:  width,
			Height: height,
		},

		ScrollPercent: 0,
		Messages:      []Message{},
	}

	ma.Draw()
	return ma
}

func (e *MessageArea) Redraw() {
	ClearArea(e.Row, e.Col, e.Width, e.Height)
	e.Draw()
}

func (e *MessageArea) Draw() {
	SetCursor(e.Row, e.Col)
	DrawBox(e.Width, e.Height, e.BoxType, e.ScrollPercent)

	rmm.MuPrint(4, 2, e.Messages)
	//TODO: draw messages
}

func (e *MessageArea) ApplySelection() {
	e.Redraw()
	rmm.ScanTab()
}

func (e *MessageArea) Select() Selectable {
	e.BoxType = ThickBoxType
	e.ApplySelection()
	return e
}

func (e *MessageArea) Deselect() Selectable {
	e.BoxType = ThinBoxType
	e.Redraw()
	return e
}

func (m *MessageArea) ScrollTo(percent float32) {
	m.ScrollPercent = percent
	m.Redraw()
}

func (m *MessageArea) AddMessage(newMessage Message) {
	m.Messages = append(m.Messages, newMessage)
	m.Redraw()
}

func (e *MessageArea) Start() {
	go e.doStuff()
}

func (e *MessageArea) doStuff() {
	for {
		rmm.MuPrint(e.Row+1, e.Col+1, "|")
		time.Sleep(100 * time.Millisecond)

		rmm.MuPrint(e.Row+1, e.Col+1, "/")
		time.Sleep(100 * time.Millisecond)

		rmm.MuPrint(e.Row+1, e.Col+1, "-")
		time.Sleep(100 * time.Millisecond)

		rmm.MuPrint(e.Row+1, e.Col+1, "\\")
		time.Sleep(100 * time.Millisecond)
	}
}

type LargeInputArea struct {
	element
	SavedStr string
	Ma       *MessageArea
}

func NewLargeInputArea(row, col, width, height int, ma *MessageArea) *LargeInputArea {
	lia := &LargeInputArea{
		element: element{
			Row:    row,
			Col:    col,
			Width:  width,
			Height: height,
		},
		SavedStr: "",
		Ma:       ma,
	}

	lia.Draw()
	return lia
}

func (e *LargeInputArea) Start() {
	// Do nothing
}

func (e *LargeInputArea) Draw() {
	SetCursor(e.Row, e.Col)
	DrawBox(e.Width, e.Height, e.BoxType, -1)
}

func (e *LargeInputArea) ApplySelection() {
	e.Redraw()
	var status int
	var a string

	// Repeat input if sent
	for status == 0 {
		status, _ = rmm.SOASatPos(&a, e.SavedStr, e.Row+1, e.Col+1)
		if status == 0 {
			e.Ma.AddMessage(MakeMessage("test", a))
		}
		e.SavedStr = ""
	}

	// If deselected, remember input
	e.SavedStr = a
	e.Deselect()
}

func (e *LargeInputArea) Select() Selectable {
	e.BoxType = ThickBoxType
	e.ApplySelection()
	return e
}

func (e *LargeInputArea) Deselect() Selectable {
	e.BoxType = ThinBoxType
	e.Redraw()
	return e
}

func (e *LargeInputArea) Redraw() {
	//ClearArea(e.Row, e.Col, e.Width, e.Height)
	e.Draw()
}

// ------- Functions ------- //

const (
	ThinBoxType = iota
	ThickBoxType
)

var (
	CursorRow int = 0
	CursorCol int = 0
)

// Strictly for making the drawing code more readable. If you can call this readable.
var boxChars = map[string]rune{
	"-w-e": 0x2500, // ─
	"n-s-": 0x2502, // │
	"--se": 0x250C, // ┌
	"-ws-": 0x2510, // ┐
	"n--e": 0x2514, // └
	"nw--": 0x2518, // ┘

	"-W-E": 0x2550, // ═
	"N-S-": 0x2551, // ║
	"--SE": 0x2554, // ╔
	"-WS-": 0x2557, // ╗
	"N--E": 0x255A, // ╚
	"NW--": 0x255D, // ╝

	"A": 0x2591, //░
	"B": 0x2592, //▒
	"C": 0x2593, //▓
	"D": 0x2588, //█
}

func SetCursor(row, col int) {
	rmm.MoveCursor(row, col)
	CursorRow = row
	CursorCol = col
}

func BoxChar(index string) error {
	if char, exists := boxChars[index]; exists {
		fmt.Printf("%s", string(char))
		CursorCol = CursorCol + 1
		return nil
	}
	return errors.New("Index " + index + " not found")
}

func BoxCharLn(index string) error {
	if char, exists := boxChars[index]; exists {
		fmt.Print(string(char))
		SetCursor(CursorRow+1, 0)
	}
	return errors.New("Index " + index + " not found")
}

func BoxCharCond(index1 string, index2 string, condition bool) error {
	var err error

	if condition {
		err = BoxChar(index1)
	} else {
		err = BoxChar(index2)
	}

	return err
}

func BoxCharLnCond(index1 string, index2 string, condition bool) error {
	var err error

	if condition {
		err = BoxCharLn(index1)
	} else {
		err = BoxCharLn(index2)
	}

	return err
}

func ClearArea(areaRow, areaCol, width, height int) {
	rmm.MuLock()
	defer rmm.MuUnLock()

	var builder strings.Builder
	builder.WriteString(strings.Repeat(" ", width))
	str := builder.String()

	rmm.MoveCursor(areaRow, areaCol)
	for i := 0; i < height; i++ {
		fmt.Print(str)
		if i < height-1 {
			rmm.MoveCursor(areaRow+i, areaCol)
		}
	}
}

func DrawBasicBox(width, height int, boxType int) {
	DrawBox(width, height, boxType, -1)
}

func DrawBox(width, height int, boxType int, scrollPercent float32) {
	rmm.MuLock()
	defer rmm.MuUnLock()

	isThin := boxType == 0
	boxRow, boxCol := CursorRow, CursorCol

	BoxCharCond("--se", "--SE", isThin)
	for col := boxCol + 1; col < boxCol+width-1; col++ {
		rmm.MoveCursor(CursorRow, col)
		BoxCharCond("-w-e", "-W-E", isThin)
	}
	BoxCharLnCond("-ws-", "-WS-", isThin)

	for row := CursorRow; row < boxRow+height-1; row++ {
		rmm.MoveCursor(CursorRow, boxCol)
		BoxCharCond("n-s-", "N-S-", isThin)
		rmm.MoveCursor(CursorRow, boxCol+width-1)

		scrollbarPos := boxRow + 1 + (int(float32(height) * scrollPercent))

		// Moves cursor to a new line, either way
		if scrollPercent != -1 {
			if row == int(scrollbarPos) {
				BoxCharLnCond("C", "D", isThin)
			} else {
				BoxCharLnCond("A", "B", isThin)
			}
		} else {
			BoxCharLnCond("n-s-", "N-S-", isThin)
		}
	}

	rmm.MoveCursor(CursorRow, boxCol)
	BoxCharCond("n--e", "N--E", isThin)

	for col := boxCol + 1; col < boxCol+width-1; col++ {
		rmm.MoveCursor(CursorRow, col)
		BoxCharCond("-w-e", "-W-E", isThin)
	}
	BoxCharCond("nw--", "NW--", isThin)

	rmm.MoveCursor(boxRow, boxCol)
}
