package fs

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/gonstr/rig/pkg/engine"
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

// UnmarshalYaml reads a path and tries to unmarshal it to yaml
func UnmarshalYaml(path string) (map[string]interface{}, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	bytes, err = engine.Render(string(bytes), nil, false)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})

	err = yaml.Unmarshal(bytes, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// ReadFiles reads all files in a directory and return them as a string array
func ReadFiles(dirOrFilePath string) ([]string, error) {
	filePaths := []string{dirOrFilePath}

	fi, err := os.Stat(dirOrFilePath)
	if err != nil {
		return nil, err
	}

	if fi.Mode().IsDir() {
		filePaths, err = filepath.Glob(path.Join(dirOrFilePath, "*"))
		if err != nil {
			return nil, err
		}
	}

	var contents []string

	for i := 0; i < len(filePaths); i++ {
		content, err := ioutil.ReadFile(filePaths[i])
		if err != nil {
			return nil, err
		}

		contents = append(contents, string(content))
	}

	return contents, nil
}
