package pattern

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/spf13/afero"
)
var mockFs afero.Fs
func init() {
	mockFs = afero.NewMemMapFs()
	mockFs.Mkdir("fixtures/global/", 0644)
	afero.WriteFile(mockFs,"fixtures/global/file.txt", []byte(""), 0755)
}

func TestBase(t *testing.T) {
	expect := assert.New(t)

	basePattern := NewBasePattern(mockFs, "fixtures/global/*")
	expect.Equal("fixtures/global", basePattern.Path)
	expect.Equal("*", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = NewBasePattern(mockFs, "fixtures/non-existing/*.*")
	expect.Equal("fixtures", basePattern.Path)
	expect.Equal("non-existing/*.*", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = NewBasePattern(mockFs, "fixtures/global/file.txt")
	expect.Equal("fixtures/global/file.txt", basePattern.Path)
	expect.Equal("", basePattern.Pattern)
	expect.False(basePattern.IsDir())
	expect.True(basePattern.IsFile())

	basePattern = NewBasePattern(mockFs, "fixtures/global/")
	expect.Equal("fixtures/global", basePattern.Path)
	expect.Equal("", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = NewBasePattern(mockFs, "*")
	expect.Equal(".", basePattern.Path)
	expect.Equal("*", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = NewBasePattern(mockFs, "./*")
	expect.Equal(".", basePattern.Path)
	expect.Equal("*", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = NewBasePattern(mockFs, "fixtures/global/(.*)")
	expect.Equal("fixtures/global", basePattern.Path)
	expect.Equal("(.*)", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = NewBasePattern(mockFs, ".")
	expect.Equal(".", basePattern.Path)
	expect.Equal("", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = NewBasePattern(mockFs, "./")
	expect.Equal(".", basePattern.Path)
	expect.Equal("", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())
}
