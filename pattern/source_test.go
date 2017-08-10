package pattern

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"regexp"
)

func TestNewSourcePattern(t *testing.T) {
	expect := assert.New(t)

	sourcePattern := NewSourcePattern("../data/fixtures/global/*")
	expect.Equal("../data/fixtures/global", sourcePattern.Path)
	expect.Equal("*", sourcePattern.Pattern)
	expect.True(sourcePattern.IsDir())
	expect.False(sourcePattern.IsFile())

	sourcePattern = NewSourcePattern("../data/fixtures/non-existing/*.*")
	expect.Equal("../data/fixtures", sourcePattern.Path)
	expect.Equal("non-existing/*.*", sourcePattern.Pattern)
	expect.True(sourcePattern.IsDir())
	expect.False(sourcePattern.IsFile())

	sourcePattern = NewSourcePattern("../data/fixtures/global/file.txt")
	expect.Equal("../data/fixtures/global/file.txt", sourcePattern.Path)
	expect.Equal("", sourcePattern.Pattern)
	expect.False(sourcePattern.IsDir())
	expect.True(sourcePattern.IsFile())

	sourcePattern = NewSourcePattern("../data/fixtures/global/")
	expect.Equal("../data/fixtures/global", sourcePattern.Path)
	expect.Equal("", sourcePattern.Pattern)
	expect.True(sourcePattern.IsDir())
	expect.False(sourcePattern.IsFile())

	sourcePattern = NewSourcePattern("*")
	expect.Equal(".", sourcePattern.Path)
	expect.Equal("*", sourcePattern.Pattern)
	expect.True(sourcePattern.IsDir())
	expect.False(sourcePattern.IsFile())

	sourcePattern = NewSourcePattern("./*")
	expect.Equal(".", sourcePattern.Path)
	expect.Equal("*", sourcePattern.Pattern)
	expect.True(sourcePattern.IsDir())
	expect.False(sourcePattern.IsFile())

	sourcePattern = NewSourcePattern("../data/fixtures/global/(.*)")
	expect.Equal("../data/fixtures/global", sourcePattern.Path)
	expect.Equal("(.*)", sourcePattern.Pattern)
	expect.True(sourcePattern.IsDir())
	expect.False(sourcePattern.IsFile())

	sourcePattern = NewSourcePattern(".")
	expect.Equal(".", sourcePattern.Path)
	expect.Equal("", sourcePattern.Pattern)
	expect.True(sourcePattern.IsDir())
	expect.False(sourcePattern.IsFile())

	sourcePattern = NewSourcePattern("./")
	expect.Equal(".", sourcePattern.Path)
	expect.Equal("", sourcePattern.Pattern)
	expect.True(sourcePattern.IsDir())
	expect.False(sourcePattern.IsFile())
}

func TestCompileGlob(t *testing.T) {
	expect := assert.New(t)
	var compiled *regexp.Regexp
	compiled, _ = NewSourcePattern("../data/fixtures/global/*").Compile()
	expect.Equal("(?i)\\.\\./data/fixtures/global/(.*)$", compiled.String())
	expect.Regexp(compiled, "../data/fixtures/global/test.txt")

	compiled, _ = NewSourcePattern("../data/fixtures/global/").Compile()
	expect.Equal("(?i)\\.\\./data/fixtures/global/(.*)$", compiled.String())
	expect.Regexp(compiled, "../data/fixtures/global/test.txt")

	compiled, _ = NewSourcePattern("../data/fixtures/global/t(*)t.(txt)").Compile()
	expect.Equal("(?i)\\.\\./data/fixtures/global/t(.*)t\\.(txt)$", compiled.String())
	expect.Regexp(compiled, "../data/fixtures/global/Test.txt")

	compiled, _ = NewSourcePattern("../data/fixtures/global/t(*)t.(txt)", CASE_SENSITIVE).Compile()
	expect.Equal("\\.\\./data/fixtures/global/t(.*)t\\.(txt)$", compiled.String())
	expect.NotRegexp(compiled, "../data/fixtures/global/Test.txt")

	pattern := NewSourcePattern("../data/fixtures/global/.*.?", CASE_SENSITIVE|USE_REAL_REGEX)
	compiled, _ = pattern.Compile()
	expect.Equal("\\.\\./data/fixtures/global/(.*.?)$", compiled.String())

	pattern = NewSourcePattern(".")
	compiled, _ = pattern.Compile()
	expect.Equal("(?i)(.*)$", compiled.String())

	pattern = NewSourcePattern("./")
	compiled, _ = pattern.Compile()
	expect.Equal("(?i)(.*)$", compiled.String())
}
