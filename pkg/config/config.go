package config

import (
	"os"
	"path/filepath"
	"runtime"
)

type Paths struct {
	Dir            string
	TokenEncrypted string
	TokenPlain     string
	HandleFile     string
	EncryptionFlag string
}

func resolveConfigDirectory() (string, error) {
	switch runtime.GOOS {
	case "windows":
		base := os.Getenv("APPDATA")
		if base == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			base = home
		}
		return filepath.Join(base, "fediverse-discord"), nil
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, "Library", "Application Support", "fediverse-discord"), nil
	default:
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, ".fediverse-discord"), nil
	}
}

func Load() (*Paths, error) {
	dir, err := resolveConfigDirectory()
	if err != nil {
		return nil, err
	}

	return &Paths{
		Dir:            dir,
		TokenEncrypted: filepath.Join(dir, "discord_token.enc"),
		TokenPlain:     filepath.Join(dir, "discord_token.txt"),
		HandleFile:     filepath.Join(dir, "fediverse_handle.txt"),
		EncryptionFlag: filepath.Join(dir, ".use_encryption"),
	}, nil
}

func (p *Paths) Initialise() error {
	return os.MkdirAll(p.Dir, 0700)
}
