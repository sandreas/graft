package newpattern

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGlobToRegex(t *testing.T) {
	expect := assert.New(t)
	expect.Equal(".*\\.jpg", GlobToRegexString("*.jpg"))
	expect.Equal("star-file-\\*\\.jpg", GlobToRegexString("star-file-\\*.jpg"))
	expect.Equal("test\\.(jpg|png)", GlobToRegexString("test.(jpg|png)"))
	expect.Equal("test\\.{1,}", GlobToRegexString("test.{1,}"))
	expect.Equal("fixtures\\(\\..*)", GlobToRegexString("fixtures\\(.*)"))
}