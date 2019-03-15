package template

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	gotmpl "text/template"

	"github.com/Masterminds/sprig"
	"github.com/ghodss/yaml"
	"github.com/gonstr/rig/pkg/fs"
	"github.com/gonstr/rig/pkg/git"

	"k8s.io/helm/pkg/chartutil"
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
	Build(filePath string) (string, error)
}

type template struct {
	scheme   string
	host     string
	owner    string
	repo     string
	template string
	version  string
}

// NewFromURI creates a template from a uri
func NewFromURI(uri string) (Template, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "" {
		return nil, fmt.Errorf("Unable to parse url scheme: %s", uri)
	}

	if u.Host == "" {
		return nil, fmt.Errorf("Unable to parse url host: %s", uri)
	}

	if u.Path == "" {
		return nil, fmt.Errorf("Unable to parse url path: %s", uri)
	}

	version := u.Fragment
	if version == "" {
		version = "master"
	}

	split := strings.Split(u.Path, "/")

	if len(split) != 4 {
		return nil, fmt.Errorf("Path does not point to a directy in the repository root: %s", u.Path)
	}

	return template{scheme: u.Scheme, host: u.Host, owner: split[1], repo: split[2], template: split[3], version: version}, nil
}

// NewFromFile returns a template by reading a file in the current working directory
func NewFromFile(path string) (Template, error) {
	m, err := readYaml(path)
	if err != nil {
		return nil, err
	}

	template := m["template"]
	if template == "" {
		return nil, fmt.Errorf("%s is malformed: contains no template", path)
	}

	version := m["version"]
	if version == "" {
		return nil, fmt.Errorf("%s is malformed: contains no version", path)
	}

	return NewFromURI(fmt.Sprintf("%s#%s", template, version))
}

// FuncMap returns funcmap for use in go templating
func FuncMap() gotmpl.FuncMap {
	f := sprig.TxtFuncMap()

	// Add some extra functionality
	extra := gotmpl.FuncMap{
		"toToml":   chartutil.ToToml,
		"toYaml":   chartutil.ToYaml,
		"fromYaml": chartutil.FromYaml,
		"toJson":   chartutil.ToJson,
		"fromJson": chartutil.FromJson,
	}

	for k, v := range extra {
		f[k] = v
	}

	return f
}

func readYaml(path string) (map[string]interface{}, error) {
	bytes, err := ioutil.ReadFile(path)
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

func (t template) ownerDir() (string, error) {
	homedir, err := fs.HomeDir()
	if err != nil {
		return "", err
	}

	return path.Join(path.Join(homedir, ".rig"), t.Host(), t.Owner()), nil
}

func (t template) repoDir() (string, error) {
	ownerDir, err := t.ownerDir()
	if err != nil {
		return "", err
	}

	return path.Join(ownerDir, t.Repo()), nil
}

func (t template) gitSCPURI() string {
	return fmt.Sprintf("git@%s:%s/%s", t.host, t.owner, t.repo)
}

func (t template) templateURL() string {
	return fmt.Sprintf("https://%s/%s/%s/%s", t.host, t.owner, t.repo, t.template)
}

func (t template) Sync() error {
	repoDir, err := t.repoDir()
	if err != nil {
		return err
	}

	ownerDir, err := t.ownerDir()
	if err != nil {
		return err
	}

	if fs.DirExists(repoDir) {
		// Dir exists - we should clean
		err := git.Clean(repoDir)
		if err != nil {
			return err
		}
	} else {
		// Dir does not exists - we should clone
		fs.EnsureDir(ownerDir)
		git.Clone(ownerDir, t.gitSCPURI())
	}

	return nil
}

const rigTmpl = `template: {{ .Template }}
version: {{ .Version }}

values:
  {{ .Values | indent 2 | trim }}
`

func (t template) Install() error {
	repoDir, err := t.repoDir()
	if err != nil {
		return err
	}

	tmpDir, err := fs.TempDir()
	if err != nil {
		return err
	}

	defer os.RemoveAll(tmpDir)

	ref := t.Version()
	if ref != "master" {
		ref = fmt.Sprintf("tags/%s#%s", t.template, t.Version())
	}

	err = git.Checkout(repoDir, tmpDir, ref, t.template)
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	values, err := ioutil.ReadFile(path.Join(tmpDir, t.Template(), "values.yaml"))
	if err != nil {
		return err
	}

	tmpl, err := gotmpl.New("rig").Funcs(FuncMap()).Parse(rigTmpl)
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

func (t template) Build(filePath string) (string, error) {
	m, err := readYaml(filePath)
	if err != nil {
		return "", err
	}

	repoDir, err := t.repoDir()
	if err != nil {
		return "", err
	}

	tmpDir, err := fs.TempDir()
	if err != nil {
		return "", err
	}

	defer os.RemoveAll(tmpDir)

	ref := t.Version()
	if ref != "master" {
		ref = fmt.Sprintf("tags/%s#%s", t.template, t.Version())
	}

	err = git.Checkout(repoDir, tmpDir, ref, t.template)
	if err != nil {
		return "", err
	}

	globPath := path.Join(tmpDir, t.template, "templates", "*")

	filePaths, err := filepath.Glob(globPath)
	if err != nil {
		return "", err
	}

	var strs []string

	for i := 0; i < len(filePaths); i++ {
		content, err := ioutil.ReadFile(filePaths[i])
		if err != nil {
			return "", err
		}

		tmpl, err := gotmpl.New("build").Option("missingkey=error").Funcs(FuncMap()).Parse(string(content))
		if err != nil {
			return "", err
		}

		var buffer bytes.Buffer
		err = tmpl.Execute(&buffer, m)
		if err != nil {
			return "", err
		}

		strs = append(strs, buffer.String())
	}

	joined := strings.Join(strs, "\n---\n")

	re := regexp.MustCompile(`(?m)^\s*$[\r\n]*|[\r\n]+\s+\z`)

	return re.ReplaceAllString(joined, ""), nil
}
