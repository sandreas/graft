package pattern

import (
	// "path/filepath"
	"strings"
	"github.com/spf13/afero"
	"path/filepath"
	"os"
	"runtime"
	"github.com/sandreas/graft/filesystem"
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
	fi, err := filesystem.Stat(p.Fs, path)
	//println(len(path))

	if !p.isLocalFs || runtime.GOOS != "windows" || len(path) < 250 || filepath.IsAbs(path) {
		return fi, err, path
	}

	var absPath string
	if absPath, err = filepath.Abs(path); err != nil {
		return nil, err, path
	}
	fi, err = filesystem.Stat(p.Fs, absPath)
	// fi, err = os.Stat(absPath)
	return fi, err, filepath.ToSlash(absPath)
}

func (p *BasePattern) IsDir() bool {
	return !p.IsFile()
}

func (p *BasePattern) IsFile() bool {
	return !p.isDirectory && p.Pattern == ""
}
