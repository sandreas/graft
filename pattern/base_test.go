package pattern_test

import (
	"testing"

	"os"

	"github.com/sandreas/graft/pattern"
	"github.com/sandreas/graft/testhelpers"
	"github.com/stretchr/testify/assert"
	"runtime"
	"strings"
	"github.com/spf13/afero"
)

func TestBase(t *testing.T) {
	sep := string(os.PathSeparator)
	expect := assert.New(t)

	// long paths on windows have to be converted to absolute ones
	veryLongRelativePath := "inetpub/wwwroot/something_4.0/node_modules/babel-preset-es2015/node_modules/babel-plugin-transform-es2015-block-scoping/node_modules/babel-traverse/node_modules/babel-code-frame/node_modules/chalk/node_modules/strip-ansi/node_modules/ansi-regex/node_modules/fake-sub-module"
	veryLongRelativePathFile := veryLongRelativePath + "/package.json"

	uncPath := sep + sep + "unc-server" + sep + "unc-share" + sep + "file.txt"

	mockFs := testhelpers.MockFileSystem(map[string]string{
		"fixtures/global/":         "",
		"fixtures/global/file.txt": "",
		veryLongRelativePathFile:   "{}",
		"C:" + sep + "Temp":        "",
		uncPath:                    "file-content",
	})

	var basePattern *pattern.BasePattern

	basePattern = pattern.NewBasePattern(mockFs, "fixtures/global/*")
	expect.Equal("fixtures"+sep+"global"+sep, basePattern.Path)
	expect.Equal("*", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = pattern.NewBasePattern(mockFs, "fixtures/global/(.*)")
	expect.Equal("fixtures"+sep+"global"+sep, basePattern.Path)
	expect.Equal("(.*)", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = pattern.NewBasePattern(mockFs, "fixtures/global/file.txt")
	expect.Equal("fixtures"+sep+"global"+sep+"file.txt", basePattern.Path)
	expect.Equal("", basePattern.Pattern)
	expect.False(basePattern.IsDir())
	expect.True(basePattern.IsFile())

	basePattern = pattern.NewBasePattern(mockFs, "fixtures/global/")
	expect.Equal("fixtures"+sep+"global"+sep, basePattern.Path)
	expect.Equal("", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = pattern.NewBasePattern(mockFs, "fixtures/non-existing/*.*")
	expect.Equal("fixtures"+sep, basePattern.Path)
	expect.Equal("non-existing/*.*", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = pattern.NewBasePattern(mockFs, "*")
	expect.Equal("."+sep, basePattern.Path)
	expect.Equal("*", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = pattern.NewBasePattern(mockFs, "")
	expect.Equal("."+sep, basePattern.Path)
	expect.Equal("", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = pattern.NewBasePattern(mockFs, "./*")
	expect.Equal("."+sep, basePattern.Path)
	expect.Equal("*", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = pattern.NewBasePattern(mockFs, ".")
	expect.Equal("."+sep, basePattern.Path)
	expect.Equal("", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = pattern.NewBasePattern(mockFs, "./")
	expect.Equal("."+sep, basePattern.Path)
	expect.Equal("", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = pattern.NewBasePattern(mockFs, veryLongRelativePath+"/*")
	expect.Equal("*", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = pattern.NewBasePattern(mockFs, "C:/TestMissingDir/")
	expect.Equal("C:"+sep, basePattern.Path)
	expect.Equal("TestMissingDir/", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = pattern.NewBasePattern(mockFs, "(*)fi(*).txt")
	expect.Equal("."+sep, basePattern.Path)
	expect.Equal("(*)fi(*).txt", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

	basePattern = pattern.NewBasePattern(mockFs, uncPath)
	if runtime.GOOS == "windows" {
		expect.Equal("\\\\unc-server\\unc-share\\file.txt", basePattern.Path)
		expect.Equal("", basePattern.Pattern)
		expect.False(basePattern.IsDir())
		expect.True(basePattern.IsFile())
	} else {
		expect.Equal("/unc-server/unc-share/file.txt", basePattern.Path)
		expect.Equal("", basePattern.Pattern)
		expect.False(basePattern.IsDir())
		expect.True(basePattern.IsFile())
	}

	slashedUncPath := strings.Replace(uncPath, "\\", "/", -1)
	basePattern = pattern.NewBasePattern(mockFs, slashedUncPath)
	if runtime.GOOS == "windows" {
		expect.Equal("\\\\unc-server\\unc-share\\file.txt", basePattern.Path)
		expect.Equal("", basePattern.Pattern)
		expect.False(basePattern.IsDir())
		expect.True(basePattern.IsFile())
	} else {
		expect.Equal("/unc-server/unc-share/file.txt", basePattern.Path)
		expect.Equal("", basePattern.Pattern)
		expect.False(basePattern.IsDir())
		expect.True(basePattern.IsFile())
	}

	basePattern = pattern.NewBasePattern(mockFs, "/")
	expect.Equal(sep, basePattern.Path)
	expect.Equal("", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())

}

func TestBaseOsFs(t *testing.T) {
	expect := assert.New(t)

	basePattern := pattern.NewBasePattern(afero.OsFs{}, "")
	expect.Equal(".\\", basePattern.Path)
	expect.Equal("", basePattern.Pattern)
	expect.True(basePattern.IsDir())
	expect.False(basePattern.IsFile())
}
