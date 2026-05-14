package storage

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jimed-rand/fediscord/pkg/config"
)

var ErrNotFound = errors.New("the requested credential was not found in the local configuration store")

func IsGPGAvailable() bool {
	if runtime.GOOS == "windows" {
		return false
	}
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
	if err := ensureParentDir(paths.TokenEncrypted); err != nil {
		return err
	}
	// --yes: avoid overwrite confirmation on stdin (stdin carries only the token).
	cmd := exec.Command("gpg", "--symmetric", "--cipher-algo", "AES256", "--yes", "--output", paths.TokenEncrypted)
	cmd.Stdin = strings.NewReader(token)
	out, err := cmd.CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(out))
		if msg != "" {
			return fmt.Errorf("%w: %s", err, msg)
		}
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
			return "", errors.New("GPG is not available on this system, however an encrypted token was detected; please install GPG to proceed")
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

func IsHandlePresent(paths *config.Paths) bool {
	return fileExists(paths.HandleFile)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func ensureParentDir(path string) error {
	return os.MkdirAll(filepath.Dir(path), 0700)
}

func writeFile(path string, data []byte, perm os.FileMode) error {
	if err := ensureParentDir(path); err != nil {
		return err
	}
	return os.WriteFile(path, data, perm)
}
