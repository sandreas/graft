package file_test

import (
	"testing"

	"github.com/sandreas/graft/file"
	"github.com/stretchr/testify/assert"
)

func TestFileTreeIntegrateAndList(t *testing.T) {
	expect := assert.New(t)

	fileTree := file.NewTreeNode("fixtures")
	expect.NoError(fileTree.Integrate("fixtures/subdir/file.txt"))
	expect.NoError(fileTree.Integrate("fixtures/other-subdir/file.txt"))
	list := fileTree.List("fixtures")
	expect.Len(list, 2)
}
