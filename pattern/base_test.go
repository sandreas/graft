package pattern_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/pattern"
	"github.com/sandreas/graft/testhelpers"
	"path/filepath"
	"runtime"
)

func TestBase(t *testing.T) {
	expect := assert.New(t)

	_, filename, _, _ := runtime.Caller(0)
	currentFilePath := filepath.ToSlash(filepath.Dir(filename))

	// long paths on windows have to be converted to absolute ones
	veryLongRelativePath := "inetpub/wwwroot/something_4.0/node_modules/babel-preset-es2015/node_modules/babel-plugin-transform-es2015-block-scoping/node_modules/babel-traverse/node_modules/babel-code-frame/node_modules/chalk/node_modules/strip-ansi/node_modules/ansi-regex/node_modules/fake-sub-module";
	veryLongAbsolutePath := currentFilePath + "/" + veryLongRelativePath
	veryLongRelativePathFile := veryLongRelativePath + "/package.json"
	veryLongAbsolutePathFile := veryLongAbsolutePath + "/package.json"

	mockFs := testhelpers.MockFileSystem(map[string]string{
		"fixtures/global/":         "",
		"fixtures/global/file.txt": "",
		veryLongRelativePathFile: "{}",
		veryLongAbsolutePathFile: "{}",
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

	basePattern = pattern.NewBasePattern(mockFs, veryLongRelativePath+"/*")
	if runtime.GOOS == "windows" {
		expect.Equal(veryLongAbsolutePath, basePattern.Path)
	} else {
		expect.Equal(veryLongRelativePath, basePattern.Path)
	}
	expect.Equal("*", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())
}
