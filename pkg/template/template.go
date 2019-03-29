package template

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gonstr/rig/pkg/engine"
	"github.com/gonstr/rig/pkg/fs"
	"github.com/gonstr/rig/pkg/git"

	"k8s.io/helm/pkg/strvals"
)

const rigTmpl = `template:
  url: {{ .URL }}
  gitref: {{ .Gitref }}
  digest: {{ .Digest }}

values:
  {{ .Values | indent 2 | trim }}
`

// Template is our interface
type Template interface {
	Sync() error
	Install(force bool) error
	Build(values []string, stringValues []string) (string, error)
}

type template struct {
	scheme string
	host   string
	owner  string
	repo   string
	path   string
	gitref string
	digest string
	values map[string]interface{}
}

// NewFromURL creates a template from a uri
func NewFromURL(templateURL string, digest string, values map[string]interface{}) (Template, error) {
	u, err := url.Parse(templateURL)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "" {
		return nil, fmt.Errorf("Unable to parse url scheme: %s", templateURL)
	}

	if u.Host == "" {
		return nil, fmt.Errorf("Unable to parse url host: %s", templateURL)
	}

	if u.Path == "" {
		return nil, fmt.Errorf("Unable to parse url path: %s", templateURL)
	}

	gitref := u.Fragment
	if gitref == "" {
		gitref = "master"
	}

	split := strings.Split(u.Path, "/")

	if len(split) < 3 {
		return nil, fmt.Errorf("Invalid git repository url: %s", u.Path)
	}

	owner := split[1]
	repo := split[2]
	path := strings.Join(split[3:], "/")

	return template{scheme: u.Scheme, host: u.Host, owner: owner, repo: repo, path: path, gitref: gitref, digest: digest, values: values}, nil
}

// NewFromPath creates a template from a local dir path
func newFromPath(path string, digest string, values map[string]interface{}) (Template, error) {
	return template{scheme: "", host: "", owner: "", repo: "", path: path, gitref: "", digest: digest, values: values}, nil
}

// New returns a template by reading a file in the current working directory
func New(filePath string, templatePath string) (Template, error) {
	file, err := fs.ReadYaml(filePath)

	if err != nil {
		// File not found, but that's okey if templatePath is specified
		if templatePath != "" {
			return newFromPath(templatePath, "", make(map[string]interface{}))
		}
		return nil, err
	}

	template, ok := file["template"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%s is malformed: could not parse template", filePath)
	}

	templatePathOk := true
	if templatePath == "" {
		templatePath, templatePathOk = template["path"].(string)
	}

	templateURL, templateURLOk := template["url"].(string)
	templateGitref, templateGitrefOk := template["gitref"].(string)

	if (!templatePathOk || templatePath == "") && (!templateURLOk || templateURL == "") && (!templateGitrefOk || templateGitref == "") {
		return nil, fmt.Errorf("%s is malformed: does not contain path or url and gitref", filePath)
	}

	templateDigest, ok := template["digest"].(string)
	templateValues, ok := file["values"].(map[string]interface{})
	if !ok {
		templateValues = make(map[string]interface{})
	}

	if templatePath != "" {
		return newFromPath(templatePath, templateDigest, templateValues)
	}

	return NewFromURL(fmt.Sprintf("%s#%s", templateURL, templateGitref), templateDigest, templateValues)
}

func (t template) ownerDir() (string, error) {
	homedir, err := fs.HomeDir()
	if err != nil {
		return "", err
	}

	return path.Join(path.Join(homedir, ".rig"), t.host, t.owner), nil
}

func (t template) repoDir() (string, error) {
	ownerDir, err := t.ownerDir()
	if err != nil {
		return "", err
	}

	return path.Join(ownerDir, t.repo), nil
}

func (t template) gitURL() string {
	return fmt.Sprintf("%s://%s/%s/%s", t.scheme, t.host, t.owner, t.repo)
}

func (t template) templateURL() string {
	url := fmt.Sprintf("%s://%s/%s/%s", t.scheme, t.host, t.owner, t.repo)

	if t.path != "" {
		url = fmt.Sprintf("%s/%s", url, t.path)
	}

	return url
}

func (t template) Sync() error {
	if t.scheme == "" {
		return nil
	}

	repoDir, err := t.repoDir()
	if err != nil {
		return err
	}

	ownerDir, err := t.ownerDir()
	if err != nil {
		return err
	}

	if fs.PathExists(repoDir) {
		err := git.Clean(repoDir)
		if err != nil {
			return err
		}
	} else {
		err := fs.EnsureDir(ownerDir)
		if err != nil {
			return err
		}

		err = git.Clone(ownerDir, t.gitURL())
		if err != nil {
			return err
		}
	}

	return nil
}

func (t template) Install(force bool) error {
	repoDir, err := t.repoDir()
	if err != nil {
		return err
	}

	tmpDir, err := fs.TempDir()
	if err != nil {
		return err
	}

	defer os.RemoveAll(tmpDir)

	err = git.Checkout(repoDir, tmpDir, t.gitref, t.path)
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	if fs.PathExists(path.Join(wd, "rig.yaml")) && !force {
		return errors.New("rig.yaml already exists. FORCE install with --force or -f")
	}

	values, err := ioutil.ReadFile(path.Join(tmpDir, t.path, "values.yaml"))
	if err != nil {
		return err
	}

	digest, err := fs.DirectoryDigest(path.Join(tmpDir, t.path, "templates"))

	tmplData := struct {
		URL    string
		Gitref string
		Digest string
		Values string
	}{
		URL:    t.templateURL(),
		Gitref: t.gitref,
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

var emptyLines = regexp.MustCompile(`(?m)^\s*$[\r\n]*|[\r\n]+\s+\z`)
var containsNonWhitespace = regexp.MustCompile(`\S+`)
var onlyASCII = regexp.MustCompile("[[:^ascii:]]")

func (t template) Build(values []string, stringValues []string) (string, error) {
	vals, err := mergeValues(t.values, values, stringValues)
	if err != nil {
		return "", err
	}

	var filePaths []string

	if t.scheme == "" {
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}

		filePaths, err = filepath.Glob(path.Join(wd, t.path, "*"))
		if err != nil {
			return "", err
		}
	} else {
		repoDir, err := t.repoDir()
		if err != nil {
			return "", err
		}

		tmpDir, err := fs.TempDir()
		if err != nil {
			return "", err
		}

		defer os.RemoveAll(tmpDir)

		err = git.Checkout(repoDir, tmpDir, t.gitref, t.path)
		if err != nil {
			return "", err
		}

		if t.digest != "" {
			newdigest, err := fs.DirectoryDigest(path.Join(tmpDir, t.path, "templates"))
			if err != nil {
				return "", err
			}
			if t.digest != newdigest {
				return "", fmt.Errorf("Template digest does not match: %s", newdigest)
			}
		}

		globPath := path.Join(tmpDir, t.path, "templates", "*")

		filePaths, err = filepath.Glob(globPath)
		if err != nil {
			return "", err
		}
	}

	var strs []string

	for i := 0; i < len(filePaths); i++ {
		content, err := ioutil.ReadFile(filePaths[i])
		if err != nil {
			return "", err
		}

		bytes, err := engine.Render(string(content), vals)
		if err != nil {
			return "", err
		}

		str := string(bytes)

		if containsNonWhitespace.MatchString(str) {
			strs = append(strs, str)
		}
	}

	str := strings.Join(strs, "\n---\n")
	str = strings.Replace(str, "<no value>", "", -1)

	return emptyLines.ReplaceAllLiteralString(str, ""), nil
}

func mergeValues(templateValues map[string]interface{}, values []string, stringValues []string) (map[string]interface{}, error) {
	vals := make(map[string]interface{})

	for k, v := range templateValues {
		vals[k] = v
	}

	// User specified a value via --value
	for _, value := range values {
		if err := strvals.ParseInto(value, vals); err != nil {
			return nil, fmt.Errorf("failed parsing --value data: %s", err)
		}
	}

	// User specified a value via --string-value
	for _, value := range stringValues {
		if err := strvals.ParseIntoString(value, vals); err != nil {
			return nil, fmt.Errorf("failed parsing --string-value data: %s", err)
		}
	}

	return map[string]interface{}{
		"values": vals,
	}, nil
}
