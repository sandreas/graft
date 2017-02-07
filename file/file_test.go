package file_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/file"
	"regexp"
)


func TestWalkPathByPattern(t *testing.T) {
	assert := assert.New(t)

	allPattern, _ := regexp.Compile("(.*)")
	txtPattern, _ := regexp.Compile("(.*)\\.txt")
	allFiles, _ := file.WalkPathByPattern("../data/fixtures/file/WalkPathByPattern", allPattern)
	txtFiles, _ := file.WalkPathByPattern("../data/fixtures/file/WalkPathByPattern", txtPattern)

	assert.Len(allFiles, 9)
	assert.Len(txtFiles, 4)
}

func TestCopyResumed(t *testing.T) {
	assert := assert.New(t)



	assert.True(true)
}

func TestFilesEqualQuick(t *testing.T) {
	assert := assert.New(t)

	file1 := "../data/fixtures/file/AreFilesEqual/equal1.txt"
	file2 := "../data/fixtures/file/AreFilesEqual/equal2.txt"
	file3 := "../data/fixtures/file/AreFilesEqual/not-equal.txt"
	assert.True(file.FilesEqualQuick(file1, file2, 5))
	assert.False(file.FilesEqualQuick(file1, file3, 5))
}
