package gpm

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/go-yaml/yaml"
)

// Options is a go-flags compatible composite structure containing
// the options for each GPM command as well as the CommonOptions across them
type Options struct {
	CommonOptions
	Init   InitOptions   `command:"init"`
	Add    AddOptions    `command:"add"`
	Sync   SyncOptions   `command:"sync"`
	Update UpdateOptions `command:"update"`
	Remove RemoveOptions `command:"remove"`
}

// InitOptions defines the options for the Init command
type InitOptions struct {
	Name       string            `short:"n" long:"name" description:"Project name"`
	License    string            `short:"l" long:"license" description:"Project license (https://spdx.org/licenses/)"`
	Repository string            `short:"r" long:"repo" description:"Project repository"`
	Homepage   string            `short:"h" long:"homepage" description:"Project homepage"`
	Meta       map[string]string `short:"m" long:"meta" description:"Project metadata (key:value pairs)"`
}

// AddOptions defines the options for the Add command
type AddOptions struct {
	Path    string `short:"o" long:"path" description:"Module path"`
	URL     string `short:"u" long:"url" description:"Module URL"`
	Version string `short:"v" long:"version" description:"Module version filter (http://semver.org/)"`
}

// SyncOptions defines the options for the Sync command
type SyncOptions struct{}

// UpdateOptions defines the options for the Update command
type UpdateOptions struct{}

// RemoveOptions defines the options for the Remove command
type RemoveOptions struct {
	Path string `short:"o" long:"path" description:"Module path"`
}

// CommonOptions defines common options for all GPM commands
type CommonOptions struct {
	BasePath  string `short:"c" long:"chdir" description:"Change base directory"`
	NoCleanup bool   `long:"no-cleanup" description:"Disable cleanup of temporary files/folders"`
	Verbose   bool   `long:"verbose" description:"Enable verbose outputs"`
}

// fullPath resolves a filename to a full path using the provided options
func (options *CommonOptions) fullPath(filename string) (string, error) {
	full := fmt.Sprintf("%s/%s", options.BasePath, filename)

	// Check path is inside basepath
	rel, err := filepath.Rel(options.BasePath, full)
	if err != nil {
		return "", err
	}
	if strings.Contains(rel, "..") {
		return "", fmt.Errorf("Invalid full path '%s' (must be within '%s')", full, options.BasePath)
	}

	return full, nil
}

// loadYaml loads a yaml file into an object using the provided options
func (options *CommonOptions) loadYaml(filename string, obj interface{}) error {
	fullpath, err := options.fullPath(filename)
	if err != nil {
		return err
	}

	d, err := ioutil.ReadFile(fullpath)
	if err != nil {
		return fmt.Errorf("Error reading file '%s' (%s)", filename, err)
	}

	err = yaml.Unmarshal(d, obj)
	if err != nil {
		return fmt.Errorf("Error decoding file (%s)", err)
	}

	return nil
}

// writeYaml saves an object into a yaml file using the provided options
func (options *CommonOptions) writeYaml(filename string, obj interface{}) error {
	fullpath, err := options.fullPath(filename)
	if err != nil {
		return err
	}

	b, err := yaml.Marshal(obj)
	if err != nil {
		return fmt.Errorf("Error marshaling object to file '%s' (%s)", filename, err)
	}

	err = ioutil.WriteFile(fullpath, b, 0644)
	if err != nil {
		return fmt.Errorf("Error writing to file '%s' (%s)", filename, err)
	}

	return nil
}
