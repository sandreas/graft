package pattern_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/pattern"
	"github.com/sandreas/graft/testhelpers"
)

func TestBase(t *testing.T) {
	expect := assert.New(t)

	mockFs := testhelpers.MockFileSystem(map[string]string{
		"fixtures/global/":         "",
		"fixtures/global/file.txt": "",
	})

	basePattern := pattern.NewBasePattern(mockFs, "fixtures/global/*")
	expect.Equal("fixtures/global", basePattern.Path)
	expect.Equal("*", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = pattern.NewBasePattern(mockFs, "fixtures/non-existing/*.*")
	expect.Equal("fixtures", basePattern.Path)
	expect.Equal("non-existing/*.*", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = pattern.NewBasePattern(mockFs, "fixtures/global/file.txt")
	expect.Equal("fixtures/global/file.txt", basePattern.Path)
	expect.Equal("", basePattern.Pattern)
	expect.False(basePattern.IsDir())
	expect.True(basePattern.IsFile())

	basePattern = pattern.NewBasePattern(mockFs, "fixtures/global/")
	expect.Equal("fixtures/global", basePattern.Path)
	expect.Equal("", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = pattern.NewBasePattern(mockFs, "*")
	expect.Equal(".", basePattern.Path)
	expect.Equal("*", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = pattern.NewBasePattern(mockFs, "")
	expect.Equal(".", basePattern.Path)
	expect.Equal("", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = pattern.NewBasePattern(mockFs, "./*")
	expect.Equal(".", basePattern.Path)
	expect.Equal("*", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = pattern.NewBasePattern(mockFs, "fixtures/global/(.*)")
	expect.Equal("fixtures/global", basePattern.Path)
	expect.Equal("(.*)", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = pattern.NewBasePattern(mockFs, ".")
	expect.Equal(".", basePattern.Path)
	expect.Equal("", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = pattern.NewBasePattern(mockFs, "./")
	expect.Equal(".", basePattern.Path)
	expect.Equal("", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())
}
