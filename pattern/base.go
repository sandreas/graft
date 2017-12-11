package pattern

import (
	"path/filepath"

	"os"

	"github.com/sandreas/graft/filesystem"
	"github.com/spf13/afero"
	"strings"
)

type BasePattern struct {
	Fs          afero.Fs
	Path        string
	Pattern     string
	isDirectory bool
	isLocalFs   bool
}

func NewBasePattern(fs afero.Fs, patternString string) *BasePattern {
	_, isOsFs := fs.(*afero.OsFs)
	_, isMemMapFs := fs.(*afero.MemMapFs)

	basePattern := &BasePattern{
		Fs:        fs,
		isLocalFs: isOsFs || isMemMapFs,
	}

	basePattern.parse(patternString)
	return basePattern
}

func (p *BasePattern) parse(patternString string) {

	sep := string(os.PathSeparator)
	patternString = strings.TrimPrefix(patternString, "./")
	pathPart := patternString

	for {
		if fi, err := p.Fs.Stat(pathPart); err == nil {
			p.Path = filesystem.CleanPath(p.Fs, pathPart)
			p.isDirectory = fi == nil || fi.IsDir()
			if p.isDirectory && !os.IsPathSeparator(p.Path[len(p.Path)-1]) {
				p.Path += sep
			}
			break
		}

		parent := filepath.Dir(pathPart)
		if parent == pathPart || parent == "." {
			p.Path = ""
			p.Pattern = pathPart
			break
		}
		pathPart = parent
	}

	if len(patternString) > 0 && pathPart != patternString {
		p.Pattern = patternString[len(p.Path):]
	}
	if p.Path == "" || p.Path == "." {
		p.Path = "." + sep
		p.isDirectory = true
	}
}

func (p *BasePattern) IsDir() bool {
	return !p.IsFile()
}

func (p *BasePattern) IsFile() bool {
	return !p.isDirectory && p.Pattern == ""
}
