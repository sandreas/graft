package file_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/file"
)

func TestFilesEqualQuick(t *testing.T) {
	assert := assert.New(t)

	file1 := "../data/fixtures/file/AreFilesEqual/equal1.txt"
	file2 := "../data/fixtures/file/AreFilesEqual/equal2.txt"
	file3 := "../data/fixtures/file/AreFilesEqual/not-equal.txt"
	assert.True(file.FilesEqualQuick(file1, file2, 5))
	assert.False(file.FilesEqualQuick(file1, file3, 5))
}
