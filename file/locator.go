package file

import (
	"github.com/sandreas/graft/pattern"
	"github.com/sandreas/graft/designpattern/observer"
	"path/filepath"
	"os"
	"github.com/sandreas/graft/matcher"
	"strings"
)

const (
	LOCATOR_INCREASE_ITEMS = 1
	LOCATOR_INCREASE_MATCHES = 2
	LOCATOR_FINISH = 3
)

type Locator struct {
	designpattern.Observable
	Src         pattern.SourcePattern
	SourceFiles []string
}


func NewLocator(pattern pattern.SourcePattern) *Locator {
	return &Locator{
		Src: pattern,
	}
}


func (t *Locator) Find(matcher *matcher.CompositeMatcher) {
	t.SourceFiles = []string{}
	if t.Src.IsFile() {
		t.SourceFiles = append(t.SourceFiles, t.Src.Path)

		t.NotifyObservers(LOCATOR_INCREASE_MATCHES)
		t.NotifyObservers(LOCATOR_FINISH)
		return
	}

	filepath.Walk(t.Src.Path, func(innerPath string, info os.FileInfo, err error) error {
		if innerPath == "." || innerPath == ".." {
			return nil
		}
		normalizedInnerPath := strings.TrimRight(filepath.ToSlash(innerPath), "/")

		// skip direct path matches (data/* should not match data/ itself)
		if normalizedInnerPath == t.Src.Path && t.Src.Pattern != "" {
			return nil
		}

		if info.IsDir() {
			normalizedInnerPath += "/"
		}

		if matcher.Matches(normalizedInnerPath) {
			t.SourceFiles = append(t.SourceFiles, normalizedInnerPath)
			t.NotifyObservers(LOCATOR_INCREASE_MATCHES)
		} else {
			t.NotifyObservers(LOCATOR_INCREASE_ITEMS)
		}

		return nil
	})

	t.NotifyObservers(LOCATOR_FINISH)
}