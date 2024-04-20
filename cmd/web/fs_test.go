package main

import (
	"io/fs"
	"testing"
)

// TestFSMount tests to see if the static and templates filesystems can
// be mounted and read
func TestFSMount(t *testing.T) {

	testCases := []struct {
		name          string
		inDevelopment bool
		tplDirPath    string
		staticDirPath string
	}{
		{"production", false, "templates", "static"}, // production embed fs doesn't need arguments
		{"development", true, "templates", "static"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			mount, err := NewFileSystem(tc.inDevelopment, tc.tplDirPath, tc.staticDirPath)
			if err != nil {
				t.Fatal(err)
			}

			d, err := fs.ReadDir(mount.StaticFS, ".")
			t.Log(d, err)

			_, err = fs.ReadFile(mount.StaticFS, "favicon.svg")
			if err != nil {
				t.Error(err)
			}

			d, err = fs.ReadDir(mount.TplFS, ".")
			t.Log(d, err)

			_, err = fs.ReadFile(mount.TplFS, "README")
			if err != nil {
				t.Error(err)
			}
		})
	}

}
