package matcher_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/matcher"
)

type TestMatcher struct {
	matcher.MatcherInterface
}
func (f *TestMatcher) Matches(subject interface{}) bool {
	return subject == "test"
}


func TestCompositeMatcher(t *testing.T) {
	expect := assert.New(t)
	testMatcher := &TestMatcher{}

	subject := matcher.NewCompositeMatcher()
	subject.Add(testMatcher)

	expect.True(subject.Matches("test"))
	expect.False(subject.Matches("other value"))
}
