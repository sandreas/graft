package filesystem

import (
	"path/filepath"

	"github.com/spf13/afero"
)

func CleanPath(fs afero.Fs, path string) string {
	if fs.Name() == NameSftpfs {
		return filepath.ToSlash(filepath.Clean(path))
	}

	return filepath.Clean(path)
}
