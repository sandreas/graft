package file

import (
	"fmt"
	"path/filepath"
	"os"
	"strings"
)

//:host
//:device
//:version
//:directory
//:name
//:type

const (
	typeNotExisting = 0
	//typeDir         = 1
	//typeFile        = 2
)

type Path struct {
	fmt.Stringer
	volume   string
	path     string
	kind     int
	absolute bool
	stat     os.FileInfo
}

func NewPath(path string) *Path {
	p := &Path{
		volume:   filepath.VolumeName(path),
		path:     filepath.Clean(path),
		kind:     typeNotExisting,
		absolute: false,
	}
	p.normalize()
	return p
}

func (path *Path) normalize() {
	dirsep := string(os.PathSeparator)
	if path.volume != "" || strings.HasPrefix(path.path, dirsep) {
		path.absolute = true
		path.volume = strings.TrimRight(path.volume, "\\/")
	} else {
		path.path = strings.TrimPrefix(path.path, "."+dirsep)
	}
}

func (path *Path) build() string {
	if path.absolute {
		return path.volume + path.path
	}
	return path.path
}

func (path *Path) Stat() (os.FileInfo, error) {
	var err error
	path.stat, err = os.Stat(path.build())
	if os.IsNotExist(err) {
		path.kind = typeNotExisting
		path.stat = nil
	}
	return path.stat, err
}

func (path *Path) String() string {
	return path.build()
}
