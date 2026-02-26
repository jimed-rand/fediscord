package storage

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/jimed-rand/fediscord/internal/config"
)

var ErrNotFound = errors.New("not found")

func IsGPGAvailable() bool {
	_, err := exec.LookPath("gpg")
	return err == nil
}

func IsEncryptionEnabled(paths *config.Paths) (bool, error) {
	data, err := os.ReadFile(paths.EncryptionFlag)
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(data)) == "true", nil
}

func SetEncryptionPreference(paths *config.Paths, useEncryption bool) error {
	value := "false"
	if useEncryption {
		value = "true"
	}
	return writeFile(paths.EncryptionFlag, []byte(value), 0600)
}

func StoreTokenEncrypted(paths *config.Paths, token string) error {
	cmd := exec.Command("gpg", "--symmetric", "--cipher-algo", "AES256", "--output", paths.TokenEncrypted)
	cmd.Stdin = strings.NewReader(token)
	if err := cmd.Run(); err != nil {
		return err
	}
	if err := os.Chmod(paths.TokenEncrypted, 0600); err != nil {
		return err
	}
	os.Remove(paths.TokenPlain)
	return nil
}

func StoreTokenPlain(paths *config.Paths, token string) error {
	if err := writeFile(paths.TokenPlain, []byte(token), 0600); err != nil {
		return err
	}
	os.Remove(paths.TokenEncrypted)
	return nil
}

func RetrieveToken(paths *config.Paths) (string, error) {
	if fileExists(paths.TokenEncrypted) {
		if !IsGPGAvailable() {
			return "", errors.New("GPG is not available but an encrypted token exists")
		}
		out, err := exec.Command("gpg", "--decrypt", "--quiet", paths.TokenEncrypted).Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(out)), nil
	}

	if fileExists(paths.TokenPlain) {
		data, err := os.ReadFile(paths.TokenPlain)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(data)), nil
	}

	return "", ErrNotFound
}

func StoreHandle(paths *config.Paths, handle string) error {
	return writeFile(paths.HandleFile, []byte(handle), 0600)
}

func RetrieveHandle(paths *config.Paths) (string, error) {
	if !fileExists(paths.HandleFile) {
		return "", ErrNotFound
	}
	data, err := os.ReadFile(paths.HandleFile)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func DeleteAll(paths *config.Paths) error {
	return os.RemoveAll(paths.Dir)
}

func IsEncryptedTokenPresent(paths *config.Paths) bool {
	return fileExists(paths.TokenEncrypted)
}

func IsPlainTokenPresent(paths *config.Paths) bool {
	return fileExists(paths.TokenPlain)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func writeFile(path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}
