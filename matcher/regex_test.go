package matcher_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"regexp"
	"github.com/sandreas/graft/matcher"
)

func TestRegexMatcher(t *testing.T) {
	expect := assert.New(t)

	compiledRegex, _ := regexp.Compile("^a(.*)$")

	subject := matcher.NewRegexMatcher(compiledRegex)

	expect.True(subject.Matches("abcd"))
	expect.False(subject.Matches("other value"))
}