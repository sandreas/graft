// +build windows

package filesystem

import (

	"os"
	"path/filepath"
	"time"

	"github.com/spf13/afero"
	"strings"
)

type OsFs struct {
	afero.Fs
}

func NewOsFs() afero.Fs {
	return &OsFs{}
}

func (OsFs) Name() string { return "OsFs" }

func (OsFs) Create(name string) (afero.File, error) {
	f, e := os.Create(makeAbsolute(name))
	if f == nil {
		// while this looks strange, we need to return a bare nil (of type nil) not
		// a nil value of type *os.File or nil won't be nil
		return nil, e
	}
	return f, e
}

func (OsFs) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(makeAbsolute(name), perm)
}

func (OsFs) MkdirAll(path string, perm os.FileMode) error {
	abs := makeAbsolute(path)

	// mkdirall fails on absolute paths, that do not exist, e.g. abs := "\\\\?\\D:\\test"
	// see https://github.com/golang/go/issues/22230
	if err:=os.Mkdir(abs, perm); err == nil {
		return err
	}
	return os.MkdirAll(abs, perm)
}

func (OsFs) Open(name string) (afero.File, error) {
	f, e := os.Open(makeAbsolute(name))
	if f == nil {
		// while this looks strange, we need to return a bare nil (of type nil) not
		// a nil value of type *os.File or nil won't be nil
		return nil, e
	}
	return f, e
}

func (OsFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	f, e := os.OpenFile(makeAbsolute(name), flag, perm)
	if f == nil {
		// while this looks strange, we need to return a bare nil (of type nil) not
		// a nil value of type *os.File or nil won't be nil
		return nil, e
	}
	return f, e
}

func (OsFs) Remove(name string) error {
	return os.Remove(makeAbsolute(name))
}

func (OsFs) RemoveAll(path string) error {
	return os.RemoveAll(makeAbsolute(path))
}

func (OsFs) Rename(oldname, newname string) error {
	return os.Rename(makeAbsolute(oldname), makeAbsolute(newname))
}

func (OsFs) Stat(name string) (os.FileInfo, error) {
	return os.Stat(makeAbsolute(name))
}

func (OsFs) Chmod(name string, mode os.FileMode) error {
	return os.Chmod(makeAbsolute(name), mode)
}

func (OsFs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return os.Chtimes(makeAbsolute(name), atime, mtime)
}


// windows cannot handle long relative paths, so relative are converted to absolute paths by default
func makeAbsolute(name string) string {
	absolutePath, err := filepath.Abs(name)
	if err == nil {
		if strings.HasPrefix(absolutePath, `\\?\UNC\`) || strings.HasPrefix(absolutePath, `\\?\`) {
			return absolutePath
		}

		if strings.HasPrefix(absolutePath, `\\`) {
			return strings.Replace(absolutePath, `\\`, `\\?\UNC\`, 1)
		}

		return `\\?\` + absolutePath
	}
	return name
}

func (fs *OsFs) Close() {
	// osfs does not need to be closed
}
