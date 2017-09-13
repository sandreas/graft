package pattern_test

import (
	"testing"

	"github.com/sandreas/graft/pattern"
	"github.com/sandreas/graft/testhelpers"
	"github.com/stretchr/testify/assert"
	"os"
)

func TestBase(t *testing.T) {
	sep := string(os.PathSeparator)
	expect := assert.New(t)

	// long paths on windows have to be converted to absolute ones
	veryLongRelativePath := "inetpub/wwwroot/something_4.0/node_modules/babel-preset-es2015/node_modules/babel-plugin-transform-es2015-block-scoping/node_modules/babel-traverse/node_modules/babel-code-frame/node_modules/chalk/node_modules/strip-ansi/node_modules/ansi-regex/node_modules/fake-sub-module"
	veryLongRelativePathFile := veryLongRelativePath + "/package.json"

	mockFs := testhelpers.MockFileSystem(map[string]string{
		"fixtures/global/":         "",
		"fixtures/global/file.txt": "",
		veryLongRelativePathFile:   "{}",
		"C:" + sep + "Temp" : "",
	})

	var basePattern *pattern.BasePattern
	basePattern = pattern.NewBasePattern(mockFs, "fixtures/global/*")
	expect.Equal("fixtures"+sep+"global", basePattern.Path)
	expect.Equal("*", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = pattern.NewBasePattern(mockFs, "fixtures/global/(.*)")
	expect.Equal("fixtures"+sep+"global", basePattern.Path)
	expect.Equal("(.*)", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = pattern.NewBasePattern(mockFs, "fixtures/global/file.txt")
	expect.Equal("fixtures"+sep+"global"+sep+"file.txt", basePattern.Path)
	expect.Equal("", basePattern.Pattern)
	expect.False(basePattern.IsDir())
	expect.True(basePattern.IsFile())

	basePattern = pattern.NewBasePattern(mockFs, "fixtures/global/")
	expect.Equal("fixtures"+sep+"global", basePattern.Path)
	expect.Equal("", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = pattern.NewBasePattern(mockFs, "fixtures/non-existing/*.*")
	expect.Equal("fixtures", basePattern.Path)
	expect.Equal("non-existing/*.*", basePattern.Pattern)
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

	basePattern = pattern.NewBasePattern(mockFs, veryLongRelativePath+"/*")
	expect.Equal("*", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = pattern.NewBasePattern(mockFs, "C:/TestMissingDir/")
	expect.Equal("C:", basePattern.Path)
	expect.Equal("TestMissingDir/", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())


	basePattern = pattern.NewBasePattern(mockFs, "(*)fi(*).txt")
	expect.Equal(".", basePattern.Path)
	expect.Equal("(*)fi(*).txt", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())


}
