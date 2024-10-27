package main

import (
	"drawino/document" // for document graphics and functionality
	"drawino/lib/rmm"  // for raw mode calculations (e.g. the cursor, the terminal)

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
