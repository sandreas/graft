package filesystem

import "github.com/spf13/afero"

type OsFs struct {
	*afero.OsFs
}

func NewOsFs() (afero.Fs, error) {
	return &OsFs{}, nil
}

func (fs *OsFs) Close() {
	// osfs does not need to be closed
}
