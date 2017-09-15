package file

import (
	"log"
	"os"

	"path/filepath"

	"github.com/sandreas/graft/designpattern/observer"
	"github.com/sandreas/graft/filesystem"
	"github.com/sandreas/graft/matcher"
	"github.com/sandreas/graft/pattern"
)

const (
	LocatorIncreaseItems   = 1
	LocatorIncreaseMatches = 2
	LocatorFinish          = 3
	LocatorIncreaseErrors  = 4
)

type Locator struct {
	designpattern.Observable
	Src         pattern.SourcePattern
	SourceFiles []string
}

func NewLocator(pattern *pattern.SourcePattern) *Locator {
	return &Locator{
		Src: *pattern,
	}
}

func (t *Locator) Find(matcher *matcher.CompositeMatcher) {
	t.SourceFiles = []string{}
	if t.Src.IsFile() {
		t.SourceFiles = append(t.SourceFiles, t.Src.Path)

		t.NotifyObservers(LocatorIncreaseMatches)
		t.NotifyObservers(LocatorFinish)
		return
	}

	walkPath := filesystem.CleanPath(t.Src.Fs, t.Src.Path)
	walkPathSeparator := string(os.PathSeparator)
	if t.Src.Fs.Name() == filesystem.NameSftpfs {
		walkPathSeparator = "/"
	}
	ferr := filesystem.Walk(t.Src.Fs, walkPath, func(innerPath string, info os.FileInfo, err error) error {
		if innerPath == "." || innerPath == ".." {
			return nil
		}

		if err != nil {
			t.NotifyObservers(LocatorIncreaseErrors)
			log.Printf("WalkError: %s, Details: %v", err.Error(), err)
			return nil
		}

		normalizedInnerPath := filesystem.CleanPath(t.Src.Fs, innerPath)

		// skip direct path matches (data/* should not match data/ itself)
		if normalizedInnerPath == walkPath && t.Src.Pattern != "" {
			return nil
		}

		if info.IsDir() {
			normalizedInnerPath += walkPathSeparator
		}

		if matcher.Matches(filepath.ToSlash(normalizedInnerPath)) {
			t.SourceFiles = append(t.SourceFiles, normalizedInnerPath)
			t.NotifyObservers(LocatorIncreaseMatches)
		} else {
			t.NotifyObservers(LocatorIncreaseItems)
		}

		return nil
	})

	if ferr != nil {
		log.Printf("Error walking files: %s\n", ferr.Error())
	}

	t.NotifyObservers(LocatorFinish)
}
