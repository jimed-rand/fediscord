package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func PrintHeader(title string) {
	fmt.Print("\033[H\033[2J")
	width := 59
	padding := width - len(title) - 2
	if padding < 0 {
		padding = 0
	}
	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║  Fediverse to Discord Connection Tool (Mastodon API)     ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════╣")
	fmt.Printf("║  %s%s║\n", title, strings.Repeat(" ", padding))
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()
}

func PrintMenu() {
	fmt.Println("1) Setup Configuration (Discord Token + Fediverse Handle)")
	fmt.Println("2) Generate Connection URL")
	fmt.Println("3) View Stored Configuration")
	fmt.Println("4) Update Discord Token")
	fmt.Println("5) Update Fediverse Handle")
	fmt.Println("6) Change Encryption Settings")
	fmt.Println("7) Delete All Data")
	fmt.Println("8) Exit")
	fmt.Println()
	fmt.Println("───────────────────────────────────────────────────────────")
	fmt.Println("Compatible Platforms (Mastodon API):")
	fmt.Println("  + Mastodon  + Akkoma  + Pleroma  + GlitchSoc  + Hometown")
	fmt.Println("Incompatible Platforms:")
	fmt.Println("  - Misskey  - Firefish  - Calckey  - Foundkey")
	fmt.Println("───────────────────────────────────────────────────────────")
	fmt.Println()
}

func Prompt(label string) string {
	fmt.Print(label)
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func PromptSecret(label string) (string, error) {
	fmt.Print(label)
	raw, err := term.ReadPassword(int(syscall.Stdin))
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
	fmt.Println("[OK] " + msg)
}

func Warn(msg string) {
	fmt.Println("[!!] " + msg)
}

func Error(msg string) {
	fmt.Println("[ERR] " + msg)
}

func Separator() {
	fmt.Println("───────────────────────────────────────────────────────────")
}
