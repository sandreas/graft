package pattern

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"path/filepath"
)

func TestBase(t *testing.T) {
	expect := assert.New(t)

	sourcePattern := NewBasePattern("../data/fixtures/global/*")
	expect.Equal("../data/fixtures/global", sourcePattern.Path)
	expect.Equal("*", sourcePattern.Pattern)
	expect.True(sourcePattern.IsDir())
	expect.False(sourcePattern.IsFile())

	sourcePattern = NewBasePattern("../data/fixtures/non-existing/*.*")
	expect.Equal("../data/fixtures", sourcePattern.Path)
	expect.Equal("non-existing/*.*", sourcePattern.Pattern)
	expect.True(sourcePattern.IsDir())
	expect.False(sourcePattern.IsFile())

	sourcePattern = NewBasePattern("../data/fixtures/global/file.txt")
	expect.Equal("../data/fixtures/global/file.txt", sourcePattern.Path)
	expect.Equal("", sourcePattern.Pattern)
	expect.False(sourcePattern.IsDir())
	expect.True(sourcePattern.IsFile())

	sourcePattern = NewBasePattern("../data/fixtures/global/")
	expect.Equal("../data/fixtures/global", sourcePattern.Path)
	expect.Equal("", sourcePattern.Pattern)
	expect.True(sourcePattern.IsDir())
	expect.False(sourcePattern.IsFile())

	sourcePattern = NewBasePattern("*")
	expect.Equal(".", sourcePattern.Path)
	expect.Equal("*", sourcePattern.Pattern)
	expect.True(sourcePattern.IsDir())
	expect.False(sourcePattern.IsFile())

	sourcePattern = NewBasePattern("./*")
	expect.Equal(".", sourcePattern.Path)
	expect.Equal("*", sourcePattern.Pattern)
	expect.True(sourcePattern.IsDir())
	expect.False(sourcePattern.IsFile())

	sourcePattern = NewBasePattern("../data/fixtures/global/(.*)")
	expect.Equal("../data/fixtures/global", sourcePattern.Path)
	expect.Equal("(.*)", sourcePattern.Pattern)
	expect.True(sourcePattern.IsDir())
	expect.False(sourcePattern.IsFile())

	sourcePattern = NewBasePattern(".")
	expect.Equal(".", sourcePattern.Path)
	expect.Equal("", sourcePattern.Pattern)
	expect.True(sourcePattern.IsDir())
	expect.False(sourcePattern.IsFile())

	sourcePattern = NewBasePattern("./")
	expect.Equal(".", sourcePattern.Path)
	expect.Equal("", sourcePattern.Pattern)
	expect.True(sourcePattern.IsDir())
	expect.False(sourcePattern.IsFile())


}

// +build windows
func TestAbsoluteWindows(t *testing.T) {

	expect := assert.New(t)

	abs := "C:/"
	absWithPattern := abs +"NotExisting/*.txt"
	sourcePattern := NewBasePattern(absWithPattern)
	expect.Equal(filepath.ToSlash(abs), sourcePattern.Path)
	expect.Equal("NotExisting/*.txt", sourcePattern.Pattern)

}
