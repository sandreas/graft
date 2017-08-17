package matcher_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/matcher"
	"time"
	"github.com/sandreas/graft/testhelpers"
)

// time.Time{}

func TestFileAgeMatcher(t *testing.T) {
	expect := assert.New(t)

	time2014 := time.Date(2014, 1, 1, 1, 1, 1, 1, time.Local)
	modTime := time.Date(2015, 1, 1, 1, 1, 1, 1, time.Local)
	time2016 := time.Date(2016, 1, 1, 1, 1, 1, 1, time.Local)
	fileToCheck := "file.txt"
	mockFs := testhelpers.MockFileSystem(map[string]string{
		fileToCheck: "",
	})
	mockFs.Chtimes(fileToCheck, modTime, modTime)


	subject := matcher.NewFileAgeMatcher( time.Time{}, time2016)
	subject.Fs = mockFs
	expect.True(subject.Matches(fileToCheck))

	subject = matcher.NewFileAgeMatcher( time2014, time.Time{})
	subject.Fs = mockFs
	expect.True(subject.Matches(fileToCheck))

	subject = matcher.NewFileAgeMatcher( time2014, time2016)
	subject.Fs = mockFs
	expect.True(subject.Matches(fileToCheck))

	subject = matcher.NewFileAgeMatcher( time.Time{}, time.Time{})
	subject.Fs = mockFs
	expect.True(subject.Matches(fileToCheck))
}

func TestFileAgeMatcherWithoutStat(t *testing.T) {
	expect := assert.New(t)
	fileToCheck := "../data/fixtures/global/file.txt"
	m := matcher.NewFileAgeMatcher(time.Time{}, time.Now())
	expect.True(m.Matches(fileToCheck))
}

func TestFileNotExists(t *testing.T) {
	expect := assert.New(t)
	fileToCheck := "not-exists.txt"
	m := matcher.NewFileAgeMatcher( time.Time{}, time.Now())
	expect.False(m.Matches(fileToCheck))
}

