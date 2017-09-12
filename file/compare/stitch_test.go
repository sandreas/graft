package compare_test

import (
	"testing"

	"github.com/sandreas/graft/file/compare"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
)

func prepareFileSystem() afero.Fs {
	appFS := afero.NewMemMapFs()
	appFS.Mkdir("src", 0644)
	appFS.Mkdir("dst", 0644)

	// not existing
	afero.WriteFile(appFS, "file1-src.txt", []byte("0123456789012345678901234567890123456789"), 0755)

	// not resumable - destination bigger
	afero.WriteFile(appFS, "file2-src.txt", []byte("0123456789012345678901234567890123456789"), 0755)
	afero.WriteFile(appFS, "file2-dst.txt", []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), 0755)

	//// already completed
	afero.WriteFile(appFS, "file3-src.txt", []byte("0123456789012345678901234567890123456789"), 0755)
	afero.WriteFile(appFS, "file3-dst.txt", []byte("0123456789012345678901234567890123456789"), 0755)
	//
	// resumable
	afero.WriteFile(appFS, "file4-src.txt", []byte("0123456789012345678901234567890123456789"), 0755)
	afero.WriteFile(appFS, "file4-dst.txt", []byte("0123456789012345"), 0755)

	// not resumable - diff at start
	afero.WriteFile(appFS, "file5-src.txt", []byte("0123456789012345678901234567890123456789"), 0755)
	afero.WriteFile(appFS, "file5-dst.txt", []byte("a123456789012345678901234567"), 0755)
	//
	// not resumable - diff at end
	afero.WriteFile(appFS, "file6-src.txt", []byte("0123456789012345678901234567890123456789"), 0755)
	afero.WriteFile(appFS, "file6-dst.txt", []byte("0123456789012345678a"), 0755)
	//
	// not resumable - diff at middle
	afero.WriteFile(appFS, "file7-src.txt", []byte("0123456789012345678901234567890123456789"), 0755)
	afero.WriteFile(appFS, "file7-dst.txt", []byte("01234567890123aaaaaaaa234567890123"), 0755)

	return appFS
}

func prepareTestSubect(fileNamePrefix string) (*compare.Stitch, error) {
	fs := prepareFileSystem()
	src, _ := fs.Open(fileNamePrefix+"-src.txt")
	dst, _ := fs.OpenFile(fileNamePrefix+"-dst.txt", os.O_RDWR|os.O_CREATE, 0755)

	return compare.NewStich(src, dst, 2)
}

func TestNotExistingFile(t *testing.T) {
	expect := assert.New(t)

	subject, err := prepareTestSubect("file1")

	expect.NoError(err)
	expect.NotNil(subject)
	expect.False(subject.IsComplete())
}

func TestBiggerDestinationFile(t *testing.T) {
	expect := assert.New(t)

	subject, err := prepareTestSubect("file2")

	expect.Error(err)
	expect.Equal("source is smaller than destination", err.Error())
	expect.Nil(subject)
}

func TestAlreadyCompletedFile(t *testing.T) {
	expect := assert.New(t)

	subject, err := prepareTestSubect("file3")

	expect.NoError(err)
	expect.NotNil(subject)
	expect.True(subject.IsComplete())
}

func TestPartialFileCanBeResumed(t *testing.T) {
	expect := assert.New(t)

	subject, err := prepareTestSubect("file4")

	expect.NoError(err)
	expect.NotNil(subject)
	expect.False(subject.IsComplete())
}

func TestFileDiffAtStartCannotBeResumed(t *testing.T) {
	expect := assert.New(t)

	subject, err := prepareTestSubect("file5")

	expect.Error(err)
	expect.Equal("source file does not match destination file", err.Error())
	expect.Nil(subject)
}

func TestFileDiffAtEndCannotBeResumed(t *testing.T) {
	expect := assert.New(t)

	subject, err := prepareTestSubect("file6")

	expect.Error(err)
	expect.Equal("source file does not match destination file", err.Error())
	expect.Nil(subject)
}

func TestFileDiffAtMiddleCannotBeResumed(t *testing.T) {
	expect := assert.New(t)

	subject, err := prepareTestSubect("file7")

	expect.Error(err)
	expect.Equal("source file does not match destination file", err.Error())
	expect.Nil(subject)
}
