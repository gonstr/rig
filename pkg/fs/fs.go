package fs

import (
	"os"

	"github.com/mitchellh/go-homedir"
)

// DirExists returns true if a directory path exists
func DirExists(path string) bool {
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
