package newmatcher

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"regexp"
)

func TestRegexMatcher(t *testing.T) {
	expect := assert.New(t)

	compiledRegex, _ := regexp.Compile("^a(.*)$")

	matcher := NewRegexMatcher(*compiledRegex)

	expect.True(matcher.Matches("abcd"))
	expect.False(matcher.Matches("other value"))
}