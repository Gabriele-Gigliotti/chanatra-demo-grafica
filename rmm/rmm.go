package rmm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

type TerminalSize struct {
	Height uint16
	Width  uint16
}

var TSize TerminalSize

func init() {
	InitTerminalSize()
}

func InitTerminalSize() error {
	var err error
	TSize, err = GetTermSize()
	if err != nil {
		return err
	}
	return nil
}

func MoveCursor(row, col int) {
	fmt.Printf("\033[%d;%dH", row, col)
}

func ClearTerm() {
	fmt.Print("\033[2J\033[H") // Clear screen and move cursor to home position
}

func ResetTerm() {
	ClearTerm()
	MoveCursor(0, 0)
}

func ScanInt(target *int) error {
	ResetTerminalMode()

	fmt.Scanln(target)

	SetRawMode()
	return nil
}

func ScanStr(target *string) error {
	ResetTerminalMode()

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	input = strings.TrimSuffix(input, "\n")
	*target = input

	SetRawMode()
	return nil
}

func GetCursorPosition() (int, int, error) {
	fmt.Print("\033[6n")

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('R')

	pos := strings.TrimSuffix(strings.TrimPrefix(response, "\033["), "R")

	parts := strings.Split(pos, ";")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("unexpected response: %s", response)
	}

	row, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}

	col, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}

	return row, col, nil
}

func GetTermSize() (TerminalSize, error) {
	fd := int(os.Stdin.Fd())

	var ts TerminalSize

	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&ts)))
	if err != 0 {
		return TerminalSize{}, err
	}

	return ts, nil
}

func SetRawMode() {
	fd := os.Stdin.Fd()
	termios := &syscall.Termios{}
	_, _, _ = syscall.Syscall6(syscall.SYS_IOCTL, fd, syscall.TCGETS, uintptr(unsafe.Pointer(termios)), 0, 0, 0)

	termios.Iflag &^= syscall.ICRNL
	termios.Lflag &^= syscall.ECHO | syscall.ICANON

	_, _, _ = syscall.Syscall6(syscall.SYS_IOCTL, fd, syscall.TCSETS, uintptr(unsafe.Pointer(termios)), 0, 0, 0)
}

func ResetTerminalMode() {
	fd := os.Stdin.Fd()
	termios := &syscall.Termios{}
	_, _, _ = syscall.Syscall6(syscall.SYS_IOCTL, fd, syscall.TCGETS, uintptr(unsafe.Pointer(termios)), 0, 0, 0)

	termios.Iflag |= syscall.ICRNL
	termios.Lflag |= syscall.ECHO | syscall.ICANON

	_, _, _ = syscall.Syscall6(syscall.SYS_IOCTL, fd, syscall.TCSETS, uintptr(unsafe.Pointer(termios)), 0, 0, 0)
}
