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

	fi, _ := mockFs.Stat(fileToCheck)

	subject := matcher.NewFileAgeMatcher(fi, time.Time{}, time2016)
	expect.True(subject.Matches(fileToCheck))

	subject = matcher.NewFileAgeMatcher(fi, time2014, time.Time{})
	expect.True(subject.Matches(fileToCheck))

	subject = matcher.NewFileAgeMatcher(fi, time2014, time2016)
	expect.True(subject.Matches(fileToCheck))

	subject = matcher.NewFileAgeMatcher(fi, time.Time{}, time.Time{})
	expect.True(subject.Matches(fileToCheck))
}

func TestFileAgeMatcherWithoutStat(t *testing.T) {
	expect := assert.New(t)
	fileToCheck := "../data/fixtures/global/file.txt"
	m := matcher.NewFileAgeMatcher(nil, time.Time{}, time.Now())
	expect.True(m.Matches(fileToCheck))
}

func TestFileNotExists(t *testing.T) {
	expect := assert.New(t)
	fileToCheck := "not-exists.txt"
	m := matcher.NewFileAgeMatcher(nil, time.Time{}, time.Now())
	expect.False(m.Matches(fileToCheck))
}

