package matcher

import (
	"github.com/spf13/afero"
)

type FileSizeMatcher struct {
	MatcherInterface
	Fs      afero.Fs
	minSize int64
	maxSize int64
}

func NewFileSizeMatcher(minSize, maxSize int64) *FileSizeMatcher {
	return &FileSizeMatcher{
		Fs:      afero.NewOsFs(),
		minSize: minSize,
		maxSize: maxSize,
	}
}

func (fsMatcher *FileSizeMatcher) Matches(subject interface{}) bool {
	fi, err := fsMatcher.Fs.Stat(subject.(string))

	if err != nil || fi.IsDir() {
		return false
	}

	if fsMatcher.minSize < 0 {
		return fi.Size() <= fsMatcher.maxSize
	}

	if fsMatcher.maxSize < 0 {
		return fi.Size() >= fsMatcher.minSize
	}

	return fi.Size() >= fsMatcher.minSize && fi.Size() <= fsMatcher.maxSize
}
