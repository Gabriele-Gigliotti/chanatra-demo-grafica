package rmm

import (
	"drawino/lib/logic"

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

// TODO: REDO
func ScanInt(target interface{}) error {
	var inputStr string
	ScanStr(&inputStr)

	if target == nil {
		return nil
	}

	val, err := strconv.Atoi(inputStr)
	if err != nil {
		return fmt.Errorf("failed to parse '%s' as an integer: %v", inputStr, err)
	}

	if intPtr, ok := target.(*int); ok {
		*intPtr = val
	} else {
		return fmt.Errorf("target is not an *int")
	}

	return nil
}

// TODO: REDO
func ScanStr(target *string) error {
	var input strings.Builder
	reader := bufio.NewReader(os.Stdin)

	for {
		b, err := reader.ReadByte()
		if err != nil {
			return err
		}

		if b == '\n' || b == '\r' {
			break
		}

		input.WriteByte(b)
		fmt.Print(string(b))
	}

	*target = input.String()
	return nil
}

func ScanStrCustom(target *string, send []rune, skip []rune, ignore []rune) error {

	var input strings.Builder
	reader := bufio.NewReader(os.Stdin)

	userCursor := 0
	inputLength := 0

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			return err
		}

		if logic.ItemExists(ignore, r) {
			continue
		} else if logic.ItemExists(send, r) {
			break
		}

		// Arrow Keys
		if r == '\x1b' {
			next1, _, _ := reader.ReadRune()
			next2, _, _ := reader.ReadRune()

			if next1 == '[' {
				if next2 == 'C' && userCursor < inputLength {
					// Right arrow key
					fmt.Print("\x1b[C")
					userCursor++
				} else if next2 == 'D' && userCursor > 0 {
					// Left arrow key
					fmt.Print("\x1b[D")
					userCursor--
				}
				continue
			} else {
				reader.UnreadRune()
				reader.UnreadRune()
			}
		}

		// Backspace key
		if r == '\x7f' {
			if userCursor > 0 {
				inputString := input.String()
				// Remove the character before the cursor
				newString := inputString[:userCursor-1] + inputString[userCursor:]
				input.Reset()
				input.WriteString(newString)

				userCursor--
				inputLength--
			}
			fmt.Print("\033[1D \033[1D")
			continue
		}

		fmt.Printf("%c", r)
		input.WriteRune(r)
		userCursor++
		inputLength++
	}

	*target = input.String()
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
