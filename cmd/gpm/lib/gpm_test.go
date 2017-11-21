package gpm

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	// Setup temporary test directory
	testDir, err := ioutil.TempDir("", "gpm-test")
	if err != nil {
		t.FailNow()
	}
	//	defer os.Remove(testDir)

	configFile := fmt.Sprintf("%s/%s", testDir, ProjectConfigName)
	t.Logf("GPM test dir: %s config: %s", testDir, configFile)

	o := CommonOptions{
		BasePath:  testDir,
		Verbose:   true,
		NoCleanup: true,
	}

	gpm := GPM{options: &o}

	po := InitOptions{
		Name:       "Test Project",
		Repository: "https://github.com/ryankurte/test-repo",
		License:    "MIT",
	}

	versionZeroOneZero := "v0.1.0"
	hashZeroOneZero := "5a83540f87665410bb6440b97f6c2f71a27471bd"

	versionZeroTwoZero := "v0.2.0"
	hashZeroTwoZero := "2075d8e0a26e6104b212428186bec59edad15e42"

	testRepo := "https://github.com/ryankurte/test-repo"

	t.Run("Initialise a project", func(t *testing.T) {
		// Initialise project with the provided options
		err := gpm.Init(&po)
		assert.Nil(t, err)

		// Check config is correct
		pc, err := gpm.loadProjectConfig()
		assert.Nil(t, err)
		assert.EqualValues(t, po.Name, pc.Name)
		assert.EqualValues(t, 0, len(pc.Dependencies))
	})

	t.Run("Re-init fails", func(t *testing.T) {
		err := gpm.Init(&po)
		assert.NotNil(t, err)
	})

	t.Run("Add a dependency (no version)", func(t *testing.T) {
		d, _ := NewDependency("test1", testRepo, "")

		err := gpm.Add(&AddOptions{d.Path, d.URL, d.Version})
		assert.Nil(t, err)

		// Check dependency got added to config
		pc, err := gpm.loadProjectConfig()
		assert.Nil(t, err)
		assert.EqualValues(t, po.Name, pc.Name)

		d.Version = versionZeroTwoZero
		if d1, ok := pc.Dependencies.Find(d.Path); ok {
			assert.EqualValues(t, d, d1)
		} else {
			t.FailNow()
		}

		locks, err := gpm.loadLockfile()
		assert.Nil(t, err)
		assert.EqualValues(t, hashZeroTwoZero, locks["test1"])
	})

	t.Run("Add a dependency (with version)", func(t *testing.T) {
		d, _ := NewDependency("test2", testRepo, versionZeroOneZero)

		err := gpm.Add(&AddOptions{d.Path, d.URL, d.Version})
		assert.Nil(t, err)

		// Check dependency got added to config
		pc, err := gpm.loadProjectConfig()
		assert.Nil(t, err)
		assert.EqualValues(t, po.Name, pc.Name)
		assert.EqualValues(t, 2, len(pc.Dependencies))

		if d1, ok := pc.Dependencies.Find(d.Path); ok {
			assert.EqualValues(t, d, d1)
		} else {
			t.FailNow()
		}

		locks, err := gpm.loadLockfile()
		assert.Nil(t, err)
		assert.EqualValues(t, hashZeroOneZero, locks["test2"])
	})

	t.Run("Updates dependencies", func(t *testing.T) {
		pc, err := gpm.loadProjectConfig()
		d, _ := NewDependency("test2", testRepo, versionZeroTwoZero)
		pc.Dependencies.Set("test2", *d)
		assert.Nil(t, gpm.writeProjectConfig(&pc))

		err = gpm.Update(&UpdateOptions{})
		assert.Nil(t, err)

		locks, err := gpm.loadLockfile()
		assert.Nil(t, err)
		assert.EqualValues(t, hashZeroTwoZero, locks["test2"])
	})

	t.Run("Syncs removed dependencies", func(t *testing.T) {
		p, _ := o.fullPath("test2")
		os.RemoveAll(p)

		err := gpm.Sync(&SyncOptions{})
		assert.Nil(t, err)

		if _, err := os.Stat(p); os.IsNotExist(err) {
			t.Errorf("Path '%s' does not exist", p)
		}
	})

	t.Run("Removes dependencies", func(t *testing.T) {
		p, _ := o.fullPath("test2")
		os.RemoveAll(p)

		err := gpm.Remove(&RemoveOptions{"test2"})
		assert.Nil(t, err)

		if _, err := os.Stat(p); !os.IsNotExist(err) {
			t.Errorf("Path '%s' does not exist", p)
		}

		pc, err := gpm.loadProjectConfig()
		assert.Nil(t, err)
		assert.EqualValues(t, 1, len(pc.Dependencies))
	})

}
