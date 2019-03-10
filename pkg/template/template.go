package template

import (
	"fmt"
	"net/url"
	"path"
	"strings"

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
	Sync() error
}

type template struct {
	scheme   string
	host     string
	owner    string
	repo     string
	template string
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

	split := strings.Split(u.Path, "/")

	if len(split) != 2 {
		return nil, fmt.Errorf("Path does not point to a directy in the repository root: %s", u.Path)
	}

	homedir, err := fs.HomeDir()
	if err != nil {
		return nil, err
	}

	rigDir := path.Join(homedir, ".rig")

	return template{scheme: u.Scheme, host: u.Host, owner: split[0], template: split[1], rigDir: rigDir}, nil
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

func (t template) ownerDir() string {
	return path.Join(t.rigDir, t.Host(), t.Owner())
}

func (t template) repoDir() string {
	return path.Join(t.ownerDir(), t.Repo())
}

func (t template) gitSCPURI() string {
	return fmt.Sprintf("git@%s:%s/%s", t.host, t.owner, t.repo)
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
		git.Clone(t.ownerDir(), t.gitSCPURI())
	}

	return nil
}
