package newmatcher

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMaxAgeMatcher(t *testing.T) {
	expect := assert.New(t)

	fileToCheck := "../data/fixtures/global/file.txt"

	matcher := NewMaxAgeMatcher(time.Now())
	expect.False(matcher.Matches(fileToCheck))

	matcher = NewMaxAgeMatcher(time.Date(2015, 1, 1, 1, 1, 1, 1, time.Local))
	expect.True(matcher.Matches(fileToCheck))
}

func TestMaxAgeMatcherWithMissingFile(t *testing.T) {
	expect := assert.New(t)

	fileToCheck := "../data/fixtures/global/not-exists.txt"

	matcher := NewMaxAgeMatcher(time.Now())
	expect.False(matcher.Matches(fileToCheck))

	matcher = NewMaxAgeMatcher(time.Date(2015, 1, 1, 1, 1, 1, 1, time.Local))
	expect.False(matcher.Matches(fileToCheck))
}
