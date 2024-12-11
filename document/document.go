package document

import (
	"drawino/document/elements" // middle management between documents and rmm
	"drawino/lib/rmm"
)

var (
	CurrentSelected interface{}
	ElementList     map[string]elements.Selectable
)

func LoadDocument() {
	ElementList = map[string]elements.Selectable{
		"ma": nil,
		"ia": nil,
	}
	ElementList["ma"] = elements.NewMessageArea(2, 1, int(rmm.TSize.Width), int(rmm.TSize.Height)-4)
	ElementList["ia"] = elements.NewLargeInputArea(int(rmm.TSize.Height)-2, 1, int(rmm.TSize.Width), 3)
}

func NewDocument() {
	rmm.SetRawMode()
	rmm.ResetTerm()
	LoadDocument()
	for _, v := range ElementList {
		v.Start()
	}

	for {
		for _, v := range ElementList {
			v.Select()
			v.Deselect()
		}
	}
}

func Select(s elements.Selectable) {
	CurrentSelected = s
	s.Select()
}

func Deselect(s elements.Selectable) {
	CurrentSelected = nil
	s.Deselect()
}
