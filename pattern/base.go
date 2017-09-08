package pattern

import (
	// "path/filepath"
	"strings"
	"github.com/spf13/afero"
	"path/filepath"
	"os"
	"runtime"
)

type BasePattern struct {
	Fs          afero.Fs
	Path        string
	Pattern     string
	isDirectory bool
}

func NewBasePattern(fs afero.Fs, patternString string) *BasePattern {
	basePattern := &BasePattern{
		Fs: fs,
	}
	basePattern.parse(patternString)
	return basePattern
}

func (p *BasePattern) parse(patternString string) {
	pathPart := patternString
	var slashIndex int
	for {
		if fi, err, path := p.AbsStat(pathPart); err == nil {
			p.Path = filepath.ToSlash(filepath.Clean(path))
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

func (p *BasePattern) AbsStat(path string) (os.FileInfo, error, string) {
	fi, err := p.Fs.Stat(path)
	//println(len(path))
	if runtime.GOOS != "windows" || len(path) < 250 || filepath.IsAbs(path) {
		return fi, err, path
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err, path
	}
	fi, err = p.Fs.Stat(absPath)
	// fi, err = os.Stat(absPath)
	return fi, err, absPath
}

func (p *BasePattern) IsDir() bool {
	return !p.IsFile()
}

func (p *BasePattern) IsFile() bool {
	return !p.isDirectory && p.Pattern == ""
}
