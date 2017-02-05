package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)



// all \ takes next char without replacement
// \* => \*

// * => .*
// . => \.
// $ at the end => *.jpg => (.*)\.jpg$

// /tmp/test.{htm,php} => test\.(htm|php)$

// specials
// !(pattern|pattern|pattern) => ^(?!ab)
// ?(pattern|pattern|pattern) => (pattern|pattern|pattern){0,1}
// +(pattern|pattern|pattern) => (pattern|pattern|pattern){1,}
// *(pattern|pattern|pattern) => (pattern|pattern|pattern){0,}
// @(pattern|pat*|pat?erN) => (pattern|pat*|pat?erN){1}

// if no sub-pattern / group is given, the pattern is treated as single group
// => /tmp/*.jpg => /tmp/(.*\.jpg)


func TestGlobToRegex(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(".*\\.jpg", GlobToRegex("*.jpg"))
	assert.Equal("star-file-\\*\\.jpg", GlobToRegex("star-file-\\*.jpg"))
	assert.Equal("test\\.(jpg|png)", GlobToRegex("test.{jpg,png}"))
}
