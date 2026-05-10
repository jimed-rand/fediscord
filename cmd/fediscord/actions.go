package main

import (
	"errors"
	"fmt"

	"github.com/jimed-rand/fediscord/pkg/config"
	"github.com/jimed-rand/fediscord/pkg/discord"
	"github.com/jimed-rand/fediscord/pkg/fediverse"
	"github.com/jimed-rand/fediscord/pkg/storage"
	"github.com/jimed-rand/fediscord/pkg/ui"
)

func askEncryptionPreference(paths *config.Paths, reuseSaved bool) (bool, error) {
	if reuseSaved {
		if en, err := storage.IsEncryptionEnabled(paths); err == nil {
			return en, nil
		}
	}

	if !storage.IsGPGAvailable() {
		ui.Warn("GPG is not installed. Token will be stored in plain text (INSECURE).")
		ui.Info("  To enable encryption, install GPG:")
		ui.Info("    Ubuntu/Debian: sudo apt install gnupg")
		ui.Info("    Fedora:        sudo dnf install gnupg")
		ui.Info("    Arch:          sudo pacman -S gnupg")
		fmt.Println()
		if !ui.Confirm("Continue with plain text storage? (yes/no): ") {
			return false, errors.New("setup cancelled")
		}
		storage.SetEncryptionPreference(paths, false)
		return false, nil
	}

	ui.Info("Choose Discord token storage method:")
	ui.Info("1) Encrypted (Recommended) - Requires GPG passphrase")
	ui.Info("2) Plain text - No encryption (NOT RECOMMENDED)")
	fmt.Println()

	for {
		choice := ui.Prompt("Select option (1 or 2): ")
		switch choice {
		case "1":
			storage.SetEncryptionPreference(paths, true)
			ui.Success("Encryption enabled")
			return true, nil
		case "2":
			ui.Warn("WARNING: Discord token will be stored in PLAIN TEXT!")
			ui.Warn("Anyone with access to your home directory can read it!")
			fmt.Println()
			if ui.Confirm("Are you absolutely sure? (yes/no): ") {
				storage.SetEncryptionPreference(paths, false)
				ui.Warn("Plain text storage enabled (INSECURE)")
				return false, nil
			}
		default:
			ui.Error("Invalid option. Please choose 1 or 2.")
		}
	}
}

func storeToken(paths *config.Paths, token string, useEncryption bool) error {
	if useEncryption {
		ui.Info("Encrypting and storing Discord token...")
		if err := storage.StoreTokenEncrypted(paths, token); err != nil {
			return fmt.Errorf("failed to encrypt token: %w", err)
		}
		ui.Success("Discord token securely encrypted and stored")
	} else {
		ui.Info("Storing Discord token in plain text...")
		if err := storage.StoreTokenPlain(paths, token); err != nil {
			return fmt.Errorf("failed to store token: %w", err)
		}
		ui.Warn("Discord token stored (UNENCRYPTED - INSECURE)")
	}
	return nil
}

func setupConfiguration(paths *config.Paths) {
	ui.PrintHeader("Setup Configuration")

	ui.Info("This will guide you through setting up:")
	ui.Info("  1. Your Discord account token")
	ui.Info("  2. Your Fediverse handle")
	fmt.Println()

	ui.Info("Step 1: Discord Token")
	ui.Separator()
	ui.Info("To get your Discord token, follow this guide:")
	ui.Info("https://gist.github.com/MarvNC/e601f3603df22f36ebd3102c501116c6")
	fmt.Println()
	ui.Warn("SECURITY WARNING:")
	ui.Warn("  Your Discord token is EXTREMELY sensitive!")
	ui.Warn("  Never share it with anyone or paste it in public places!")
	ui.Warn("  It gives FULL access to your Discord account!")
	fmt.Println()

	useEncryption, err := askEncryptionPreference(paths, true)
	if err != nil {
		ui.Error(err.Error())
		ui.PressEnter()
		return
	}

	token, err := ui.PromptSecret("Enter your Discord token (input hidden): ")
	if err != nil || token == "" {
		ui.Error("Discord token cannot be empty")
		ui.PressEnter()
		return
	}

	if err := storeToken(paths, token, useEncryption); err != nil {
		ui.Error(err.Error())
		ui.PressEnter()
		return
	}
	fmt.Println()

	ui.Info("Step 2: Fediverse Handle")
	ui.Separator()
	ui.Info("Enter your Fediverse handle from a Mastodon API-compatible instance")
	ui.Info("Examples:")
	ui.Info("  @jimedrand@fe.disroot.org (Mastodon)")
	ui.Info("  @user@social.example.com  (Akkoma)")
	ui.Info("  @alice@pleroma.site       (Pleroma)")
	fmt.Println()

	handle := ui.Prompt("Enter your Fediverse handle: ")
	validated, err := fediverse.ValidateHandle(handle)
	if err != nil {
		ui.Error(err.Error())
		ui.PressEnter()
		return
	}

	instance := fediverse.ExtractInstance(validated)
	ui.Success("Instance: " + instance)

	version, err := fediverse.CheckMastodonAPISupport(instance)
	if err != nil {
		ui.Warn(err.Error())
		if !ui.Confirm("Do you want to continue anyway? (yes/no): ") {
			ui.Info("Setup cancelled")
			ui.PressEnter()
			return
		}
	} else {
		ui.Success("Instance is running: " + version)
		ui.Success("Instance appears to support Mastodon API")
	}

	if err := storage.StoreHandle(paths, validated); err != nil {
		ui.Error("Failed to save handle: " + err.Error())
		ui.PressEnter()
		return
	}

	fmt.Println()
	ui.Separator()
	ui.Success("Configuration completed successfully!")
	ui.Success("  Discord token: Stored")
	ui.Success("  Fediverse handle: @" + validated)
	ui.Separator()
	fmt.Println()
	ui.Info("Next step: Use option 2 to generate connection URL")
	ui.PressEnter()
}

func generateConnectionURL(paths *config.Paths) {
	ui.PrintHeader("Generate Connection URL")

	token, err := storage.RetrieveToken(paths)
	if err != nil {
		ui.Error("No Discord token found. Please setup configuration first (Option 1)")
		ui.PressEnter()
		return
	}

	handle, err := storage.RetrieveHandle(paths)
	if err != nil {
		ui.Error("No Fediverse handle found. Please setup configuration first (Option 1)")
		ui.PressEnter()
		return
	}

	ui.Success("Found stored configuration")
	ui.Info("  Fediverse handle: @" + handle)
	fmt.Println()

	authURL, err := discord.GenerateConnectionURL(handle, token)
	if err != nil {
		ui.Error(err.Error())
		ui.Info("")
		ui.Info("Possible reasons:")
		ui.Info("  1. Invalid Discord token")
		ui.Info("  2. Discord API endpoint changed")
		ui.Info("  3. Network connectivity issues")
		ui.PressEnter()
		return
	}

	fmt.Println()
	ui.Separator()
	ui.Success("Authorization URL generated successfully!")
	ui.Separator()
	fmt.Println()
	ui.Info("Authorization URL:")
	fmt.Println()
	fmt.Println(authURL)
	fmt.Println()
	ui.Separator()
	ui.Info("Instructions:")
	ui.Info("  1. Copy the URL above")
	ui.Info("  2. Paste it in your browser")
	ui.Info("  3. Log in to your Fediverse account if needed")
	ui.Info("  4. Authorize the connection")
	ui.Info("  5. Wait for confirmation")
	ui.Separator()
	fmt.Println()
	ui.PressEnter()
}

func viewConfiguration(paths *config.Paths) {
	ui.PrintHeader("Stored Configuration")

	hasConfig := false

	token, err := storage.RetrieveToken(paths)
	if err == nil {
		hasConfig = true
		ui.Success("Discord Token: [STORED]")
		if storage.IsEncryptedTokenPresent(paths) {
			ui.Info("  Storage: Encrypted (GPG) - SECURE")
		} else if storage.IsPlainTokenPresent(paths) {
			ui.Warn("  Storage: Plain text - INSECURE")
		}
		preview := token
		if len(preview) > 10 {
			preview = preview[:10]
		}
		ui.Info("  Preview: " + preview + "...")
	} else {
		ui.Error("Discord Token: [NOT SET]")
	}

	fmt.Println()

	handle, err := storage.RetrieveHandle(paths)
	if err == nil {
		hasConfig = true
		ui.Success("Fediverse Handle: @" + handle)
		ui.Info("  Instance: " + fediverse.ExtractInstance(handle))
	} else {
		ui.Error("Fediverse Handle: [NOT SET]")
	}

	fmt.Println()

	if !hasConfig {
		ui.Separator()
		ui.Warn("No configuration found. Please use Option 1 to setup.")
		ui.Separator()
	}

	fmt.Println()
	ui.PressEnter()
}

func updateDiscordToken(paths *config.Paths) {
	ui.PrintHeader("Update Discord Token")

	ui.Info("This will replace your current Discord token.")
	fmt.Println()

	useEncryption, _ := storage.IsEncryptionEnabled(paths)

	token, err := ui.PromptSecret("Enter new Discord token (input hidden): ")
	if err != nil || token == "" {
		ui.Error("Token cannot be empty")
		ui.PressEnter()
		return
	}

	if err := storeToken(paths, token, useEncryption); err != nil {
		ui.Error(err.Error())
		ui.PressEnter()
		return
	}

	fmt.Println()
	ui.Success("Discord token updated successfully!")
	fmt.Println()
	ui.PressEnter()
}

func updateFediverseHandle(paths *config.Paths) {
	ui.PrintHeader("Update Fediverse Handle")

	ui.Info("This will replace your current Fediverse handle.")
	fmt.Println()

	handle := ui.Prompt("Enter new Fediverse handle: ")
	validated, err := fediverse.ValidateHandle(handle)
	if err != nil {
		ui.Error(err.Error())
		ui.PressEnter()
		return
	}

	instance := fediverse.ExtractInstance(validated)
	ui.Success("Instance: " + instance)

	version, err := fediverse.CheckMastodonAPISupport(instance)
	if err != nil {
		ui.Warn(err.Error())
		if !ui.Confirm("Do you want to continue anyway? (yes/no): ") {
			ui.Info("Update cancelled")
			ui.PressEnter()
			return
		}
	} else {
		ui.Success("Instance is running: " + version)
	}

	if err := storage.StoreHandle(paths, validated); err != nil {
		ui.Error("Failed to save handle: " + err.Error())
		ui.PressEnter()
		return
	}

	fmt.Println()
	ui.Success("Fediverse handle updated to: @" + validated)
	fmt.Println()
	ui.PressEnter()
}

func changeEncryption(paths *config.Paths) {
	ui.PrintHeader("Change Encryption Settings")

	_ = storage.SetEncryptionPreference(paths, false)

	if storage.IsEncryptedTokenPresent(paths) || storage.IsPlainTokenPresent(paths) {
		ui.Warn("Existing Discord token found.")
		ui.Warn("  It will be re-encrypted or decrypted based on your new choice.")
		fmt.Println()

		token, err := storage.RetrieveToken(paths)
		if err != nil {
			ui.Error("Failed to read existing token: " + err.Error())
			ui.PressEnter()
			return
		}

		useEncryption, err := askEncryptionPreference(paths, false)
		if err != nil {
			ui.Error(err.Error())
			ui.PressEnter()
			return
		}

		if err := storeToken(paths, token, useEncryption); err != nil {
			ui.Error(err.Error())
			ui.PressEnter()
			return
		}

		ui.Success("Encryption settings updated successfully!")
	} else {
		ui.Info("No existing Discord token found.")
		fmt.Println()
		_, err := askEncryptionPreference(paths, false)
		if err != nil {
			ui.Error(err.Error())
			ui.PressEnter()
			return
		}
		ui.Success("Encryption preference saved for future tokens")
	}

	fmt.Println()
	ui.PressEnter()
}

func showHelpScreen(paths *config.Paths) {
	ui.PrintHeader("Help")

	ui.Info("What this tool does:")
	ui.Info("  Discord Connections can verify a Mastodon-compatible Fediverse account.")
	ui.Info("  fediscord calls Discord API v9 mastodon authorize and prints the URL you open")
	ui.Info("  in your browser to complete the link.")
	fmt.Println()

	ui.Info("Discord token (advanced):")
	ui.Info("  You need your Discord account user token (what the Discord client uses), not a bot token.")
	ui.Info("  It is equivalent to password-level access; never share or commit it.")
	ui.Info("  How to retrieve it:")
	ui.Info("  https://gist.github.com/MarvNC/e601f3603df22f36ebd3102c501116c6")
	fmt.Println()

	ui.Info("Fediverse handle:")
	ui.Info("  Shape: username@instance.example (leading @ optional).")
	ui.Info("  This tool checks /api/v1/instance; Misskey forks are incompatible.")
	fmt.Println()

	ui.Info("Where data is stored (paths only; no secrets):")
	ui.Separator()
	ui.Info("  Directory:          " + paths.Dir)
	ui.Info("  Token (plain):      " + paths.TokenPlain)
	ui.Info("  Token (encrypted):  " + paths.TokenEncrypted)
	ui.Info("  Handle:             " + paths.HandleFile)
	ui.Info("  Encryption flag:    " + paths.EncryptionFlag)
	ui.Separator()
	fmt.Println()

	ui.Info("Menu: choose 1-8, or type h / ? for this help screen.")
	fmt.Println()
	ui.PressEnter()
}

func deleteAllData(paths *config.Paths) {
	ui.PrintHeader("Delete All Data")

	ui.Warn("WARNING: This will PERMANENTLY delete:")
	ui.Info("  Your Discord token")
	ui.Info("  Your Fediverse handle")
	ui.Info("  All encryption settings")
	ui.Info("  All configuration files")
	fmt.Println()
	ui.Warn("This action CANNOT be undone!")
	fmt.Println()

	confirm := ui.Prompt("Type 'DELETE' to confirm: ")
	if confirm == "DELETE" {
		if err := storage.DeleteAll(paths); err != nil {
			ui.Error("Failed to delete data: " + err.Error())
		} else {
			ui.Success("All data deleted successfully")
			ui.Info("  Configuration directory removed: " + paths.Dir)
		}
	} else {
		ui.Info("Deletion cancelled")
	}

	fmt.Println()
	ui.PressEnter()
}
