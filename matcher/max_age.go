package matcher

import (
	"time"
	"os"
)

type MaxAgeMatcher struct {
	MatcherInterface
	MaxAge time.Time
}

func NewMaxAgeMatcher(MaxAge time.Time) *MaxAgeMatcher {
	return &MaxAgeMatcher{
		MaxAge: MaxAge,
	}
}

func (f *MaxAgeMatcher) Matches(subject interface{}) bool {
	fi, err := os.Stat(subject.(string))

	if err != nil {
		return false
	}

	return f.MaxAge.Before(fi.ModTime())
}