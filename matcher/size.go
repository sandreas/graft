package matcher

import (
	"os"
)

type FileSizeMatcher struct {
	MatcherInterface
	fi      os.FileInfo
	minSize int64
	maxSize int64
}

func NewFileSizeMatcher(fi os.FileInfo, minSize, maxSize int64) *FileSizeMatcher {
	return &FileSizeMatcher{
		fi:      fi,
		minSize: minSize,
		maxSize: maxSize,
	}
}

func (f *FileSizeMatcher) Matches(subject interface{}) bool {
	var err error
	if f.fi == nil {
		f.fi, err = os.Stat(subject.(string))
	}

	if err != nil || f.fi.IsDir(){
		return false
	}

	if f.minSize < 0 {
		return f.fi.Size() <= f.maxSize
	}

	if f.maxSize < 0 {
		return f.fi.Size() >= f.minSize
	}

	return f.fi.Size() >= f.minSize && f.fi.Size() <= f.maxSize
}