package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jimed-rand/fediscord/pkg/terminal"
)

func ClearScreen() {
	fmt.Print(terminal.ClearScreen())
}

func PrintHeader(title string) {
	ClearScreen()
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

func PrintOverview() {
	fmt.Println("This tool builds the Discord authorisation URL used to link your")
	fmt.Println("Fediverse account (Mastodon-style API) to your Discord profile.")
	fmt.Println("Your token and handle are stored only on this computer.")
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
	fmt.Printf("Status  Discord token: %-3s    Fediverse handle: %-3s\n", tok, hdl)

	var next string
	switch {
	case !p.HasToken && !p.HasHandle:
		next = "Choose 1 — Set Up Configuration — to begin."
	case p.HasToken && !p.HasHandle:
		next = "Choose 1 (or 5) to add your Fediverse handle."
	case !p.HasToken && p.HasHandle:
		next = "Choose 1 or 4 to add your Discord token."
	default:
		next = "Choose 2 — Generate Connection URL — then open it in your browser."
	}
	fmt.Println("Next    " + next)
	fmt.Println()
}

func PrintMenu() {
	Separator()
	fmt.Println("Getting started")
	Separator()
	fmt.Println("1) Set Up Configuration (Discord Token + Fediverse Handle)")
	Separator()
	fmt.Println("Daily use")
	Separator()
	fmt.Println("2) Generate Connection URL")
	fmt.Println("3) View Stored Configuration")
	Separator()
	fmt.Println("Maintenance")
	Separator()
	fmt.Println("4) Update Discord Token")
	fmt.Println("5) Update Fediverse Handle")
	fmt.Println("6) Change Encryption Settings")
	Separator()
	fmt.Println("Danger zone")
	Separator()
	fmt.Println("7) Delete All Data")
	Separator()
	fmt.Println("")
	fmt.Println("8) Exit")
	fmt.Println()
	Separator()
	fmt.Println("Compatible Platforms (Mastodon API):")
	fmt.Println("  + Mastodon  + Akkoma  + Pleroma  + GlitchSoc  + Hometown")
	fmt.Println("Incompatible Platforms:")
	fmt.Println("  - Misskey  - Firefish  - Calckey  - Foundkey")
	Separator()
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
