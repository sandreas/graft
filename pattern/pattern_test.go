package pattern_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/pattern"
	"regexp"
)

func TestNormalizeDirSep(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("/tmp/dir/subdir", pattern.NormalizeDirSep("/tmp\\dir\\subdir"))
}

func TestParsePathPattern(t *testing.T) {
	assert := assert.New(t)

	path, pat := pattern.ParsePathPattern("../data/fixtures/global/*")
	assert.Equal("../data/fixtures/global", path)
	assert.Equal("*", pat)

	path, pat = pattern.ParsePathPattern("../data/fixtures/non-existing/*.*")
	assert.Equal("../data/fixtures", path)
	assert.Equal("non-existing/*.*", pat)
}

func TestGlobToRegex(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(".*\\.jpg", pattern.GlobToRegex("*.jpg"))
	assert.Equal("star-file-\\*\\.jpg", pattern.GlobToRegex("star-file-\\*.jpg"))
	assert.Equal("test\\.(jpg|png)", pattern.GlobToRegex("test.{jpg,png}"))
}

func TestBuildMatchList(t *testing.T) {
	assert := assert.New(t)
	compiled, _ := regexp.Compile("data/fixtures/global/(.*)(\\.txt)$")

	list := pattern.BuildMatchList(compiled, "data/fixtures/global/documents (2010)/document (2010).txt")


	assert.Equal(2, len(list))
	assert.Equal("documents (2010)/document (2010)", list[0])
	assert.Equal(".txt", list[1])
}

func TestCompileNormalizedPathPattern(t *testing.T) {
	assert := assert.New(t)
	compiled, _ := pattern.CompileNormalizedPathPattern("data\\fixtures/global", "(.*)")
	assert.Regexp(compiled, "data/fixtures/global/test.txt")
}