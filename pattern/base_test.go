package pattern_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/pattern"
	"github.com/sandreas/graft/testhelpers"
	"os"
	"path/filepath"
)

func TestBase(t *testing.T) {
	expect := assert.New(t)

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.ToSlash(filepath.Dir(ex))


	mockFs := testhelpers.MockFileSystem(map[string]string{
		"fixtures/global/":         "",
		"fixtures/global/file.txt": "",
		exPath + "/inetpub/wwwroot/something_4.0/node_modules/babel-preset-es2015/node_modules/babel-plugin-transform-es2015-block-scoping/node_modules/babel-traverse/node_modules/babel-code-frame/node_modules/chalk/node_modules/strip-ansi/node_modules/ansi-regex/node_modules/fake-sub-module/package.json":"{}",
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

	// C:/inetpub/wwwroot/something_4.0/node_modules/babel-preset-es2015/node_modules/babel-plugin-transform-es2015-block-scoping/node_modules/babel-traverse/node_modules/babel-code-frame/node_modules/chalk/node_modules/strip-ansi/node_modules/ansi-regex/


	//basePattern := pattern.NewBasePattern(mockFs, "inetpub/wwwroot/something_4.0/node_modules/babel-preset-es2015/node_modules/babel-plugin-transform-es2015-block-scoping/node_modules/babel-traverse/node_modules/babel-code-frame/node_modules/chalk/node_modules/strip-ansi/node_modules/ansi-regex/node_modules/fake-sub-module/*")
	//expect.Equal(exPath +"/inetpub/wwwroot/something_4.0/node_modules/babel-preset-es2015/node_modules/babel-plugin-transform-es2015-block-scoping/node_modules/babel-traverse/node_modules/babel-code-frame/node_modules/chalk/node_modules/strip-ansi/node_modules/ansi-regex/node_modules/fake-sub-module", basePattern.Path)
	//expect.Equal("*", basePattern.Pattern)
	//expect.True(basePattern.IsDir())
	//expect.False(basePattern.IsFile())
}
