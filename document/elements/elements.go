package elements

import (
	"drawino/rmm"
	"errors"
	"fmt"
)

type element struct {
	Row     int
	Col     int
	Width   int
	Height  int
	BoxMode int
}

func (e *element) ApplySelection() {
	//default does nothing
}

func (e *element) Draw() {
	//default does nothing
}

func (e *element) Redraw() {
	//default does nothing
}

func (e *element) SetSelection(boxMode int) {
	e.BoxMode = boxMode
	e.ApplySelection()
}

func (e *element) Select() {
	e.BoxMode = ThickBoxType
	e.ApplySelection()
}

func (e *element) Deselect() {
	e.BoxMode = ThinBoxType
	e.ApplySelection()
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

func (e *MessageArea) ApplySelection() {
	//default does nothing
}

func (e *MessageArea) Draw() {
	SetCursor(e.Row, e.Col)
	DrawBox(e.Width, e.Height, e.BoxMode, e.ScrollPercent)
}

func (e *MessageArea) Redraw() {
	ClearArea(e.Row, e.Col, e.Width, e.Height)
	e.Draw()
}

func (m *MessageArea) ScrollTo(percent float32) {
	m.ScrollPercent = percent
}

func (m *MessageArea) AddMessage(newMessage Message) {
	m.Messages = append(m.Messages, newMessage)
}

type LargeInputArea struct {
	element
}

func NewLargeInputArea(row, col, width, height int) *LargeInputArea {
	return &LargeInputArea{
		//TODO
	}
}

func (e *LargeInputArea) ApplySelection() {
	//default does nothing
}

func (e *LargeInputArea) Draw() {
	//default does nothing
}

func (e *LargeInputArea) Redraw() {
	//default does nothing
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
		SetCursor(CursorRow+1, 1)
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

func ClearArea(row, col, width, height int) {

	for row := CursorRow; row < CursorRow+height; row++ {
		SetCursor(row, col)
		rmm.ClearLine()
	}
}

func DrawBasicBox(width, height int, boxType int) {
	DrawBox(width, height, boxType, -1)
}

func DrawBox(width, height int, boxType int, scrollPercent float32) {
	isThin := boxType == 0
	boxRow, boxCol := CursorRow, CursorCol

	BoxCharCond("--se", "--SE", isThin)
	for col := CursorCol; col < boxCol+width-1; col++ {
		rmm.MoveCursor(CursorRow, col)
		BoxCharCond("-w-e", "-W-E", isThin)
	}
	BoxCharLnCond("-ws-", "-WS-", isThin)

	for row := CursorRow; row < boxRow+height-1; row++ {
		BoxCharCond("n-s-", "N-S-", isThin)
		rmm.MoveCursor(CursorRow, width)

		// Moves cursor to a new line, either way
		if scrollPercent != -1 {
			if row == int(float32(height)*scrollPercent) {
				BoxCharLnCond("C", "D", isThin)
			} else {
				BoxCharLnCond("A", "B", isThin)
			}
		} else {
			BoxCharLnCond("n-s-", "N-S-", isThin)
		}
	}

	BoxCharCond("n--e", "N--E", isThin)
	for col := CursorCol; col < boxCol+width-1; col++ {
		rmm.MoveCursor(CursorRow, col)
		BoxCharCond("-w-e", "-W-E", isThin)
	}
	BoxCharCond("nw--", "NW--", isThin)

	rmm.MoveCursor(boxRow, boxCol)
}
