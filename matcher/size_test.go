package matcher_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/matcher"
)

func TestSizeMatcher(t *testing.T) {
	expect := assert.New(t)

	fileToCheck := "../data/fixtures/global/file.txt"

	m := matcher.NewSizeMatcher(-1, 2)
	expect.False(m.Matches(fileToCheck))

	m = matcher.NewSizeMatcher(-1, 4)
	expect.True(m.Matches(fileToCheck))

	m = matcher.NewSizeMatcher(-1, 15)
	expect.True(m.Matches(fileToCheck))

	m = matcher.NewSizeMatcher(0, -1)
	expect.True(m.Matches(fileToCheck))

	m = matcher.NewSizeMatcher(5, -1)
	expect.False(m.Matches(fileToCheck))


	m = matcher.NewSizeMatcher(3, 4)
	expect.True(m.Matches(fileToCheck))

	m = matcher.NewSizeMatcher(3, 3)
	expect.False(m.Matches(fileToCheck))

	m = matcher.NewSizeMatcher(5, 4)
	expect.False(m.Matches(fileToCheck))

	m = matcher.NewSizeMatcher(-1, -1)
	expect.False(m.Matches(fileToCheck))


	dirToCheck := "../data/fixtures/global"
	m = matcher.NewSizeMatcher(0, 5)
	expect.False(m.Matches(dirToCheck))


}

