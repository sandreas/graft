package pattern_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/testhelpers"
	"regexp"
	"github.com/sandreas/graft/pattern"
	"os"
)
var sep = string(os.PathSeparator)

func TestNewSourcePattern(t *testing.T) {
	expect := assert.New(t)

	mockFs := testhelpers.MockFileSystem(map[string]string{
		"fixtures/global/":         "",
		"fixtures/global/file.txt": "",
	})

	sourcePattern := pattern.NewSourcePattern(mockFs, "fixtures/global/*")
	expect.Equal("fixtures"+sep+"global", sourcePattern.Path)
	expect.Equal("*", sourcePattern.Pattern)
	expect.True(sourcePattern.IsDir())
	expect.False(sourcePattern.IsFile())

	sourcePattern = pattern.NewSourcePattern(mockFs, "fixtures/non-existing/*.*")
	expect.Equal("fixtures", sourcePattern.Path)
	expect.Equal("non-existing/*.*", sourcePattern.Pattern)
	expect.True(sourcePattern.IsDir())
	expect.False(sourcePattern.IsFile())

	sourcePattern = pattern.NewSourcePattern(mockFs, "fixtures/global/file.txt")
	expect.Equal("fixtures"+sep+"global"+sep+"file.txt", sourcePattern.Path)
	expect.Equal("", sourcePattern.Pattern)
	expect.False(sourcePattern.IsDir())
	expect.True(sourcePattern.IsFile())

	sourcePattern = pattern.NewSourcePattern(mockFs, "fixtures/global/")
	expect.Equal("fixtures"+sep+"global", sourcePattern.Path)
	expect.Equal("", sourcePattern.Pattern)
	expect.True(sourcePattern.IsDir())
	expect.False(sourcePattern.IsFile())

	sourcePattern = pattern.NewSourcePattern(mockFs, "*")
	expect.Equal(".", sourcePattern.Path)
	expect.Equal("*", sourcePattern.Pattern)
	expect.True(sourcePattern.IsDir())
	expect.False(sourcePattern.IsFile())

	sourcePattern = pattern.NewSourcePattern(mockFs, "./*")
	expect.Equal(".", sourcePattern.Path)
	expect.Equal("*", sourcePattern.Pattern)
	expect.True(sourcePattern.IsDir())
	expect.False(sourcePattern.IsFile())

	sourcePattern = pattern.NewSourcePattern(mockFs, "fixtures/global/(.*)")
	expect.Equal("fixtures"+sep+"global", sourcePattern.Path)
	expect.Equal("(.*)", sourcePattern.Pattern)
	expect.True(sourcePattern.IsDir())
	expect.False(sourcePattern.IsFile())

	sourcePattern = pattern.NewSourcePattern(mockFs, ".")
	expect.Equal(".", sourcePattern.Path)
	expect.Equal("", sourcePattern.Pattern)
	expect.True(sourcePattern.IsDir())
	expect.False(sourcePattern.IsFile())

	sourcePattern = pattern.NewSourcePattern(mockFs, "./")
	expect.Equal(".", sourcePattern.Path)
	expect.Equal("", sourcePattern.Pattern)
	expect.True(sourcePattern.IsDir())
	expect.False(sourcePattern.IsFile())
}

func TestCompileSimple(t *testing.T) {
	expect := assert.New(t)
	mockFs := testhelpers.MockFileSystem(map[string]string{
		"fixtures/global/":         "",
		"fixtures/global/file.txt": "",
	})

	var compiled *regexp.Regexp
	compiled, _ = pattern.NewSourcePattern(mockFs, "fixtures/global/file.txt").Compile()
	expect.Equal("(?i)fixtures"+sep+sep+"global"+sep+sep+"file\\.txt$", compiled.String())
	expect.Regexp(compiled, "fixtures"+sep+"global"+sep+"file.txt")

}

func TestCompileGlob(t *testing.T) {
	expect := assert.New(t)
	mockFs := testhelpers.MockFileSystem(map[string]string{
		"fixtures/global/":         "",
		"fixtures/global/file.txt": "",
	})
	var compiled *regexp.Regexp
	compiled, _ = pattern.NewSourcePattern(mockFs, "fixtures/global/*").Compile()
	expect.Equal("(?i)fixtures"+sep+sep+"global"+sep+sep+"(.*)$", compiled.String())
	expect.Regexp(compiled, "fixtures"+sep+"global"+sep+"test.txt")

	compiled, _ = pattern.NewSourcePattern(mockFs, "fixtures/global/").Compile()
	expect.Equal("(?i)fixtures"+sep+sep+"global"+sep+sep+"(.*)$", compiled.String())
	expect.Regexp(compiled, "fixtures"+sep+"global"+sep+"test.txt")

	compiled, _ = pattern.NewSourcePattern(mockFs, "fixtures/global/t(*)t.(txt)").Compile()
	expect.Equal("(?i)fixtures"+sep+sep+"global"+sep+sep+"t(.*)t\\.(txt)$", compiled.String())
	expect.Regexp(compiled, "fixtures"+sep+"global"+sep+"Test.txt")

	compiled, _ = pattern.NewSourcePattern(mockFs, "fixtures/global/t(*)t.(txt)", pattern.CASE_SENSITIVE).Compile()
	expect.Equal("fixtures"+sep+sep+"global"+sep+sep+"t(.*)t\\.(txt)$", compiled.String())
	expect.NotRegexp(compiled, "fixtures"+sep+"global"+sep+"Test.txt")

	sourcePattern := pattern.NewSourcePattern(mockFs, "fixtures/global/.*.?", pattern.CASE_SENSITIVE|pattern.USE_REAL_REGEX)
	compiled, _ = sourcePattern.Compile()
	expect.Equal("fixtures"+sep+sep+"global"+sep+sep+"(.*.?)$", compiled.String())

	sourcePattern = pattern.NewSourcePattern(mockFs, ".")
	compiled, _ = sourcePattern.Compile()
	expect.Equal("(?i)(.*)$", compiled.String())

	sourcePattern = pattern.NewSourcePattern(mockFs, "./")
	compiled, _ = sourcePattern.Compile()
	expect.Equal("(?i)(.*)$", compiled.String())
}
