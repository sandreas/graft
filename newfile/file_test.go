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
	expect.Equal("../data/fixtures/global/file.txt", foundFiles[0])
}

func TestFindFilesBySourcePatternWithDirectory(t *testing.T) {
	expect := assert.New(t)

	sourcePattern := newpattern.NewSourcePattern("../data/fixtures/file/WalkPathByPattern/dir")
	compiledRegex, _ := sourcePattern.Compile()
	matcher := newmatcher.NewRegexMatcher(*compiledRegex)

	foundFiles, _ := FindFilesBySourcePattern(*sourcePattern, matcher)
	expect.Equal(4, len(foundFiles))
	expect.Equal("../data/fixtures/file/WalkPathByPattern/dir/", foundFiles[0])
	expect.Equal("../data/fixtures/file/WalkPathByPattern/dir/dirfile.txt", foundFiles[1])
	expect.Equal("../data/fixtures/file/WalkPathByPattern/dir/subdir/", foundFiles[2])
	expect.Equal("../data/fixtures/file/WalkPathByPattern/dir/subdir/subdirfile.log", foundFiles[3])
}

func TestFindFilesBySourcePatternWithGlob(t *testing.T) {
	expect := assert.New(t)

	sourcePattern := newpattern.NewSourcePattern("../data/fixtures/file/WalkPathByPattern/dir/*irfile*")
	compiledRegex, _ := sourcePattern.Compile()
	matcher := newmatcher.NewRegexMatcher(*compiledRegex)

	foundFiles, _ := FindFilesBySourcePattern(*sourcePattern, matcher)
	expect.Equal(2, len(foundFiles))
	expect.Equal("../data/fixtures/file/WalkPathByPattern/dir/dirfile.txt", foundFiles[0])
	expect.Equal("../data/fixtures/file/WalkPathByPattern/dir/subdir/subdirfile.log", foundFiles[1])
}