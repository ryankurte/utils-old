package gpm

import ()

// Dependency is a git based project dependency
type Dependency struct {
	Path    string
	URL     string // URL is the git repository URL
	Version string `yaml:",omitempty"` // Version is a semver version or range for matching to git tags
}

// Dependencies are a map of project dependencies
type Dependencies []Dependency

// NewDependency creates a new dependency instance
func NewDependency(path, url, version string) (*Dependency, error) {
	return &Dependency{Path: path, URL: url, Version: version}, nil
}

// Find finds a dependency by path
func (d *Dependencies) Find(path string) (*Dependency, bool) {
	for _, v := range *d {
		if v.Path == path {
			return &v, true
		}
	}
	return nil, false
}

// Delete removes a dependency by path
func (d *Dependencies) Delete(path string) {
	for k, v := range *d {
		if v.Path == path {
			*d = append((*d)[:k], (*d)[k+1:]...)
			break
		}
	}
}

// Set a dependency with the matching path
func (d *Dependencies) Set(path string, dep Dependency) {
	for k, v := range *d {
		if v.Path == path {
			(*d)[k] = dep
			break
		}
	}
}
