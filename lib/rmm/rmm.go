package rmm

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"unsafe"
)

type TerminalSize struct {
	Height uint16
	Width  uint16
}

type ScanStrSettings struct {
	send   []rune
	ignore []rune

	deselect []rune
	saved    string

	row int
	col int
}

var (
	TSize TerminalSize
	OS    string
	mu    sync.Mutex
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

func MuLock() {
	mu.Lock()
}
func MuUnLock() {
	mu.Unlock()
}

func MuPrint(row, col int, a ...any) {
	mu.Lock()
	MoveCursor(row, col)
	fmt.Print(a...)
	mu.Unlock()
}

func MuPrintf(row, col int, format string, a ...any) {
	mu.Lock()
	MoveCursor(row, col)
	fmt.Printf(format, a...)
	mu.Unlock()
}

func ScanInt(target interface{}) (int, error) {
	var a string
	status, err := ScanStr(&a)
	if err != nil {
		return status, err
	}

	switch v := target.(type) {
	case *int:
		parsedInt, err := strconv.Atoi(a)
		if err != nil {
			return status, fmt.Errorf("failed to parse integer: %v", err)
		}
		*v = parsedInt
	default:
		return status, fmt.Errorf("invalid target type for ScanInt")
	}

	return status, nil
}

// the double tab problem does not come from here
func ScanTab() {
	reader := bufio.NewReader(os.Stdin)
	for {
		r, _, _ := reader.ReadRune()
		if r == '\t' {
			return
		}
	}

}

func ScanStr(target *string) (int, error) {
	cRow, cCol, _ := GetCursorPosition()
	return ScanStrCustom(target, ScanStrSettings{nil, nil, []rune{'\t'}, "", cRow, cCol})
}

func ScanOrAppendStr(target *string, saved string) (int, error) {
	cRow, cCol, _ := GetCursorPosition()
	return ScanStrCustom(target, ScanStrSettings{nil, nil, []rune{'\t'}, saved, cRow, cCol})
}

func SOASatPos(target *string, saved string, row, col int) (int, error) {
	return ScanStrCustom(target, ScanStrSettings{nil, nil, []rune{'\t'}, saved, row, col})
}

/*
 * STATUS
 * 0 = sent
 * 1 = deselected
 */
func ScanStrCustom(target *string, settings ScanStrSettings) (status int, err error) {
	send := settings.send
	ignore := settings.ignore
	deselect := settings.deselect
	input := []rune(settings.saved)

	var cRow, cCol = settings.row, settings.col
	reader := bufio.NewReader(os.Stdin)

	userCursor := len(input)
	inputLength := len(input)
	var isScannable bool

	for {
		isScannable = true
		var r rune
		r, _, err = reader.ReadRune()
		if err != nil {
			return
		}

		if r == '\n' || r == '\r' || itemExists(send, r) {
			break
		}

		if itemExists(ignore, r) {
			continue
		}

		if itemExists(deselect, r) {
			*target = string(input)
			status = 1
			return
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

						printCursoredString(cRow, cCol, userCursor, input)
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

				printCursoredString(cRow, cCol, userCursor, input)
			}
		}

		if isScannable {
			input = append(input[:userCursor], append([]rune{r}, input[userCursor:]...)...)
			userCursor++
			inputLength++
		}

		printCursoredString(cRow, cCol, userCursor, input)
	}

	*target = string(input) // Set target to the final input string
	MoveCursor(cRow, cCol)
	for i := 0; i < len(string(input)); i++ {
		fmt.Print(" ")
	}
	MoveCursor(cRow, cCol)
	return
}

func printCursoredString(cRow, cCol, userCursor int, input []rune) {
	mu.Lock()
	MoveCursor(cRow, cCol)

	for i := 0; i <= len(input); i++ {
		switch i {
		case len(input):
			if i == userCursor {
				fmt.Print("\x1b[7m")
				fmt.Print(" ")
				fmt.Print("\x1b[27m")
			}
		case userCursor:
			fmt.Print("\x1b[7m")
			fmt.Print(string(input[i]))
			fmt.Print("\x1b[27m")
		default:
			fmt.Print(string(input[i]))
		}
	}
	fmt.Print(" ")
	mu.Unlock()
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
	fmt.Print("\033[?25l")
}

func ResetTerminalMode() {
	fd := os.Stdin.Fd()
	termios := &syscall.Termios{}
	_, _, _ = syscall.Syscall6(syscall.SYS_IOCTL, fd, syscall.TCGETS, uintptr(unsafe.Pointer(termios)), 0, 0, 0)

	termios.Iflag |= syscall.ICRNL
	termios.Lflag |= syscall.ECHO | syscall.ICANON

	_, _, _ = syscall.Syscall6(syscall.SYS_IOCTL, fd, syscall.TCSETS, uintptr(unsafe.Pointer(termios)), 0, 0, 0)
	defer fmt.Print("\033[?25h")
}
