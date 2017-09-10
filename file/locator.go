package file

import (
	"log"
	"os"

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

	filesystem.Walk(t.Src.Fs, t.Src.Path, func(innerPath string, info os.FileInfo, err error) error {
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
		if normalizedInnerPath == t.Src.Path && t.Src.Pattern != "" {
			return nil
		}

		if info.IsDir() {
			normalizedInnerPath += string(os.PathSeparator)
		}

		if matcher.Matches(normalizedInnerPath) {
			t.SourceFiles = append(t.SourceFiles, normalizedInnerPath)
			t.NotifyObservers(LocatorIncreaseMatches)
		} else {
			t.NotifyObservers(LocatorIncreaseItems)
		}

		return nil
	})

	t.NotifyObservers(LocatorFinish)
}
