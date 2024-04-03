package web

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
)

// fileSystem implements a simple filesystem abstraction for accessing files
// for static web serving and templates suitable for embedding or
// running live off a local machine during development.
//
// an http file server endpoint is also provided for the static content
type fileSystem struct {
	inDevelopment bool
	TplFS         fs.FS // filesystem for templates
	StaticFS      fs.FS // filesystem for static files
}

//go:embed templates
var templatesFS embed.FS

//go:embed static
var StaticFS embed.FS

// NewFileSystem returns a new fileSystem
func NewFileSystem(inDevelopment bool, tplDirPath, staticDirPath string) (*fileSystem, error) {

	var err error
	f := fileSystem{inDevelopment: inDevelopment}

	if !inDevelopment {
		// return embedded filesystem
		f.TplFS, err = fs.Sub(templatesFS, tplDirPath)
		if err != nil {
			return &f, err
		}
		f.StaticFS, err = fs.Sub(StaticFS, staticDirPath)
		if err != nil {
			return &f, err
		}

	} else {
		dirOK := func(d string) bool {
			if d == "" {
				return false
			}
			if _, err := os.Stat(d); os.IsNotExist(err) {
				return false
			}
			return true
		}
		if tplDirPath == "" || !dirOK(tplDirPath) {
			return &f, fmt.Errorf("tplDirPath %s could not be mounted", tplDirPath)
		}
		if staticDirPath == "" || !dirOK(staticDirPath) {
			return &f, fmt.Errorf("staticDirPath %s could not be mounted", staticDirPath)
		}
		f.TplFS = os.DirFS(tplDirPath)
		f.StaticFS = os.DirFS(staticDirPath)
	}

	return &f, nil
}
