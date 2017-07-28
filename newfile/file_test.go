package newfile

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/newpattern"
	"github.com/sandreas/graft/newmatcher"
)

func TestFindFilesBySourcePatternWithFile(t *testing.T) {
	expect := assert.New(t)

	sourcePattern := newpattern.NewSourcePattern("../data/fixtures/global/file.txt")
	compiledRegex, _ := sourcePattern.Compile()
	matcher := newmatcher.NewRegexMatcher(*compiledRegex)

	foundFiles,_ := FindFilesBySourcePattern(*sourcePattern, matcher)
	expect.Equal(1, len(foundFiles))
	expect.Equal("../data/fixtures/global/file.txt", foundFiles["../data/fixtures/global/file.txt"])
}

func TestFindFilesBySourcePatternWithDirectory(t *testing.T) {
	expect := assert.New(t)

	sourcePattern := newpattern.NewSourcePattern("../data/fixtures/file/WalkPathByPattern/dir")
	compiledRegex, _ := sourcePattern.Compile()
	matcher := newmatcher.NewRegexMatcher(*compiledRegex)

	foundFiles, _ := FindFilesBySourcePattern(*sourcePattern, matcher)
	expect.Equal(4, len(foundFiles))
	expect.Equal("../data/fixtures/file/WalkPathByPattern/dir/", foundFiles["../data/fixtures/file/WalkPathByPattern/dir/"])
	expect.Equal("../data/fixtures/file/WalkPathByPattern/dir/dirfile.txt", foundFiles["../data/fixtures/file/WalkPathByPattern/dir/dirfile.txt"])
	expect.Equal("../data/fixtures/file/WalkPathByPattern/dir/subdir/", foundFiles["../data/fixtures/file/WalkPathByPattern/dir/subdir/"])
	expect.Equal("../data/fixtures/file/WalkPathByPattern/dir/subdir/subdirfile.log", foundFiles["../data/fixtures/file/WalkPathByPattern/dir/subdir/subdirfile.log"])
}

func TestFindFilesBySourcePatternWithGlob(t *testing.T) {
	expect := assert.New(t)

	sourcePattern := newpattern.NewSourcePattern("../data/fixtures/file/WalkPathByPattern/dir/*irfile*")
	compiledRegex, _ := sourcePattern.Compile()
	matcher := newmatcher.NewRegexMatcher(*compiledRegex)

	foundFiles, _ := FindFilesBySourcePattern(*sourcePattern, matcher)
	expect.Equal(2, len(foundFiles))
	expect.Equal("../data/fixtures/file/WalkPathByPattern/dir/dirfile.txt", foundFiles["../data/fixtures/file/WalkPathByPattern/dir/dirfile.txt"])
	expect.Equal("../data/fixtures/file/WalkPathByPattern/dir/subdir/subdirfile.log", foundFiles["../data/fixtures/file/WalkPathByPattern/dir/subdir/subdirfile.log"])
}