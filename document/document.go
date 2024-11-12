package document

import (
	"drawino/document/elements" // middle management between documents and rmm
	"drawino/lib/rmm"
	"fmt"
)

var (
	CurrentSelected interface{}
	ElementList     map[string]elements.Selectable
)

func LoadDocument() {
	ElementList["ma"] = elements.NewMessageArea(2, 1, int(rmm.TSize.Width), int(rmm.TSize.Height)-4)
	ElementList["ia"] = elements.NewLargeInputArea(int(rmm.TSize.Height)-2, 1, int(rmm.TSize.Width), 3)

	Select(ElementList["ia"])

	elements.SetCursor(int(rmm.TSize.Height)-1, 2)
	var a string
	rmm.ScanStr(&a)

	rmm.ResetTerm()
	fmt.Println(a)
	for i := 0; i < len(a); i++ {
		fmt.Printf("%d:\t%08b\t%d\t%c\n", i, a[i], a[i], a[i])
	}
	fmt.Scan(&a)
}

func NewDocument() {
	rmm.SetRawMode()
	rmm.ResetTerm()
	ElementList = map[string]elements.Selectable{
		"ma": nil,
		"ia": nil,
	}
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
