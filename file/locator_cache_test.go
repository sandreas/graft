package file_test

import (
	"github.com/spf13/afero"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/file"
)


func TestLoadNonExisting(t *testing.T) {
	expect := assert.New(t)

	cacheFile := "locator-cache.txt"
	subject := file.NewLocatorCache(cacheFile)
	subject.Fs = prepareFilesystemTest("", "")
	err := subject.Load()
	expect.Error(err)
	expect.Len(subject.Items, 0)
}

func TestLoadExisting(t *testing.T) {
	expect := assert.New(t)

	cacheFile := "locator-cache.txt"
	subject := file.NewLocatorCache(cacheFile)
	subject.Fs = prepareFilesystemTest(cacheFile, "data/file1.txt\ndata/file2.txt\n")
	err := subject.Load()

	expect.Nil(err)
	expect.Len(subject.Items, 2)
}

func TestSaveErrorWhenFileExists(t *testing.T) {
	expect := assert.New(t)

	cacheFile := "locator-cache.txt"
	cacheContent := "original-content"
	subject := file.NewLocatorCache(cacheFile)
	subject.Fs = prepareFilesystemTest(cacheFile, cacheContent)
	err := subject.Save()
	expect.Error(err)
	actual, err := afero.ReadFile(subject.Fs, cacheFile)
	expect.Nil(err)
	expect.Equal("original-content", string(actual))
}

func TestSave(t *testing.T) {
	expect := assert.New(t)

	cacheFile := "locator-cache.txt"
	subject := file.NewLocatorCache(cacheFile)
	subject.Fs = prepareFilesystemTest("", "")


	subject.Items = []string{
		"data/file1.txt",
		"data/file2.txt\n",
		"data/file3.txt\r\n",
		"data/file3.txt",
	}
	err := subject.Save()
	expect.Nil(err)

	actual, err := afero.ReadFile(subject.Fs, cacheFile)
	expect.Nil(err)
	expect.Equal("data/file1.txt\ndata/file2.txt\ndata/file3.txt\ndata/file3.txt\n", string(actual))
}

func prepareFilesystemTest(cacheFile, content string) afero.Fs {
	appFS := afero.NewMemMapFs()
	if cacheFile != "" {
		afero.WriteFile(appFS, cacheFile, []byte(content), 0644)
	}
	return appFS
}