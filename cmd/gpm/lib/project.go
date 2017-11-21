package gpm

// ProjectConfig is a project configuration object from a project file
type ProjectConfig struct {
	Name         string            `yaml:",omitempty"` // Project name
	License      string            `yaml:",omitempty"` // Project license (https://spdx.org/licenses/)
	Repository   string            `yaml:",omitempty"` // Project home repository
	Homepage     string            `yaml:",omitempty"` // Project homepage
	Meta         map[string]string `yaml:",omitempty"` // Project metadata
	Dependencies Dependencies      `yaml:",omitempty"` // List of project dependencies
}

const (
	// ProjectConfigName is the default project config file
	ProjectConfigName = ".gpm.yml"
)
