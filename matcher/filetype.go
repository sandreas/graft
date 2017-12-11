package matcher

import (
	"github.com/sandreas/graft/filesystem"
	"github.com/spf13/afero"
)

const (
	TypeFile      = "f"
	TypeDirectory = "d"
)

type FileTypeMatcher struct {
	MatcherInterface
	Fs            afero.Fs
	matchingTypes []string
}

func NewFileTypeMatcher(matchingTypes ...string) *FileTypeMatcher {
	return &FileTypeMatcher{
		Fs:            filesystem.NewOsFs(),
		matchingTypes: matchingTypes,
	}
}

func (faMatcher *FileTypeMatcher) Matches(subject interface{}) bool {
	fi, err := faMatcher.Fs.Stat(subject.(string))
	if err != nil {
		return false
	}

	if fi.IsDir() {
		return contains(faMatcher.matchingTypes, TypeDirectory)
	}

	return contains(faMatcher.matchingTypes, TypeFile)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
