package gpm

import (
	"fmt"
	"log"
	"os"
)

// GPM is the core GoodPackageManager engine
type GPM struct {
	options *CommonOptions
}

// NewGPM creates a new GPM instance with the provided options
func NewGPM(o *CommonOptions) *GPM {
	return &GPM{options: o}
}

// Init initialises a GPM project with the provided ProjectOptions
func (gpm *GPM) Init(po *InitOptions) error {
	if gpm.options.Verbose {
		log.Printf("Initialising project '%s' at '%s' \n", po.Name, gpm.options.BasePath)
	}

	// Build new project information
	pc := ProjectConfig{
		Name:       po.Name,
		License:    po.License,
		Homepage:   po.Homepage,
		Repository: po.Repository,
	}

	// Attempt to load project config
	if _, err := gpm.loadProjectConfig(); err == nil {
		return fmt.Errorf("Project configuration already exists, try `gpm sync`")
	}

	// Validate ProjectOptions
	if po == nil || (po.Name == "" || po.Repository == "") {
		return fmt.Errorf("Project options must include a project name and repository")
	}

	// Save project file
	if err := gpm.writeProjectConfig(&pc); err != nil {
		return err
	}

	// Save lock file
	locks := make(Locks)
	if err := gpm.writeLockfile(&locks); err != nil {
		return err
	}

	if gpm.options.Verbose {
		log.Printf("Init created project config '%s' in dir: '%s'\n", ProjectConfigName, gpm.options.BasePath)
	}

	return nil
}

// Add adds a repository to the current project
func (gpm *GPM) Add(ao *AddOptions) error {
	pc, err := gpm.loadProjectConfig()
	if err != nil {
		return err
	}

	if ao.Path == "" || ao.URL == "" {
		return fmt.Errorf("Module path and url fields cannot be empty")
	}

	if gpm.options.Verbose {
		log.Printf("Adding dependency '%s' with url: '%s' at version: '%s'\n", ao.Path, ao.URL, ao.Version)
	}

	// Determine full module path
	modulePath, err := gpm.options.fullPath(ao.Path)
	if err != nil {
		return err
	}

	// Create and clone a new repository
	repo := NewRepo(modulePath, ao.URL)
	if err := repo.Clone(); err != nil {
		return err
	}

	// Create a tag map
	tags, err := repo.GetTags()
	if err != nil {
		return err
	}

	if gpm.options.Verbose {
		log.Printf("Tags: \n")
		for k, v := range tags {
			log.Printf("\t- %s: %+v", k, v)
		}
	}

	// Fetch latest matching tag if available
	latestTag, latestHash, err := tags.GetLatest(ao.Version)
	if err != nil {
		return err
	}
	if latestTag != "" && latestHash != "" && gpm.options.Verbose {
		log.Printf("Add (%s) Latest matching tag: %s hash: %s", ao.Path, latestTag, latestHash)
	}
	if ao.Version == "" {
		ao.Version = latestTag
	}

	// Sync latest matching hash into worktree
	if err := repo.SyncHash(latestHash); err != nil {
		return err
	}

	// Create new dependency
	d, _ := NewDependency(ao.Path, ao.URL, ao.Version)

	// Update project config file
	pc.Dependencies = append(pc.Dependencies, *d)
	if err := gpm.writeProjectConfig(&pc); err != nil {
		return err
	}

	if gpm.options.Verbose {
		log.Printf("Add (%s) Updated project config '%s' in dir: '%s'\n", ao.Path, ProjectConfigName, gpm.options.BasePath)
	}

	// Update lock file
	locks, err := gpm.loadLockfile()
	if err != nil {
		return err
	}

	locks[ao.Path] = latestHash

	if err := gpm.writeLockfile(&locks); err != nil {
		return err
	}

	if gpm.options.Verbose {
		log.Printf("Add (%s) Updated lockfile '%s' in dir: '%s'\n", ao.Path, LockfileName, gpm.options.BasePath)
	}

	// TODO: Add module path to .gitignore

	return err
}

// Sync pulls and updates dependencies to match lockfile version hashes
func (gpm *GPM) Sync(so *SyncOptions) error {
	pc, err := gpm.loadProjectConfig()
	if err != nil {
		return err
	}
	locks, err := gpm.loadLockfile()
	if err != nil {
		return err
	}

	for _, v := range pc.Dependencies {
		// Resolve full module path
		modulePath, err := gpm.options.fullPath(v.Path)
		if err != nil {
			return err
		}

		// Create repo object to work with
		repo := NewRepo(modulePath, v.URL)

		// Clone or open and update repo depending on current state
		if !repo.Exists() {
			if gpm.options.Verbose {
				log.Printf("Sync (%s) not found, cloning\n", v.Path)
			}
			if err := repo.Clone(); err != nil {
				return err
			}
		} else {
			if gpm.options.Verbose {
				log.Printf("Sync (%s) opening\n", v.Path)
			}
			if err := repo.Open(); err != nil {
				return err
			}
			if gpm.options.Verbose {
				log.Printf("Sync (%s) updating\n", v.Path)
			}
			if err := repo.Fetch(); err != nil {
				return err
			}
		}

		// Locate the matching lock hash
		hash, ok := locks[v.Path]
		if !ok {
			if gpm.options.Verbose {
				log.Printf("Sync (%s) could not find locked hash, try 'gpm update'\n", v.Path)
			}
			return fmt.Errorf("Missing lock hash for module '%s'", v.Path)
		}

		// Sync hash to repo
		if gpm.options.Verbose {
			log.Printf("Sync (%s) checking out hash '%s'\n", v.Path, hash)
		}
		if err := repo.SyncHash(hash); err != nil {
			return err
		}
	}

	// TODO: Check module path exists in .gitignore

	return nil
}

// Update updates lockfile hashes based on the current semver range
func (gpm *GPM) Update(uo *UpdateOptions) error {
	pc, err := gpm.loadProjectConfig()
	if err != nil {
		return err
	}
	locks, err := gpm.loadLockfile()
	if err != nil {
		return err
	}

	for _, v := range pc.Dependencies {
		// Resolve full module path
		modulePath, err := gpm.options.fullPath(v.Path)
		if err != nil {
			return err
		}

		// Create repo object to work with
		repo := NewRepo(modulePath, v.URL)

		// Clone or open and update repo depending on current state
		if !repo.Exists() {
			if gpm.options.Verbose {
				log.Printf("Sync (%s) not found, cloning\n", v.Path)
			}
			if err := repo.Clone(); err != nil {
				return err
			}
		} else {
			if gpm.options.Verbose {
				log.Printf("Sync (%s) opening\n", v.Path)
			}
			if err := repo.Open(); err != nil {
				return err
			}
			if gpm.options.Verbose {
				log.Printf("Sync (%s) updating\n", v.Path)
			}
			if err := repo.Fetch(); err != nil {
				return err
			}
		}

		// Create a tag map
		tags, err := repo.GetTags()
		if err != nil {
			return err
		}

		// Fetch latest matching tag if available
		latestTag, latestHash, err := tags.GetLatest(v.Version)
		if err != nil {
			return err
		}
		if latestTag != "" && latestHash != "" && gpm.options.Verbose {
			log.Printf("Add (%s) Latest matching tag: %s hash: %s", v.Path, latestTag, latestHash)
		}
		if v.Version == "" {
			v.Version = latestTag
		}

		// Sync latest matching hash into worktree
		if err := repo.SyncHash(latestHash); err != nil {
			return err
		}

		locks[v.Path] = latestHash
	}

	if err = gpm.writeLockfile(&locks); err != nil {
		return err
	}

	return nil
}

// Remove removes the specified dependency
func (gpm *GPM) Remove(rm *RemoveOptions) error {
	pc, err := gpm.loadProjectConfig()
	if err != nil {
		return err
	}

	if gpm.options.Verbose {
		log.Printf("Remove (%s) removing module\n", rm.Path)
	}

	// Check dependency exists in the list
	dep, ok := pc.Dependencies.Find(rm.Path)
	if !ok {
		return fmt.Errorf("No dependency bound to location '%s'", rm.Path)
	}

	// Resolve full module path
	modulePath, err := gpm.options.fullPath(rm.Path)
	if err != nil {
		return err
	}

	// TODO: check path is ok
	repo := NewRepo(modulePath, dep.URL)

	// Remove repo from location if it exists
	if repo.Exists() {
		// TODO: check if repo is dirty first

		// Remove repo folder
		err = os.Remove(modulePath)
		if err != nil {
			return err
		}
	}

	// Remove from project config
	pc.Dependencies.Delete(rm.Path)
	err = gpm.writeProjectConfig(&pc)
	if err != nil {
		return err
	}

	// Remove from lockfile
	locks, err := gpm.loadLockfile()
	if err != nil {
		return err
	}
	delete(locks, rm.Path)
	err = gpm.writeLockfile(&locks)
	if err != nil {
		return err
	}

	// TODO: Remove module path from .gitignore

	return nil
}

func (gpm *GPM) loadProjectConfig() (pc ProjectConfig, err error) {
	err = gpm.options.loadYaml(ProjectConfigName, &pc)
	if err != nil && gpm.options.Verbose {
		log.Printf("Error loading project config '%s", ProjectConfigName)
	}
	if pc.Dependencies == nil {
		pc.Dependencies = make(Dependencies, 0)
	}
	return pc, err
}

func (gpm *GPM) writeProjectConfig(pc *ProjectConfig) error {
	return gpm.options.writeYaml(ProjectConfigName, pc)
}

func (gpm *GPM) loadLockfile() (locks Locks, err error) {
	err = gpm.options.loadYaml(LockfileName, &locks)
	if err != nil && gpm.options.Verbose {
		log.Printf("Error loading lock file '%s", LockfileName)
	}
	return locks, err
}

func (gpm *GPM) writeLockfile(locks *Locks) error {
	return gpm.options.writeYaml(LockfileName, locks)
}
