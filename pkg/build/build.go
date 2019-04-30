package build

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/gonstr/rig/pkg/context"
	"github.com/gonstr/rig/pkg/engine"
	"github.com/gonstr/rig/pkg/fs"
	"github.com/gonstr/rig/pkg/git"
)

var containsNonWhitespace = regexp.MustCompile(`\S+`)

// FromString builds a rig template from a string
func FromString(str string, valueMap map[string]interface{}, values []string, stringValues []string) (string, error) {
	vals, err := createValueMap(valueMap, values, stringValues)
	if err != nil {
		return "", err
	}

	bytes, err := engine.Render(str, vals)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(bytes)), nil
}

// FromTemplatesPath builds a template from a path
func FromTemplatesPath(templatesPath string, valueMap map[string]interface{}, values []string, stringValues []string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	files, err := fs.ReadFiles(path.Join(wd, templatesPath))
	if err != nil {
		return "", err
	}

	var rendered []string
	for i := 0; i < len(files); i++ {
		str, err := FromString(files[i], valueMap, values, stringValues)
		if err != nil {
			return "", err
		}

		if containsNonWhitespace.MatchString(str) {
			rendered = append(rendered, str)
		}
	}

	concatinated := strings.Join(rendered, "\n---\n")
	trimmed := strings.TrimSpace(concatinated)

	return trimmed, nil
}

// FromRigFile builds a template from rig.yaml
func FromRigFile(filePath string, values []string, stringValues []string) (string, error) {
	ctx, err := context.FromFile(filePath)
	if err != nil {
		return "", err
	}

	if ctx.Scheme() == "" {
		return FromTemplatesPath(ctx.Path(), ctx.Values(), values, stringValues)
	}

	ownerDir, err := ctx.OwnerDir()
	if err != nil {
		return "", err
	}

	repoDir, err := ctx.RepoDir()
	if err != nil {
		return "", err
	}

	gitURL, err := ctx.RepoURL()
	if err != nil {
		return "", err
	}

	err = git.Sync(ownerDir, repoDir, gitURL)
	if err != nil {
		return "", err
	}

	tmpDir, err := fs.TempDir()
	if err != nil {
		return "", err
	}

	defer os.RemoveAll(tmpDir)

	err = git.Checkout(repoDir, tmpDir, ctx.Gitref(), ctx.Path())
	if err != nil {
		return "", err
	}

	if ctx.Digest() != "" {
		newdigest, err := fs.DirectoryDigest(path.Join(tmpDir, ctx.Path(), "templates"))
		if err != nil {
			return "", err
		}
		if ctx.Digest() != newdigest {
			return "", fmt.Errorf("Template digest does not match: %s", newdigest)
		}
	}

	files, err := fs.ReadFiles(path.Join(path.Join(tmpDir, ctx.Path(), "templates")))
	if err != nil {
		return "", err
	}

	var rendered []string
	for i := 0; i < len(files); i++ {
		str, err := FromString(files[i], ctx.Values(), values, stringValues)
		if err != nil {
			return "", err
		}

		if containsNonWhitespace.MatchString(str) {
			rendered = append(rendered, str)
		}
	}

	concatinated := strings.Join(rendered, "\n---\n")
	trimmed := strings.TrimSpace(concatinated)

	return trimmed, nil
}
