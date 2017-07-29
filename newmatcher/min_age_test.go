package newmatcher

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMinAgeMatcher(t *testing.T) {
	expect := assert.New(t)

	fileToCheck := "../data/fixtures/global/file.txt"

	matcher := NewMinAgeMatcher(time.Now())
	expect.True(matcher.Matches(fileToCheck))

	matcher = NewMinAgeMatcher(time.Date(2015, 1, 1, 1, 1, 1, 1, time.Local))
	expect.False(matcher.Matches(fileToCheck))
}

func TestMinAgeMatcherWithMissingFile(t *testing.T) {
	expect := assert.New(t)

	fileToCheck := "../data/fixtures/global/not-exists.txt"

	matcher := NewMinAgeMatcher(time.Now())
	expect.False(matcher.Matches(fileToCheck))

	matcher = NewMinAgeMatcher(time.Date(2015, 1, 1, 1, 1, 1, 1, time.Local))
	expect.False(matcher.Matches(fileToCheck))
}
