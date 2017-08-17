package matcher

import (
	"os"
	"time"
)

type FileAgeMatcher struct {
	MatcherInterface
	fi     os.FileInfo
	minAge time.Time
	maxAge time.Time
}

func NewFileAgeMatcher(fi os.FileInfo, minAge, maxAge time.Time) *FileAgeMatcher {
	return &FileAgeMatcher{
		fi:     fi,
		minAge: minAge,
		maxAge: maxAge,
	}
}

func (f *FileAgeMatcher) Matches(subject interface{}) bool {
	var err error
	if f.fi == nil {
		f.fi, err = os.Stat(subject.(string))
	}

	if err != nil {
		return false
	}

	if f.maxAge.IsZero() {
		return  f.minAge.Before(f.fi.ModTime())
	}

	if f.minAge.IsZero() {
		return  f.maxAge.After(f.fi.ModTime())
	}

	return f.maxAge.After(f.fi.ModTime()) && f.minAge.Before(f.fi.ModTime())

}