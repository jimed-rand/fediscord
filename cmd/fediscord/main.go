package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jimed-rand/fediscord/pkg/config"
	"github.com/jimed-rand/fediscord/pkg/storage"
	"github.com/jimed-rand/fediscord/pkg/ui"
)

var version = "dev"

func main() {
	var showHelp bool
	flag.BoolVar(&showHelp, "help", false, "Show usage and exit")
	flag.BoolVar(&showHelp, "h", false, "Show usage and exit (shorthand)")
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "Print version and exit")
	var printConfigPath bool
	flag.BoolVar(&printConfigPath, "print-config-path", false, "Print config directory and file paths, then exit")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Run with no options to open the interactive menu.\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	if len(flag.Args()) > 0 {
		fmt.Fprintf(os.Stderr, "Unknown arguments: %q\n", flag.Args())
		flag.Usage()
		os.Exit(2)
	}
	if showVersion {
		fmt.Println(version)
		os.Exit(0)
	}
	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	paths, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration paths: %v\n", err)
		os.Exit(1)
	}

	if printConfigPath {
		fmt.Println("Config directory:     " + paths.Dir)
		fmt.Println("Discord token (plain):    " + paths.TokenPlain)
		fmt.Println("Discord token (encrypted): " + paths.TokenEncrypted)
		fmt.Println("Fediverse handle file:    " + paths.HandleFile)
		fmt.Println("Encryption preference:  " + paths.EncryptionFlag)
		os.Exit(0)
	}

	if err := paths.Initialize(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize config directory: %v\n", err)
		os.Exit(1)
	}

	for {
		ui.PrintHeader("Main Menu")
		ui.PrintOverview()
		presence := ui.MenuPresence{
			HasToken:  storage.IsEncryptedTokenPresent(paths) || storage.IsPlainTokenPresent(paths),
			HasHandle: storage.IsHandlePresent(paths),
		}
		ui.PrintStatusStrip(presence)
		ui.PrintMenu()

		choice := ui.Prompt("Select an option (1-8, h or ? for help): ")

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
			ui.ClearScreen()
			ui.Separator()
			ui.Info("Thank you for using Fediverse to Discord Connection Tool!")
			ui.Separator()
			os.Exit(0)
		case "h", "H", "?":
			showHelpScreen(paths)
		default:
			ui.Error("Invalid option. Please choose 1-8, or h / ? for help.")
		}
	}
}
