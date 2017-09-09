// +build windows

package filesystem

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/afero"
)

type OsFs struct {
	afero.Fs
}

func NewOsFs() afero.Fs {
	return &OsFs{}
}

func (OsFs) Name() string { return "OsFs" }

func (OsFs) Create(name string) (afero.File, error) {
	f, e := os.Create(normalizePath(name))
	if f == nil {
		// while this looks strange, we need to return a bare nil (of type nil) not
		// a nil value of type *os.File or nil won't be nil
		return nil, e
	}
	return f, e
}

func (OsFs) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(normalizePath(name), perm)
}

func (OsFs) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(normalizePath(path), perm)
}

func (OsFs) Open(name string) (afero.File, error) {
	f, e := os.Open(normalizePath(name))
	if f == nil {
		// while this looks strange, we need to return a bare nil (of type nil) not
		// a nil value of type *os.File or nil won't be nil
		return nil, e
	}
	return f, e
}

func (OsFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	f, e := os.OpenFile(normalizePath(name), flag, perm)
	if f == nil {
		// while this looks strange, we need to return a bare nil (of type nil) not
		// a nil value of type *os.File or nil won't be nil
		return nil, e
	}
	return f, e
}

func (OsFs) Remove(name string) error {
	return os.Remove(normalizePath(name))
}

func (OsFs) RemoveAll(path string) error {
	return os.RemoveAll(normalizePath(path))
}

func (OsFs) Rename(oldname, newname string) error {
	return os.Rename(normalizePath(oldname), normalizePath(newname))
}

func (OsFs) Stat(name string) (os.FileInfo, error) {
	return os.Stat(normalizePath(name))
}

func (OsFs) Chmod(name string, mode os.FileMode) error {
	return os.Chmod(normalizePath(name), mode)
}

func (OsFs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return os.Chtimes(normalizePath(name), atime, mtime)
}

func normalizePath(path string) string {
	if filepath.IsAbs(path) || len(path) < 227 {
		return path
	}

	if absPath, err := filepath.Abs(path); err == nil {
		return absPath
	} else {
		log.Printf("Could not determine absolute path for %s - this can lead to misbehaviour\n", path)
		return path
	}
}

func (fs *OsFs) Close() {
	// osfs does not need to be closed
}
