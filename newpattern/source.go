package newpattern

import (
	"os"
	"path/filepath"
	"strings"
)

type SourcePattern struct {
	Path string
	Pattern string
	IsDirectory bool
}

func NewSourcePattern(patternString string) *SourcePattern {
	sourcePattern := &SourcePattern{}
	sourcePattern.Parse(patternString)
	return sourcePattern
}
func (p *SourcePattern) Parse(patternString string) {
	p.Path = ""
	p.Pattern = ""
	path := patternString
	for {
		if fi, err := os.Stat(path); err == nil {
			p.Path = filepath.ToSlash(path)

			p.IsDirectory = fi.IsDir()

			startIndex := len(p.Path)+1
			if path == "." {
				p.Pattern = strings.TrimPrefix(patternString, "./")
			} else if len(patternString) > startIndex {
				p.Pattern = patternString[startIndex:]
			}
			break
		}
		path = filepath.Dir(path)
	}
}