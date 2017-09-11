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
	pathPart := patternString
	for {
		if fi, err := p.Fs.Stat(pathPart); err == nil {
			p.Path = strings.TrimRight(filesystem.CleanPath(p.Fs, pathPart), string(os.PathSeparator))
			p.isDirectory = fi == nil || fi.IsDir()
			break
		}
		parent := filepath.Dir(pathPart)
		if parent == pathPart {
			p.Path = "."
			p.Pattern = pathPart
			break
		}
		pathPart = parent
	}

	if pathPart != patternString {
		if pathPart == "." && len(patternString) == 1 {
			p.Path = "."
			p.Pattern = patternString
		} else {
			p.Pattern = patternString[len(p.Path)+1:]
		}
	}

}

func (p *BasePattern) IsDir() bool {
	return !p.IsFile()
}

func (p *BasePattern) IsFile() bool {
	return !p.isDirectory && p.Pattern == ""
}
