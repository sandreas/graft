package newmatcher

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

type TestMatcher struct {
	MatcherInterface
}
func (f *TestMatcher) Matches(subject interface{}) bool {
	return subject == "test"
}


func TestCompositeMatcher(t *testing.T) {
	expect := assert.New(t)
	testMatcher := &TestMatcher{}

	matcher := NewCompositeMatcher()
	matcher.Add(testMatcher)

	expect.True(matcher.Matches("test"))
	expect.False(matcher.Matches("other value"))
}
