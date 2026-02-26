//go:build !windows

package terminal

import (
	"golang.org/x/term"
	"os"
)

func ReadPassword() ([]byte, error) {
	return term.ReadPassword(int(os.Stdin.Fd()))
}

func ClearScreen() string {
	return "\033[H\033[2J"
}

func IsTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}
