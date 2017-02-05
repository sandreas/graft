package main

import "testing"



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

func expectHelper(expected string, actual string, t *testing.T) {
	if(actual != expected) {
		t.Error("expected " + expected + ", got " + actual)
	}
}

func TestGlobToRegex(t *testing.T) {
	expectHelper(".*\\.jpg", GlobToRegex("*.jpg"), t)
	expectHelper("star-file-\\*\\.jpg", GlobToRegex("star-file-\\*.jpg"), t)
	expectHelper("test\\.(jpg|png)", GlobToRegex("test.{jpg,png}"), t)
}
