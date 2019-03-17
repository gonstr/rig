package fs

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

// PathExists returns true if a directory path exists
func PathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// EnsureDir ensures a directory exists
func EnsureDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// HomeDir return the users home directory
func HomeDir() (string, error) {
	return homedir.Dir()
}

// TempDir return a temp dir
func TempDir() (string, error) {
	dir, err := ioutil.TempDir("", "rig")
	if err != nil {
		return "", err
	}

	return dir, nil
}

// DirectoryDigest returns a sha 256 digest of all files in a directory
func DirectoryDigest(path string) (string, error) {
	hash := sha256.New()

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !info.IsDir() {
			bytes, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			io.WriteString(hash, string(bytes))
		}

		return nil
	})

	if err != nil {
		return "", nil
	}

	return fmt.Sprintf("sha256:%x", hash.Sum(nil)), nil
}
