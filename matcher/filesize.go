package matcher

import (
	"os"
	"github.com/spf13/afero"
)

type FileSizeMatcher struct {
	MatcherInterface
	Fs     afero.Fs
	fi      os.FileInfo
	minSize int64
	maxSize int64
}

func NewFileSizeMatcher(fi os.FileInfo, minSize, maxSize int64) *FileSizeMatcher {
	return &FileSizeMatcher{
		Fs:     afero.NewOsFs(),
		fi:      fi,
		minSize: minSize,
		maxSize: maxSize,
	}
}

func (fsMatcher *FileSizeMatcher) Matches(subject interface{}) bool {
	var err error
	if fsMatcher.fi == nil {
		fsMatcher.fi, err = fsMatcher.Fs.Stat(subject.(string))
	}

	if err != nil || fsMatcher.fi.IsDir(){
		return false
	}

	if fsMatcher.minSize < 0 {
		return fsMatcher.fi.Size() <= fsMatcher.maxSize
	}

	if fsMatcher.maxSize < 0 {
		return fsMatcher.fi.Size() >= fsMatcher.minSize
	}

	return fsMatcher.fi.Size() >= fsMatcher.minSize && fsMatcher.fi.Size() <= fsMatcher.maxSize
}