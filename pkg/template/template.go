package template

import (
	"bytes"
	"fmt"
	gotmpl "html/template"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/Masterminds/sprig"
	"github.com/gonstr/rig/pkg/fs"
	"github.com/gonstr/rig/pkg/git"
)

// Template is our interface
type Template interface {
	Scheme() string
	Host() string
	Owner() string
	Repo() string
	Template() string
	Version() string
	Sync() error
	Install() error
}

type template struct {
	scheme   string
	host     string
	owner    string
	repo     string
	template string
	version  string
	rigDir   string
}

// NewTemplate creates a template from an endpoint
func NewTemplate(endpoint string) (Template, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "" {
		return nil, fmt.Errorf("Unable to parse url scheme: %s", endpoint)
	}

	if u.Host == "" {
		return nil, fmt.Errorf("Unable to parse url host: %s", endpoint)
	}

	if u.Path == "" {
		return nil, fmt.Errorf("Unable to parse url path: %s", endpoint)
	}

	version := u.Fragment
	if version == "" {
		version = "master"
	}

	split := strings.Split(u.Path, "/")

	if len(split) != 4 {
		return nil, fmt.Errorf("Path does not point to a directy in the repository root: %s", u.Path)
	}

	homedir, err := fs.HomeDir()
	if err != nil {
		return nil, err
	}

	rigDir := path.Join(homedir, ".rig")

	return template{scheme: u.Scheme, host: u.Host, owner: split[1], repo: split[2], template: split[3], version: version, rigDir: rigDir}, nil
}

func (t template) Scheme() string {
	return t.scheme
}

func (t template) Host() string {
	return t.host
}

func (t template) Owner() string {
	return t.owner
}

func (t template) Repo() string {
	return t.repo
}

func (t template) Template() string {
	return t.template
}

func (t template) Version() string {
	return t.version
}

func (t template) ownerDir() string {
	return path.Join(t.rigDir, t.Host(), t.Owner())
}

func (t template) repoDir() string {
	return path.Join(t.ownerDir(), t.Repo())
}

func (t template) gitSCPURI() string {
	return fmt.Sprintf("git@%s:%s/%s", t.host, t.owner, t.repo)
}

func (t template) templateURL() string {
	return fmt.Sprintf("https://%s/%s/%s/%s", t.host, t.owner, t.repo, t.template)
}

func (t template) Sync() error {
	if fs.DirExists(t.repoDir()) {
		// Dir exists - we should clean
		err := git.Clean(t.repoDir())
		if err != nil {
			return err
		}
	} else {
		// Dir does not exists - we should clone
		fs.EnsureDir(t.ownerDir())
		git.Clone(t.ownerDir(), t.gitSCPURI())
	}

	return nil
}

const rigTmpl = `template: {{ .Template }}
version: {{ .Version }}

values:
{{ .Values | indent 2 | trim }}
`

func (t template) Install() error {
	// Temp dir
	tmpDir, err := fs.TempDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer os.RemoveAll(tmpDir)

	// Create version
	ref := t.Version()
	if ref != "master" {
		ref = fmt.Sprintf("tags/%s#%s", t.template, t.Version())
	}

	// Copy to temp dir
	err = git.Checkout(t.repoDir(), tmpDir, ref, t.template)
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Read values.yaml to input stream
	values, err := ioutil.ReadFile(path.Join(tmpDir, t.Template(), "values.yaml"))
	if err != nil {
		return err
	}

	tmpl, err := gotmpl.New("rig").Funcs(sprig.FuncMap()).Parse(rigTmpl)
	if err != nil {
		return err
	}

	tmplData := struct {
		Template string
		Version  string
		Values   string
	}{
		Template: t.templateURL(),
		Version:  t.Version(),
		Values:   string(values),
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, tmplData)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(wd, "rig.yaml"), buffer.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}
