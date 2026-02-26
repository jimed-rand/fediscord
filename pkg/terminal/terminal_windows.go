//go:build windows

package terminal

import (
	"golang.org/x/term"
	"os"
)

func ReadPassword() ([]byte, error) {
	return term.ReadPassword(int(os.Stdin.Fd()))
}

func ClearScreen() string {
	return ""
}

func IsTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}
