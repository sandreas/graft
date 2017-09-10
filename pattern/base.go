package pattern

import (
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
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
	var slashIndex int
	for {
		if fi, err := p.Fs.Stat(pathPart); err == nil {
			p.Path = strings.TrimRight(filepath.Clean(pathPart), "\\/")
			p.isDirectory = fi == nil || fi.IsDir()
			break
		}
		slashIndex = strings.LastIndexAny(pathPart, "\\/")
		if slashIndex == -1 {
			p.Path = "."
			p.Pattern = pathPart
			break
		}
		pathPart = pathPart[0:slashIndex]
	}

	if pathPart != patternString {
		p.Pattern = patternString[len(pathPart)+1:]
	}

}

func (p *BasePattern) IsDir() bool {
	return !p.IsFile()
}

func (p *BasePattern) IsFile() bool {
	return !p.isDirectory && p.Pattern == ""
}
