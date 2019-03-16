package git

import (
	"errors"
	"fmt"
	"os/exec"
)

// Clone clones an uri in to dir
func Clone(dir string, uri string) error {
	cmd := exec.Command("git", "clone", uri)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}

	return nil
}

// Clean cleans a directory
func Clean(dir string) error {
	cmd := exec.Command("git", "fetch", "--tags")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}

	cmd = exec.Command("git", "checkout", "master")
	cmd.Dir = dir
	out, err = cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}

	cmd = exec.Command("git", "clean", "-d", "-f")
	cmd.Dir = dir
	out, err = cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}

	cmd = exec.Command("git", "pull", "-q", "origin", "master")
	cmd.Dir = dir
	out, err = cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}

	return nil
}

// Checkout does a git checkout of a local repo/folder to a target directory
func Checkout(repoDir string, targetDir string, ref string, path string) error {
	if path == "" {
		path = "."
	}

	cmd := exec.Command("git", fmt.Sprintf("--work-tree=%s", targetDir), "checkout", ref, "--", path)
	cmd.Dir = repoDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}
	return nil
}
