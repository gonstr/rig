package context

import (
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/gonstr/rig/pkg/fs"
)

// Context is our interface
type Context interface {
	Scheme() string
	Host() string
	Owner() string
	Repo() string
	Path() string
	Gitref() string
	Digest() string
	RepoURL() (string, error)
	URL() (string, error)
	Values() map[string]interface{}
	OwnerDir() (string, error)
	RepoDir() (string, error)
}

type context struct {
	scheme string
	host   string
	owner  string
	repo   string
	path   string
	gitref string
	digest string
	values map[string]interface{}
}

// FromURL returns a new Context from an url string
func FromURL(urlString string) (Context, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "" {
		return nil, fmt.Errorf("Unable to parse url scheme: %s", urlString)
	}

	if u.Host == "" {
		return nil, fmt.Errorf("Unable to parse url host: %s", urlString)
	}

	if u.Path == "" {
		return nil, fmt.Errorf("Unable to parse url path: %s", urlString)
	}

	gitref := u.Fragment
	if gitref == "" {
		gitref = "master"
	}

	splitPath := strings.Split(u.Path, "/")

	if len(splitPath) < 3 {
		return nil, fmt.Errorf("Invalid git repository url: %s", u.Path)
	}

	owner := splitPath[1]
	repo := splitPath[2]
	path := strings.Join(splitPath[3:], "/")

	return context{scheme: u.Scheme, host: u.Host, owner: owner, repo: repo, path: path, gitref: gitref, digest: "", values: nil}, nil
}

// FromFile returns a new context from a rig file
func FromFile(filePath string) (Context, error) {
	file, err := fs.UnmarshalYaml(filePath)
	if err != nil {
		return nil, err
	}

	template, ok := file["template"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%s is malformed: could not parse template", filePath)
	}

	templatePath, templatePathOk := template["path"].(string)
	templateURL, templateURLOk := template["url"].(string)
	templateGitref, templateGitrefOk := template["gitref"].(string)

	if (!templatePathOk || templatePath == "") && (!templateURLOk || templateURL == "") && (!templateGitrefOk || templateGitref == "") {
		return nil, fmt.Errorf("%s is malformed: does not contain path or url and gitref", filePath)
	}

	templateDigest, _ := template["digest"].(string)

	templateValues, templateValuesOk := file["values"].(map[string]interface{})
	if !templateValuesOk {
		templateValues = make(map[string]interface{})
	}

	if templateURLOk {
		ctx, err := FromURL(templateURL)
		if err != nil {
			return nil, err
		}

		return context{scheme: ctx.Scheme(), host: ctx.Host(), owner: ctx.Owner(), repo: ctx.Repo(), path: ctx.Path(), gitref: ctx.Gitref(), digest: templateDigest, values: templateValues}, nil
	}

	return context{scheme: "", host: "", owner: "", repo: "", path: templatePath, gitref: "", digest: "", values: templateValues}, nil
}

func (c context) Scheme() string {
	return c.scheme
}

func (c context) Host() string {
	return c.host
}

func (c context) Owner() string {
	return c.owner
}

func (c context) Repo() string {
	return c.repo
}

func (c context) Path() string {
	return c.path
}

func (c context) Gitref() string {
	return c.gitref
}

func (c context) Digest() string {
	return c.digest
}

func (c context) URL() (string, error) {
	if c.scheme == "" {
		return "", errors.New("Context contains no URL")
	}

	url := fmt.Sprintf("%s://%s/%s/%s", c.scheme, c.host, c.owner, c.repo)

	if c.path != "" {
		url = fmt.Sprintf("%s/%s", url, c.path)
	}

	return url, nil
}

func (c context) RepoURL() (string, error) {
	if c.scheme == "" {
		return "", errors.New("Context contains no URL")
	}

	url := fmt.Sprintf("%s://%s/%s/%s", c.scheme, c.host, c.owner, c.repo)

	return url, nil
}

func (c context) Values() map[string]interface{} {
	if c.values == nil {
		return make(map[string]interface{})
	}
	return c.values
}

func (c context) OwnerDir() (string, error) {
	if c.scheme == "" {
		return "", errors.New("Can not resolve owner dir since context has no URL")
	}

	homedir, err := fs.HomeDir()
	if err != nil {
		return "", err
	}

	return path.Join(path.Join(homedir, ".rig"), c.host, c.owner), nil
}

func (c context) RepoDir() (string, error) {
	if c.scheme == "" {
		return "", errors.New("Can not resolve repo dir since context has no URL")
	}

	ownerDir, err := c.OwnerDir()
	if err != nil {
		return "", err
	}

	return path.Join(ownerDir, c.repo), nil
}
