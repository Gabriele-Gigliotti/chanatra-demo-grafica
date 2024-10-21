package main

import (
	"drawino/document" // for document graphics and functionality
	"drawino/rmm"      // for raw mode calculations (e.g. the cursor, the terminal)
	"fmt"
)

func main() {
	defer rmm.ResetTerminalMode()

	document.NewDocument()

	fmt.Print(" Connesso alla Stanza: DEV")

	document.SetCursor(2, 1)
	document.AddBox(int(rmm.TSize.Width), int(rmm.TSize.Height)-4, document.ThinBox)

	document.SetCursor(int(rmm.TSize.Height)-2, 1)
	document.AddBox(int(rmm.TSize.Width), 3, document.ThinBox)

	for {
		document.SetCursor(int(rmm.TSize.Height)-1, 2)
		var a int
		rmm.ScanInt(&a)
	}

}
