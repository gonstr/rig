package install

import (
	"errors"
	"io/ioutil"
	"os"
	"path"

	"github.com/gonstr/rig/pkg/context"
	"github.com/gonstr/rig/pkg/engine"
	"github.com/gonstr/rig/pkg/fs"
	"github.com/gonstr/rig/pkg/git"
)

const rigTmpl = `template:
  url: {{ .URL }}
  gitref: {{ .Gitref }}
  digest: {{ .Digest }}

values:
  {{ .Values | indent 2 | trim }}
`

// FromURL installs a rig template from an url
func FromURL(url string, force bool) error {
	ctx, err := context.FromURL(url)
	if err != nil {
		return err
	}

	ownerDir, err := ctx.OwnerDir()
	if err != nil {
		return err
	}

	repoDir, err := ctx.RepoDir()
	if err != nil {
		return err
	}

	gitURL, err := ctx.RepoURL()
	if err != nil {
		return err
	}

	err = git.Sync(ownerDir, repoDir, gitURL)
	if err != nil {
		return err
	}

	tmpDir, err := fs.TempDir()
	if err != nil {
		return err
	}

	defer os.RemoveAll(tmpDir)

	err = git.Checkout(repoDir, tmpDir, ctx.Gitref(), ctx.Path())
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	if fs.PathExists(path.Join(wd, "ctx.yaml")) && !force {
		return errors.New("ctx.yaml already exists. FORCE install with --force or -f")
	}

	values, err := ioutil.ReadFile(path.Join(tmpDir, ctx.Path(), "values.yaml"))
	if err != nil {
		return err
	}

	digest, err := fs.DirectoryDigest(path.Join(tmpDir, ctx.Path(), "templates"))

	fullURL, err := ctx.URL()
	if err != nil {
		return err
	}

	tmplData := struct {
		URL    string
		Gitref string
		Digest string
		Values string
	}{
		URL:    fullURL,
		Gitref: ctx.Gitref(),
		Digest: digest,
		Values: string(values),
	}

	bytes, err := engine.Render(rigTmpl, tmplData)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(wd, "rig.yaml"), bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
