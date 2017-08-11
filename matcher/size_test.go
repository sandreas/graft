package matcher

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSizeMatcher(t *testing.T) {
	expect := assert.New(t)

	fileToCheck := "../data/fixtures/global/file.txt"

	m := NewSizeMatcher(-1, 2)
	expect.False(m.Matches(fileToCheck))

	m = NewSizeMatcher(-1, 4)
	expect.True(m.Matches(fileToCheck))

	m = NewSizeMatcher(-1, 15)
	expect.True(m.Matches(fileToCheck))

	m = NewSizeMatcher(0, -1)
	expect.True(m.Matches(fileToCheck))

	m = NewSizeMatcher(5, -1)
	expect.False(m.Matches(fileToCheck))


	m = NewSizeMatcher(3, 4)
	expect.True(m.Matches(fileToCheck))

	m = NewSizeMatcher(3, 3)
	expect.False(m.Matches(fileToCheck))

	m = NewSizeMatcher(5, 4)
	expect.False(m.Matches(fileToCheck))

	m = NewSizeMatcher(-1, -1)
	expect.False(m.Matches(fileToCheck))


	dirToCheck := "../data/fixtures/global"
	m = NewSizeMatcher(0, 5)
	expect.False(m.Matches(dirToCheck))


}

