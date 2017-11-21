package gpm

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type Repo struct {
	path       string
	url        string
	repository *git.Repository
}

// NewRepo creates a new repo instance
func NewRepo(path, url string) *Repo {
	return &Repo{path: path, url: url, repository: nil}
}

// Exists checks if a repo exists on disk
func (repo *Repo) Exists() bool {
	// Check folder exists

	// Attempt loading repository and config
	r, err := git.PlainOpen(repo.path)
	if err != nil {
		return false
	}
	_, err = r.Config()
	if err != nil {
		return false
	}

	return true
}

// Clone populates a new repo on disk
func (repo *Repo) Clone() error {
	cloneOpts := git.CloneOptions{URL: repo.url, Depth: 1, Tags: git.AllTags}
	r, err := git.PlainClone(repo.path, false, &cloneOpts)
	if err != nil {
		return err
	}
	repo.repository = r
	return nil
}

// Open loads a repo from the provided path
func (repo *Repo) Open() error {
	r, err := git.PlainOpen(repo.path)
	if err != nil {
		return err
	}

	repo.repository = r
	return nil
}

// Sync updates a repo and remotes
func (repo *Repo) Sync() error {
	_, err := repo.repository.Config()
	if err != nil {
		return err
	}

	// Update remote if mismatched
	remote, err := repo.repository.Remote("origin")
	origin := config.RemoteConfig{Name: "origin", URLs: []string{repo.url}}
	origin.Validate()
	if err == nil && remote.Config().URLs[0] != repo.url {
		repo.repository.DeleteRemote("origin")
		repo.repository.CreateRemote(&origin)
	} else if err != nil {
		repo.repository.CreateRemote(&origin)
	}

	return nil
}

// Fetch updates the tags in a given repo
func (repo *Repo) Fetch() error {
	err := repo.repository.Fetch(&git.FetchOptions{RemoteName: "origin", Tags: git.AllTags})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}
	return nil
}

// GetTags fetches semver compliant tags from a repo
func (repo *Repo) GetTags() (Tags, error) {
	return NewTagsFromRepo(repo.repository)
}

// SyncHash syncs a repo to a given commit hash
func (repo *Repo) SyncHash(hash string) error {
	// Checkout matching hash into worktree
	worktree, err := repo.repository.Worktree()
	if err != nil {
		return err
	}
	err = worktree.Checkout(&git.CheckoutOptions{Hash: plumbing.NewHash(hash)})
	if err != nil {
		return err
	}

	return nil
}

func (repo *Repo) Update(path, version string) error {
	// Fetch matching tags
	tags, err := NewTagsFromRepo(repo.repository)
	if err != nil {
		return err
	}

	// Fetch latest matching tag if available
	_, _, err = tags.GetLatest(version)
	if err != nil {
		return err
	}

	return nil
}
