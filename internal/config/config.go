package config

import (
	"os"
	"path/filepath"
)

type Paths struct {
	Dir            string
	TokenEncrypted string
	TokenPlain     string
	HandleFile     string
	EncryptionFlag string
}

func Load() (*Paths, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(home, ".fediverse-discord")

	return &Paths{
		Dir:            dir,
		TokenEncrypted: filepath.Join(dir, "discord_token.enc"),
		TokenPlain:     filepath.Join(dir, "discord_token.txt"),
		HandleFile:     filepath.Join(dir, "fediverse_handle.txt"),
		EncryptionFlag: filepath.Join(dir, ".use_encryption"),
	}, nil
}

func (p *Paths) Initialize() error {
	if err := os.MkdirAll(p.Dir, 0700); err != nil {
		return err
	}
	return nil
}
