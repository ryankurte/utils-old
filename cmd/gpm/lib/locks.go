package gpm

import ()

// Locks is a list mapping modules to git commit hashes
type Locks map[string]string

const (
	// LockfileName is the default project lock file name
	LockfileName = ".lock.yml"
)
