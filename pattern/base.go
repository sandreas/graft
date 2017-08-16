package pattern

import (
	"path/filepath"
	"strings"
	"github.com/spf13/afero"
)

type BasePattern struct {
	fs          afero.Fs
	Path        string
	Pattern     string
	isDirectory bool
}

func NewBasePattern(fs afero.Fs, patternString string) *BasePattern {
	basePattern := &BasePattern{
		fs: fs,
	}
	basePattern.parse(patternString)
	return basePattern
}

func (p *BasePattern) parse(patternString string) {
	if fi, err := p.fs.Stat(patternString); err != nil {
		pathPart := patternString
		path := ""
		for {
			slashIndex := strings.IndexAny(pathPart, "\\/")
			if slashIndex == -1 {
				break
			}
			pathCandidate := path + pathPart[0:slashIndex+1]
			fi, err := p.fs.Stat(pathCandidate)
			if err != nil {
				break
			}

			path = pathCandidate
			pathPart = pathPart[slashIndex+1:]
			p.isDirectory = fi.IsDir()
		}
		p.Path = filepath.ToSlash(filepath.Clean(path))
		p.Pattern = strings.TrimPrefix(patternString, path)
	} else {
		p.Path = filepath.ToSlash(filepath.Clean(patternString))
		p.Pattern = ""
		p.isDirectory = fi.IsDir()
	}
}

func (p *BasePattern) IsDir() bool {
	return !p.IsFile()
}

func (p *BasePattern) IsFile() bool {
	return !p.isDirectory && p.Pattern == ""
}
