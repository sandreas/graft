package newfile

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/newpattern"
)

func TestFindFilesBySourcePatternWithFile(t *testing.T) {
	expect := assert.New(t)

	sourcePattern := newpattern.NewSourcePattern("../data/fixtures/global/file.txt")

	foundFiles,_ := FindFilesBySourcePattern(*sourcePattern)
	expect.Equal(1, len(foundFiles))
	expect.Equal("../data/fixtures/global/file.txt", foundFiles["../data/fixtures/global/file.txt"])
}

func TestFindFilesBySourcePatternWithDirectory(t *testing.T) {
	expect := assert.New(t)

	sourcePattern := newpattern.NewSourcePattern("../data/fixtures/file/WalkPathByPattern/dir")

	foundFiles, _ := FindFilesBySourcePattern(*sourcePattern)
	expect.Equal(4, len(foundFiles))
	expect.Equal("../data/fixtures/file/WalkPathByPattern/dir/", foundFiles["../data/fixtures/file/WalkPathByPattern/dir/"])
	expect.Equal("../data/fixtures/file/WalkPathByPattern/dir/dirfile.txt", foundFiles["../data/fixtures/file/WalkPathByPattern/dir/dirfile.txt"])
	expect.Equal("../data/fixtures/file/WalkPathByPattern/dir/subdir/", foundFiles["../data/fixtures/file/WalkPathByPattern/dir/subdir/"])
	expect.Equal("../data/fixtures/file/WalkPathByPattern/dir/subdir/subdirfile.log", foundFiles["../data/fixtures/file/WalkPathByPattern/dir/subdir/subdirfile.log"])
}

func TestFindFilesBySourcePatternWithGlob(t *testing.T) {
	expect := assert.New(t)

	sourcePattern := newpattern.NewSourcePattern("../data/fixtures/file/WalkPathByPattern/dir/*irfile*")

	foundFiles, _ := FindFilesBySourcePattern(*sourcePattern)
	expect.Equal(2, len(foundFiles))
	expect.Equal("../data/fixtures/file/WalkPathByPattern/dir/dirfile.txt", foundFiles["../data/fixtures/file/WalkPathByPattern/dir/dirfile.txt"])
	expect.Equal("../data/fixtures/file/WalkPathByPattern/dir/subdir/subdirfile.log", foundFiles["../data/fixtures/file/WalkPathByPattern/dir/subdir/subdirfile.log"])
}