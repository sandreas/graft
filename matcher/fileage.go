package matcher

import (
	"time"
	"github.com/spf13/afero"
	"os"
)

type FileAgeMatcher struct {
	MatcherInterface
	Fs     afero.Fs
	fi     os.FileInfo
	minAge time.Time
	maxAge time.Time
}

func NewFileAgeMatcher(fi os.FileInfo, minAge, maxAge time.Time) *FileAgeMatcher {
	return &FileAgeMatcher{
		Fs:     afero.NewOsFs(),
		fi:     fi,
		minAge: minAge,
		maxAge: maxAge,
	}
}

func (faMatcher *FileAgeMatcher) Matches(subject interface{}) bool {
	var err error
	if faMatcher.fi == nil {
		faMatcher.fi, err = faMatcher.Fs.Stat(subject.(string))
	}

	if err != nil {
		return false
	}

	if faMatcher.maxAge.IsZero() {
		return faMatcher.minAge.Before(faMatcher.fi.ModTime())
	}

	if faMatcher.minAge.IsZero() {
		return faMatcher.maxAge.After(faMatcher.fi.ModTime())
	}

	return faMatcher.maxAge.After(faMatcher.fi.ModTime()) && faMatcher.minAge.Before(faMatcher.fi.ModTime())

}
