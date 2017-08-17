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

	fi, _ := mockFs.Stat(fileToCheck)

	m := matcher.NewFileSizeMatcher(fi, -1, 2)
	expect.False(m.Matches(fileToCheck))

	m = matcher.NewFileSizeMatcher(fi, -1, 4)
	expect.True(m.Matches(fileToCheck))

	m = matcher.NewFileSizeMatcher(fi, -1, 15)
	expect.True(m.Matches(fileToCheck))

	m = matcher.NewFileSizeMatcher(fi, 0, -1)
	expect.True(m.Matches(fileToCheck))

	m = matcher.NewFileSizeMatcher(fi, 5, -1)
	expect.False(m.Matches(fileToCheck))

	m = matcher.NewFileSizeMatcher(fi, 3, 4)
	expect.True(m.Matches(fileToCheck))

	m = matcher.NewFileSizeMatcher(fi, 3, 3)
	expect.False(m.Matches(fileToCheck))

	m = matcher.NewFileSizeMatcher(fi, 5, 4)
	expect.False(m.Matches(fileToCheck))

	m = matcher.NewFileSizeMatcher(fi, -1, -1)
	expect.False(m.Matches(fileToCheck))
}

func TestSizeMatcherForFile(t *testing.T) {
	expect := assert.New(t)
	fileToCheck := "../data/fixtures/global/file.txt"
	m := matcher.NewFileSizeMatcher(nil, 3, 4)
	expect.True(m.Matches(fileToCheck))
}

func TestSizeMatcherForDirWithStat(t *testing.T) {
	expect := assert.New(t)
	dirToCheck := "fixtures/"

	mockFs := testhelpers.MockFileSystem(map[string]string{
		dirToCheck:  "",
	})
	di, _ := mockFs.Stat(dirToCheck)

	m := matcher.NewFileSizeMatcher(di, 0, 5)
	expect.False(m.Matches(dirToCheck))

}
