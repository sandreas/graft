package newpattern

import (
	"os"
	"path/filepath"
	"strings"
)

type BasePattern struct {
	Path string
	Pattern string
	isDirectory bool
}

func NewBasePattern(patternString string) *BasePattern {
	basePattern := &BasePattern{}
	basePattern.parse(patternString)
	return basePattern
}

func (p *BasePattern) parse(patternString string) {
	path := patternString
	for {
		if fi, err := os.Stat(path); err == nil {
			p.Path = filepath.ToSlash(path)

			p.isDirectory = fi.IsDir()


			startIndex := len(p.Path)+1
			if path == "." {
				p.Pattern = strings.TrimPrefix(patternString, ".")
				p.Pattern = strings.TrimPrefix(p.Pattern, "/")
			} else if len(patternString) > startIndex {
				p.Pattern = patternString[startIndex:]
			}
			break
		}
		path = filepath.Dir(path)
	}

	p.Path = strings.TrimSuffix(p.Path, "/")
}

func (p *BasePattern) IsDir() bool {
	return !p.IsFile()
}

func (p *BasePattern) IsFile() bool {
	return !p.isDirectory && p.Pattern == ""
}


