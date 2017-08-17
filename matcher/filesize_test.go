package matcher_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/matcher"
	"github.com/sandreas/graft/testhelpers"
)

func TestSizeMatcherForFileWithStat(t *testing.T) {
	expect := assert.New(t)
	fileToCheck := "file.txt"

	mockFs := testhelpers.MockFileSystem(map[string]string{
		fileToCheck: "file",
	})

	m := matcher.NewFileSizeMatcher(-1, 2)
	m.Fs = mockFs
	expect.False(m.Matches(fileToCheck))

	m = matcher.NewFileSizeMatcher(-1, 4)
	m.Fs = mockFs
	expect.True(m.Matches(fileToCheck))

	m = matcher.NewFileSizeMatcher(-1, 15)
	m.Fs = mockFs
	expect.True(m.Matches(fileToCheck))

	m = matcher.NewFileSizeMatcher(0, -1)
	m.Fs = mockFs
	expect.True(m.Matches(fileToCheck))

	m = matcher.NewFileSizeMatcher(5, -1)
	m.Fs = mockFs
	expect.False(m.Matches(fileToCheck))

	m = matcher.NewFileSizeMatcher(3, 4)
	m.Fs = mockFs
	expect.True(m.Matches(fileToCheck))

	m = matcher.NewFileSizeMatcher(3, 3)
	m.Fs = mockFs
	expect.False(m.Matches(fileToCheck))

	m = matcher.NewFileSizeMatcher(5, 4)
	m.Fs = mockFs
	expect.False(m.Matches(fileToCheck))

	m = matcher.NewFileSizeMatcher(-1, -1)
	m.Fs = mockFs
	expect.False(m.Matches(fileToCheck))
}

func TestSizeMatcherForFile(t *testing.T) {
	expect := assert.New(t)
	fileToCheck := "../data/fixtures/global/file.txt"
	m := matcher.NewFileSizeMatcher(3, 4)
	expect.True(m.Matches(fileToCheck))
}

func TestSizeMatcherForDirWithStat(t *testing.T) {
	expect := assert.New(t)
	dirToCheck := "fixtures/"

	mockFs := testhelpers.MockFileSystem(map[string]string{
		dirToCheck: "",
	})

	m := matcher.NewFileSizeMatcher(0, 5)
	m.Fs = mockFs
	expect.False(m.Matches(dirToCheck))
}
