package engine

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"text/template"

	"github.com/Masterminds/sprig"
	"k8s.io/helm/pkg/chartutil"
)

// FuncMap returns funcmap for use in go templating
func FuncMap() template.FuncMap {
	f := sprig.TxtFuncMap()

	// Add some extra functionality
	extra := template.FuncMap{
		"toToml":   chartutil.ToToml,
		"toYaml":   chartutil.ToYaml,
		"fromYaml": chartutil.FromYaml,
		"toJson":   chartutil.ToJson,
		"fromJson": chartutil.FromJson,

		// We want to error on env or expandenv if the env values does no exist
		"env": func(s string) (string, error) {
			e := os.Getenv(s)
			if e == "" {
				return "", fmt.Errorf("Environment variable '%s' does not exists", s)
			}
			return e, nil
		},
		"expandenv": func(s string) (string, error) {
			e := os.ExpandEnv(s)
			if e == "" {
				return "", fmt.Errorf("Environment variable '%s' does not exists", s)
			}
			return e, nil
		},
	}

	for k, v := range extra {
		f[k] = v
	}

	return f
}

// Render a tmp string with templates values
func Render(str string, vals interface{}) ([]byte, error) {
	// Go templates fails to render funny unicode characters to we replace any non
	// ascii characters with an empty string
	str = regexp.MustCompile("[[:^ascii:]]").ReplaceAllLiteralString(str, "")

	tmpl, err := template.New("tmpl").Option("missingkey=error").Funcs(FuncMap()).Parse(str)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, vals)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
