package matcher

import (
	"time"
	"github.com/spf13/afero"
)

type FileAgeMatcher struct {
	MatcherInterface
	Fs     afero.Fs
	minAge time.Time
	maxAge time.Time
}

func NewFileAgeMatcher(minAge, maxAge time.Time) *FileAgeMatcher {
	return &FileAgeMatcher{
		Fs:     afero.NewOsFs(),
		minAge: minAge,
		maxAge: maxAge,
	}
}

func (faMatcher *FileAgeMatcher) Matches(subject interface{}) bool {
	fi, err := faMatcher.Fs.Stat(subject.(string))
	if err != nil {
		return false
	}

	if faMatcher.maxAge.IsZero() {
		return faMatcher.minAge.Before(fi.ModTime())
	}

	if faMatcher.minAge.IsZero() {
		return faMatcher.maxAge.After(fi.ModTime())
	}

	return faMatcher.maxAge.After(fi.ModTime()) && faMatcher.minAge.Before(fi.ModTime())

}
