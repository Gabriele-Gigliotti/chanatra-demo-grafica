package document

import (
	"drawino/document/elements" // middle management between documents and rmm
	"drawino/lib/rmm"
	"fmt"
)

var (
	CurrentSelected interface{}
)

func LoadDocument() {
	elements.NewMessageArea(2, 1, int(rmm.TSize.Width), int(rmm.TSize.Height)-4)

	lia := elements.NewLargeInputArea(int(rmm.TSize.Height)-2, 1, int(rmm.TSize.Width), 3)
	Select(lia)

	elements.SetCursor(int(rmm.TSize.Height)-1, 2)
	var a string
	fmt.Scan(&a)
}

func NewDocument() {
	rmm.SetRawMode()
	rmm.ResetTerm()
	LoadDocument()
}

func Select(s elements.Selectable) {
	CurrentSelected = s
	s.Select()
}

func Deselect(s elements.Selectable) {
	CurrentSelected = nil
	s.Deselect()
}
