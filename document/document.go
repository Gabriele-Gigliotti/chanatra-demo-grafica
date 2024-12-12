package document

import (
	"drawino/document/elements" // middle management between documents and rmm
	"drawino/lib/rmm"
	"time"
)

var (
	CurrentSelected interface{}
	ElementList     map[string]elements.Selectable
)

func LoadDocument() {
	rmm.ResetTerm()
	ElementList = map[string]elements.Selectable{
		"ma": nil,
		"ia": nil,
	}

	messageArea := elements.NewMessageArea(2, 1, int(rmm.TSize.Width), int(rmm.TSize.Height)-4)
	ElementList["ma"] = messageArea
	ElementList["ia"] = elements.NewLargeInputArea(int(rmm.TSize.Height)-2, 1, int(rmm.TSize.Width), 3, messageArea)
}

func NewDocument() {
	rmm.SetRawMode()

	LoadDocument()
	go ResizeWithTerm()

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

func ResizeWithTerm() {
	var formerSize rmm.TerminalSize
	for {
		formerSize = rmm.TSize
		rmm.InitTerminalSize()
		if formerSize != rmm.TSize {
			LoadDocument()
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}
