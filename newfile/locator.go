package newfile

import (
	"github.com/sandreas/graft/newpattern"
	"github.com/sandreas/graft/newdesignpattern/observer"
	"path/filepath"
	"os"
	"github.com/sandreas/graft/newmatcher"
)

const (
	LOCATOR_INCREASE_ITEMS = 1
	LOCATOR_INCREASE_MATCHES = 2
	LOCATOR_FINISH = 3
)

type Locator struct {
	newdesignpattern.Observable
	src newpattern.SourcePattern
	SourceFiles []string
}


func NewLocator(pattern newpattern.SourcePattern) *Locator {
	return &Locator{
		src: pattern,
	}
}


func (t *Locator) Find(matcher *newmatcher.CompositeMatcher) {
	t.SourceFiles = []string{}
	if t.src.IsFile() {
		t.SourceFiles = append(t.SourceFiles, t.src.Path)

		t.NotifyObservers(LOCATOR_INCREASE_MATCHES)
		t.NotifyObservers(LOCATOR_FINISH)
		return
	}

	filepath.Walk(t.src.Path, func(innerPath string, info os.FileInfo, err error) error {
		if innerPath == "." || innerPath == ".." {
			return nil
		}


		normalizedInnerPath := filepath.ToSlash(innerPath)
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