package document

import (
	"drawino/rmm"
	"errors"
	"fmt"
)

const (
	ThinBox = iota
	ThickBox
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
}

func init() {

}

func NewDocument() {
	rmm.SetRawMode()
	rmm.ResetTerm()
}

func SetCursor(row, col int) {
	rmm.MoveCursor(row, col)
	CursorRow = row
	CursorCol = col
}

func PrintBoxChar(index string) error {
	if char, exists := boxChars[index]; exists {
		fmt.Printf("%s", string(char))
		return nil
	}
	return errors.New("Index " + index + " not found")
}

func PrintlnBoxChar(index string) error {
	if char, exists := boxChars[index]; exists {
		fmt.Printf("%s\n", string(char))
	}
	return errors.New("Index " + index + " not found")
}

func AddBox(width, height int, boxType int) {
	if boxType == ThinBox {
		PrintBoxChar("--se")
		for col := CursorCol + 1; col < width; col++ {
			rmm.MoveCursor(CursorRow, col)
			PrintBoxChar("-w-e")
		}
		PrintlnBoxChar("-ws-")

		for row := CursorCol + 1; row < height; row++ {
			PrintlnBoxChar("n-s-")
			rmm.MoveCursor(CursorRow+row-1, width)
			PrintlnBoxChar("n-s-")
		}

		PrintBoxChar("n--e")
		for col := CursorCol + 1; col < width; col++ {
			rmm.MoveCursor(CursorRow+height-1, col)
			PrintBoxChar("-w-e")
		}
		PrintBoxChar("nw--")

		return
	}
}
