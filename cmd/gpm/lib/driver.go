package gpm

// Driver defines the interface to be implemented by gpm drivers
type Driver interface {
	Add(Dependency) error
	Sync(Dependency) error
	Update(Dependency) error
	Remove(Dependency) error
}
