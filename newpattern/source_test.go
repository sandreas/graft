package newpattern

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewSourcePattern(t *testing.T) {
	expect := assert.New(t)
	sourcePattern := NewSourcePattern("../data/fixtures/global/*")
	expect.Equal("../data/fixtures/global", sourcePattern.Path)
	expect.Equal("*", sourcePattern.Pattern)
	expect.True(sourcePattern.IsDirectory)
}

func TestParse(t *testing.T) {
	expect := assert.New(t)

	sourcePattern := &SourcePattern{}

	sourcePattern.Parse("../data/fixtures/global/*")
	expect.Equal("../data/fixtures/global", sourcePattern.Path)
	expect.Equal("*", sourcePattern.Pattern)
	expect.True(sourcePattern.IsDirectory)

	sourcePattern.Parse("../data/fixtures/non-existing/*.*")
	expect.Equal("../data/fixtures", sourcePattern.Path)
	expect.Equal("non-existing/*.*", sourcePattern.Pattern)
	expect.True(sourcePattern.IsDirectory)

	sourcePattern.Parse("../data/fixtures/global/file.txt")
	expect.Equal("../data/fixtures/global/file.txt", sourcePattern.Path)
	expect.Equal("", sourcePattern.Pattern)
	expect.False(sourcePattern.IsDirectory)

	sourcePattern.Parse("*")
	expect.Equal(".", sourcePattern.Path)
	expect.Equal("*", sourcePattern.Pattern)
	expect.True(sourcePattern.IsDirectory)

	sourcePattern.Parse("./*")
	expect.Equal(".", sourcePattern.Path)
	expect.Equal("*", sourcePattern.Pattern)
	expect.True(sourcePattern.IsDirectory)
}