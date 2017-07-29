package newmatcher

import (
	"time"
	"os"
)

type MinAgeMatcher struct {
	MatcherInterface
	minAge time.Time
}

func NewMinAgeMatcher(minAge time.Time) *MinAgeMatcher {
	return &MinAgeMatcher{
		minAge: minAge,
	}
}

func (f *MinAgeMatcher) Matches(subject interface{}) bool {
	fi, err := os.Stat(subject.(string))

	if err != nil {
		return false
	}

	return f.minAge.After(fi.ModTime())
}