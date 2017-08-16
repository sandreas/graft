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
	typeDir  = 1
	typeFile = 2
	separator = string(os.PathSeparator)
)

type Path struct {
	fmt.Stringer
	volume   string
	path     string
	name     string
	kind     int
	absolute bool
}

func NewPath(path string) *Path {
	vol := filepath.VolumeName(path)
	path = strings.TrimPrefix(path, vol)
	kind := getPathType(path)
	path = filepath.Clean(path)
	name := ""

	if kind == typeFile {
		dir := filepath.Dir(path) + separator
		name = strings.TrimPrefix(path, dir)
		path = strings.TrimSuffix(dir, separator)
	}
	if path == "." {
		path = ""
	}

	p := &Path{
		volume:   filepath.FromSlash(vol),
		path:     path,
		name:	name,
		kind:     kind,
		absolute: false,
	}

	p.normalize()
	return p
}

func getPathType(path string) int {
	if strings.HasSuffix(path, "/") || strings.HasSuffix(path, "\\") {
		return typeDir
	}
	return typeFile
}

func (path *Path) normalize() {
	if path.volume != "" || strings.HasPrefix(path.path, separator) {
		path.absolute = true
		path.volume = strings.TrimRight(path.volume, "\\/")
	} else {
		path.path = strings.TrimPrefix(path.path, "."+separator)
	}
}

func (path *Path) build() string {
	if path.absolute {
		return path.volume + path.path + separator + path.name
	}
	return path.path + separator + path.name
}

func (path *Path) String() string {
	return path.build()
}

func (path *Path) IsDir() bool {
	return path.kind == typeDir
}

func (path *Path) IsFile() bool {
	return path.kind == typeFile
}

func (path *Path) IsAbs() bool {
	return path.absolute
}
