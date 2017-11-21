package gpm

import (
	"fmt"
	"sort"

	"github.com/Masterminds/semver"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// Tags type is a map of semver compatible git annotated tags and hashes
type Tags map[string]string

// NewTagsFromRepo creates a new Tags object from the provided git repository
func NewTagsFromRepo(repo *git.Repository) (Tags, error) {
	tags := make(Tags)
	tagIter, _ := repo.TagObjects()
	err := tagIter.ForEach(func(t *object.Tag) error {
		commit, err := t.Commit()
		if err != nil {
			return nil
		}
		tags[t.Name] = commit.Hash.String()
		return nil
	})
	return tags, err
}

// Filter filters tags by a given semver filter
func (tags *Tags) Filter(filter string) (Tags, error) {
	filtered := make(Tags)
	var err error

	// Load range filter
	var constraints *semver.Constraints
	if filter != "" {
		constraints, err = semver.NewConstraint(filter)
		if err != nil {
			return nil, err
		}
	}

	// Iterate and filter tags
	for t, h := range *tags {
		// Attempt to parse value
		v, err := semver.NewVersion(t)
		if err != nil {
			fmt.Printf("Tag: %s parsing failed (%s)\n", t, err)
			continue
		}
		// Apply filter if available
		if constraints != nil && !constraints.Check(v) {
			continue
		}
		// Add to list if within range
		filtered[t] = h
	}

	return filtered, nil
}

// Sort creates a sorted array based of tags
func (tags *Tags) Sort() ([]string, error) {
	// Create array for sorting
	tagArr := make([]string, len(*tags))
	i := 0
	for t := range *tags {
		tagArr[i] = t
		i++
	}

	// Sort into ascending tag order
	sort.SliceStable(tagArr, func(i, j int) bool {
		vi, _ := semver.NewVersion(tagArr[i])
		vj, _ := semver.NewVersion(tagArr[j])
		return vi.LessThan(vj)
	})

	return tagArr, nil
}

// GetLatest filters tags by the specified range then finds the latest matching tag
func (tags *Tags) GetLatest(filter string) (string, string, error) {
	filtered, err := tags.Filter(filter)
	if err != nil {
		return "", "", err
	}
	ordered, err := filtered.Sort()
	if err != nil {
		return "", "", err
	}
	if len(ordered) == 0 {
		return "", "", nil
	}

	latest := ordered[len(ordered)-1]
	return latest, filtered[latest], nil
}
