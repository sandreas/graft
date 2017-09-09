// +build !windows

package filesystem

import (
	"path/filepath"

	"github.com/spf13/afero"
)

func Walk(fs afero.Fs, root string, walkFn filepath.WalkFunc) error {
	return afero.Walk(fs, root, walkFn)
}
