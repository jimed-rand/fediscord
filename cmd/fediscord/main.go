package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jimed-rand/fediscord/internal/config"
	"github.com/jimed-rand/fediscord/internal/ui"
)

func main() {
	paths, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration paths: %v\n", err)
		os.Exit(1)
	}

	if err := paths.Initialize(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize config directory: %v\n", err)
		os.Exit(1)
	}

	for {
		ui.PrintHeader("Main Menu")
		ui.PrintMenu()

		choice := ui.Prompt("Select an option (1-8): ")

		switch choice {
		case "1":
			setupConfiguration(paths)
		case "2":
			generateConnectionURL(paths)
		case "3":
			viewConfiguration(paths)
		case "4":
			updateDiscordToken(paths)
		case "5":
			updateFediverseHandle(paths)
		case "6":
			changeEncryption(paths)
		case "7":
			deleteAllData(paths)
		case "8":
			fmt.Print("\033[H\033[2J")
			ui.Separator()
			ui.Info("Thank you for using Fediverse to Discord Connection Tool!")
			ui.Separator()
			os.Exit(0)
		default:
			ui.Error("Invalid option. Please choose 1-8.")
			time.Sleep(2 * time.Second)
		}
	}
}
