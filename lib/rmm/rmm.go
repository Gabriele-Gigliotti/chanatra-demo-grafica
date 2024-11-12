package rmm

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

type TerminalSize struct {
	Height uint16
	Width  uint16
}

var (
	TSize TerminalSize
	OS    string
)

func init() {
	InitTerminalSize()

	OS = runtime.GOOS
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

func ClearLine() {
	fmt.Print("\033[2K")
}

func ResetTerm() {
	fmt.Print("\033[2J\033[H")
}

func OSClear() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd = exec.Command("clear")
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		// tries to clear with ANSI
		fmt.Print("\033[2J\033[H")
		return
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func ScanInt(target interface{}) error {
	var a string
	err := ScanStr(&a)
	if err != nil {
		return err
	}

	switch v := target.(type) {
	case *int:
		parsedInt, err := strconv.Atoi(a)
		if err != nil {
			return fmt.Errorf("failed to parse integer: %v", err)
		}
		*v = parsedInt
	default:
		return fmt.Errorf("invalid target type for ScanInt")
	}

	return nil
}

func ScanStr(target *string) error {
	return ScanStrCustom(target, nil, []rune{'\t'})
}

func ScanStrCustom(target *string, send []rune, ignore []rune) error {

	var cRow, cCol, _ = GetCursorPosition()
	var input []rune
	reader := bufio.NewReader(os.Stdin)

	userCursor := 0
	inputLength := 0
	var isScannable bool

	for {
		isScannable = true
		r, _, err := reader.ReadRune()
		if err != nil {
			return err
		}

		if r == '\n' || r == '\r' || itemExists(send, r) {
			break
		}

		if itemExists(ignore, r) {
			continue
		}

		// Arrow keys and delete
		if r == '\x1b' {
			isScannable = false
			next1, _, _ := reader.ReadRune()
			next2, _, _ := reader.ReadRune()

			if next1 == '[' {
				switch next2 {
				case 'C': // Right
					if userCursor < inputLength {
						fmt.Print("\x1b[C")
						userCursor++
					}
				case 'D': // Left
					if userCursor > 0 {
						fmt.Print("\x1b[D")
						userCursor--
					}
				case '3': // Delete key (ESC [3~ sequence)
					next3, _, _ := reader.ReadRune()
					if next3 == '~' && userCursor < inputLength {
						input = append(input[:userCursor], input[userCursor+1:]...)
						inputLength--

						MoveCursor(cRow, cCol)
						fmt.Print(string(input) + " ")
					}
				}
			}
		}

		// Backspace
		if r == '\x7f' {
			isScannable = false
			if userCursor > 0 {
				input = append(input[:userCursor-1], input[userCursor:]...)
				userCursor--
				inputLength--

				MoveCursor(cRow, cCol)
				fmt.Print(string(input) + " ")
			}
		}

		if isScannable {
			input = append(input[:userCursor], append([]rune{r}, input[userCursor:]...)...)
			userCursor++
			inputLength++
		}

		MoveCursor(cRow, cCol)
		fmt.Print(string(input))
		MoveCursor(cRow, cCol+userCursor)
	}

	*target = string(input) // Set target to the final input string
	MoveCursor(cRow, cCol)
	return nil
}

func itemExists(list []rune, item rune) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
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
