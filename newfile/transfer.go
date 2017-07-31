package newfile

import (
	"github.com/sandreas/graft/newpattern"
	"github.com/sandreas/graft/newdesignpattern/observer"
	"path/filepath"
	"os"
	"github.com/sandreas/graft/newmatcher"
)

const (
	OBSERVER_INCREASE_ITEMS = 1
	OBSERVER_INCREASE_MATCHES = 2
	OBSERVER_FINISH = 3
)

type Transfer struct {
	newdesignpattern.Observable
	src newpattern.SourcePattern
	SourceFiles []string
}


func NewTransfer(pattern newpattern.SourcePattern) *Transfer {
	return &Transfer{
		src: pattern,
	}
}


func (t *Transfer) Find(matcher *newmatcher.CompositeMatcher) {
	t.SourceFiles = []string{}
	if t.src.IsFile() {
		t.SourceFiles = append(t.SourceFiles, t.src.Path)

		t.NotifyObservers(OBSERVER_INCREASE_MATCHES)
		t.NotifyObservers(OBSERVER_FINISH)
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
			t.NotifyObservers(OBSERVER_INCREASE_MATCHES)
		} else {
			t.NotifyObservers(OBSERVER_INCREASE_ITEMS)
		}

		return nil
	})

	t.NotifyObservers(OBSERVER_FINISH)
}


func (t *Transfer) CopyTo(dst *newpattern.BasePattern) {

	if dst.IsFile() {

	}

}

//func (t *Transfer) moveTo(dst string) {
//
//}
//
//func (t *Transfer) remove(dst string) {
//
//}