package main

import (
	"drawino/document"          // for document graphics and functionality
	"drawino/document/elements" // middle management between documents and rmm
	"drawino/lib/rmm"           // for raw mode calculations (e.g. the cursor, the terminal)

	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	defer gracefulStop()
	go gracefulInterrupt()

	document.NewDocument()

	/*
		fmt.Print(" Connesso alla Stanza: DEV")

		elements.SetCursor(2, 1)
		elements.DrawBox(int(rmm.TSize.Width), int(rmm.TSize.Height)-4, elements.ThinBoxType, 0.50)

		elements.SetCursor(int(rmm.TSize.Height)-2, 1)
		elements.DrawBasicBox(int(rmm.TSize.Width), 3, elements.ThinBoxType)

		for {
			elements.SetCursor(int(rmm.TSize.Height)-1, 2)
			var a int
			rmm.ScanInt(&a)
		}
	*/

	ma := elements.NewMessageArea(1, 1, 10, 5)
	mb := elements.NewMessageArea(1, 11, 10, 5)

	var a string

	for {
		elements.SetCursor(11, 1)
		rmm.ScanStrCustom(&a, []rune{'\n', '\r'}, []rune{'\t'}, []rune{})
		fmt.Print("\n", a)

		mb.Deselect()
		ma.Select()

		elements.SetCursor(11, 1)
		rmm.ScanStrCustom(&a, []rune{'\n', '\r'}, []rune{'\t'}, []rune{})

		ma.Deselect()
		mb.Select()
	}

}

func gracefulInterrupt() {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	gracefulStop()
	os.Exit(0)
}

func gracefulStop() {
	rmm.OSClear()
	rmm.ResetTerminalMode()
	fmt.Println("Chanatra was closed successfully.")
}
