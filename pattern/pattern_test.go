package pattern_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/pattern"
	"regexp"
)

func TestNormalizeDirSep(t *testing.T) {
	expect := assert.New(t)
	expect.Equal("/tmp/dir/subdir", pattern.NormalizeDirSep("/tmp\\dir\\subdir"))
}

func TestParsePathPattern(t *testing.T) {
	expect := assert.New(t)

	path, pat := pattern.ParsePathPattern("../data/fixtures/global/*")
	expect.Equal("../data/fixtures/global", path)
	expect.Equal("*", pat)

	path, pat = pattern.ParsePathPattern("../data/fixtures/non-existing/*.*")
	expect.Equal("../data/fixtures", path)
	expect.Equal("non-existing/*.*", pat)
}

func TestGlobToRegex(t *testing.T) {
	expect := assert.New(t)
	expect.Equal(".*\\.jpg", pattern.GlobToRegex("*.jpg"))
	expect.Equal("star-file-\\*\\.jpg", pattern.GlobToRegex("star-file-\\*.jpg"))
	expect.Equal("test\\.(jpg|png)", pattern.GlobToRegex("test.{jpg,png}"))

	expect.Equal("fixtures\\(\\..*)", pattern.GlobToRegex("fixtures\\(.*)"))
}

func TestBuildMatchList(t *testing.T) {
	expect := assert.New(t)
	compiled, _ := regexp.Compile("data/fixtures/global/(.*)(\\.txt)$")

	list := pattern.BuildMatchList(compiled, "data/fixtures/global/documents (2010)/document (2010).txt")


	expect.Equal(2, len(list))
	expect.Equal("documents (2010)/document (2010)", list[0])
	expect.Equal(".txt", list[1])
}

func TestCompileNormalizedPathPattern(t *testing.T) {
	expect := assert.New(t)
	compiled, _ := pattern.CompileNormalizedPathPattern("data\\fixtures/global", "(.*)")
	expect.Equal("data/fixtures/global/(.*)", compiled.String())
	expect.Regexp(compiled, "data/fixtures/global/test.txt")


	compiled, _ = pattern.CompileNormalizedPathPattern("", "(.*\\.jpg)")
	expect.Equal("(.*\\.jpg)", compiled.String())
	expect.Regexp(compiled, "data/fixtures/global/test.jpg")

}