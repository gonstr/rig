package git

import (
	"fmt"
	"os/exec"
)

// Clone clones an uri in to dir
func Clone(dir string, uri string) error {
	cmd := exec.Command(fmt.Sprintf("git clone %s", uri))
	cmd.Dir = dir
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// Clean cleans a directory
func Clean(dir string) error {
	cmd := exec.Command("git", "fetch", "--tags")
	cmd.Dir = dir
	err := cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("git", "checkout", "master")
	cmd.Dir = dir
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("git", "clean", "-d", "-f")
	cmd.Dir = dir
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("git", "pull", "-q", "origin", "master")
	cmd.Dir = dir
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
