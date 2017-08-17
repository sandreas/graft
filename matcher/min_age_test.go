package matcher_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/matcher"
)

func TestMinAgeMatcher(t *testing.T) {
	expect := assert.New(t)

	fileToCheck := "../data/fixtures/global/file.txt"

	subject := matcher.NewMinAgeMatcher(time.Now())
	expect.True(subject.Matches(fileToCheck))

	subject = matcher.NewMinAgeMatcher(time.Date(2015, 1, 1, 1, 1, 1, 1, time.Local))
	expect.False(subject.Matches(fileToCheck))
}

func TestMinAgeMatcherWithMissingFile(t *testing.T) {
	expect := assert.New(t)

	fileToCheck := "../data/fixtures/global/not-exists.txt"

	subject := matcher.NewMinAgeMatcher(time.Now())
	expect.False(subject.Matches(fileToCheck))

	subject = matcher.NewMinAgeMatcher(time.Date(2015, 1, 1, 1, 1, 1, 1, time.Local))
	expect.False(subject.Matches(fileToCheck))
}
