package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jimed-rand/fediscord/pkg/terminal"
)

const dividerWidth = 64

func rule(ch rune) {
	fmt.Println(strings.Repeat(string(ch), dividerWidth))
}

func ClearScreen() {
	fmt.Print(terminal.ClearScreen())
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
	fmt.Println("  Discord can show your Fediverse account on your profile once it has verified")
	fmt.Println("  that you control that identity. Discord asks you to approve that link via a")
	fmt.Println("  special URL (OAuth-style flow hosted by Discord and/or your instance).")
	fmt.Println()
	fmt.Println("How fediscord helps:")
	fmt.Println("  Options 1-3 configure your Discord user token and Fediverse handle on disk.")
	fmt.Println("  Option 2 calls Discord APIs with your saved token only to obtain that link;")
	fmt.Println("  copy it into your browser, sign into your Fediverse account if prompted, then")
	fmt.Println("  approve. fediscord does not run that browser step or store your Mastodon login.")
	fmt.Println()
	fmt.Println("Privacy:")
	fmt.Println("  Your Discord token and Fediverse handle live only on this machine and are")
	fmt.Println("  written under your OS config folder (plain text or GPG, per your choices).")
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
	fmt.Println("Fediverse stacks and Mastodon-compatible API:")
	fmt.Println("  Discord Mastodon Connections expect a Mastodon-style REST API (for example")
	fmt.Println("  `/api/v1/instance`). Your instance must expose that protocol for this tool to")
	fmt.Println("  validate your handle and for linking to succeed.")
	fmt.Println()
	fmt.Println("  Usually okay (known to ship Mastodon API compatibility):")
	fmt.Println("    Mastodon, Akkoma, Pleroma, GlitchSoc, Hometown, and peers that emulate")
	fmt.Println("    the same Mastodon v1 endpoints.")
	fmt.Println()
	fmt.Println("  Not usable here (different API; Misskey-derived, not Mastodon-compatible):")
	fmt.Println("    Misskey, Firefish, Calckey, Foundkey - they use another API layout, so")
	fmt.Println("    ownership checks Discord expects cannot be completed.")
	fmt.Println()
	fmt.Println("  Not sure? Use option 1: the tool probes your instance; you can bail out.")
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
