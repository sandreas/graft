package matcher

import (
	"os"
)

type SizeMatcher struct {
	MatcherInterface
	MinSize int64
	MaxSize	int64
}

func NewSizeMatcher(minSize, maxSize int64) *SizeMatcher {
	return &SizeMatcher{
		MinSize: minSize,
		MaxSize: maxSize,
	}
}

func (f *SizeMatcher) Matches(subject interface{}) bool {
	fi, err := os.Stat(subject.(string))

	if err != nil || fi.IsDir(){
		return false
	}

	if f.MinSize < 0 {
		return fi.Size() <= f.MaxSize
	}

	if f.MaxSize < 0 {
		return fi.Size() >= f.MinSize
	}

	return fi.Size() >= f.MinSize && fi.Size() <= f.MaxSize
}