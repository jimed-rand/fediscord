package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jimed-rand/fediscord/pkg/terminal"
)

const dividerWidth = 64

// One reader for stdin so buffered lines stay in sync with term.ReadPassword on the same fd.
var stdin = bufio.NewReader(os.Stdin)

func rule(ch rune) {
	fmt.Println(strings.Repeat(string(ch), dividerWidth))
}

func ClearScreen() {
	if terminal.IsTerminal() {
		fmt.Print(terminal.ClearScreen())
	}
}

func PrintHeader(title string) {
	ClearScreen()
	rule('=')
	fmt.Println("  Fediverse to Discord Connection Tool (Mastodon API)")
	rule('-')
	fmt.Printf("  %s\n", title)
	rule('=')
	fmt.Println()
}

func PrintOverview() {
	fmt.Println("What this does:")
	fmt.Println("  Link a Mastodon-style Fediverse account to Discord: you save a Discord user")
	fmt.Println("  token and your handle here; option 2 prints the URL you open in a browser to")
	fmt.Println("  finish verification on Discord and your instance.")
	fmt.Println()
	fmt.Println("Privacy:")
	fmt.Println("  Token and handle stay on this machine (plain text or GPG, under your OS config")
	fmt.Println("  folder). fediscord does not run the browser step or store your instance login.")
	fmt.Println()
}

type MenuPresence struct {
	HasToken  bool
	HasHandle bool
}

func PrintStatusStrip(p MenuPresence) {
	tok := "no"
	if p.HasToken {
		tok = "yes"
	}
	hdl := "no"
	if p.HasHandle {
		hdl = "yes"
	}
	fmt.Printf("Status: Discord token saved: %s  |  Fediverse handle saved: %s\n", tok, hdl)

	var next string
	switch {
	case !p.HasToken && !p.HasHandle:
		next = "Next: choose 1 (Set Up Configuration) to begin."
	case p.HasToken && !p.HasHandle:
		next = "Next: choose 1 or 5 to add your Fediverse handle."
	case !p.HasToken && p.HasHandle:
		next = "Next: choose 1 or 4 to add your Discord token."
	default:
		next = "Next: choose 2 (Generate Connection URL), then open it in your browser."
	}
	fmt.Println(next)
	fmt.Println()
}

func PrintMenu() {
	fmt.Println("Getting started")
	fmt.Println("  1) Set Up Configuration (Discord Token + Fediverse Handle)")
	fmt.Println()

	fmt.Println("Daily use")
	fmt.Println("  2) Generate Connection URL")
	fmt.Println("  3) View Stored Configuration")
	fmt.Println()

	fmt.Println("Maintenance")
	fmt.Println("  4) Update Discord Token")
	fmt.Println("  5) Update Fediverse Handle")
	fmt.Println("  6) Change Encryption Settings")
	fmt.Println()

	fmt.Println("Danger zone")
	fmt.Println("  7) Delete All Data")
	fmt.Println()

	fmt.Println("Exit")
	fmt.Println("  8) Exit")
	fmt.Println()

	Separator()
	fmt.Println("  Which instances work? Discord expects a Mastodon-compatible API (e.g. /api/v1/instance).")
	fmt.Println("  Mastodon, Akkoma, Pleroma, and similar usually work; Misskey-style APIs do not.")
	fmt.Println("  Press h or ? for full help, paths, and compatibility notes.")
	Separator()
	fmt.Println()
}

func Prompt(label string) string {
	fmt.Print(label)
	line, _ := stdin.ReadString('\n')
	return strings.TrimSpace(line)
}

func PromptSecret(label string) (string, error) {
	fmt.Print(label)
	raw, err := terminal.ReadPassword()
	fmt.Println()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(raw)), nil
}

func Confirm(label string) bool {
	answer := Prompt(label)
	return strings.ToLower(answer) == "yes"
}

func PressEnter() {
	Prompt("Press Enter to continue...")
}

func Info(msg string) {
	fmt.Println(msg)
}

func Success(msg string) {
	fmt.Println("OK - " + msg)
}

func Warn(msg string) {
	fmt.Println("Warning - " + msg)
}

func Error(msg string) {
	fmt.Println("Error - " + msg)
}

func Separator() {
	fmt.Println(strings.Repeat("-", dividerWidth))
}
