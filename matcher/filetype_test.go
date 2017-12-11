package matcher_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/matcher"
	"github.com/sandreas/graft/testhelpers"
)

func TestTypeMatcherForFileWithStat(t *testing.T) {
	expect := assert.New(t)
	fileToCheck := "file.txt"
	dirToCheck := "dir/"

	mockFs := testhelpers.MockFileSystem(map[string]string{
		fileToCheck: "file",
		dirToCheck: "",
	})

	m := matcher.NewFileTypeMatcher(matcher.TypeDirectory)
	m.Fs = mockFs
	expect.False(m.Matches(fileToCheck))
	expect.True(m.Matches(dirToCheck))

	m = matcher.NewFileTypeMatcher(matcher.TypeFile)
	m.Fs = mockFs
	expect.True(m.Matches(fileToCheck))
	expect.False(m.Matches(dirToCheck))
}
