package newfile

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/newpattern"
	"github.com/sandreas/graft/newmatcher"
	"github.com/sandreas/graft/newprogress"
)

type FakeProgressHandler struct {
	newprogress.ProgressHandler
	increaseItemsCalls int64
	increaseMatchesCalls int64
	finishCalls int64
}

func(ph *FakeProgressHandler) IncreaseItems(){
	ph.increaseItemsCalls++
}

func(ph *FakeProgressHandler) IncreaseMatches(){
	ph.increaseMatchesCalls++
}

func(ph *FakeProgressHandler) Finish(){
	ph.finishCalls++
}

func TestFindFilesBySourcePatternWithFile(t *testing.T) {
	expect := assert.New(t)


	sourcePattern := newpattern.NewSourcePattern("../data/fixtures/global/file.txt")
	compiledRegex, _ := sourcePattern.Compile()
	matcher := newmatcher.NewRegexMatcher(*compiledRegex)

	progressHandler := &FakeProgressHandler{}

	foundFiles,_ := FindFilesBySourcePattern(*sourcePattern, matcher, progressHandler)
	expect.Equal(1, len(foundFiles))
	expect.Equal("../data/fixtures/global/file.txt", foundFiles[0])

	expect.Equal(int64(0), progressHandler.increaseItemsCalls)
	expect.Equal(int64(1), progressHandler.increaseMatchesCalls)
	expect.Equal(int64(1), progressHandler.finishCalls)
}

func TestFindFilesBySourcePatternWithDirectory(t *testing.T) {
	expect := assert.New(t)

	sourcePattern := newpattern.NewSourcePattern("../data/fixtures/file/WalkPathByPattern/dir")
	compiledRegex, _ := sourcePattern.Compile()
	matcher := newmatcher.NewRegexMatcher(*compiledRegex)

	progressHandler := &FakeProgressHandler{}
	foundFiles, _ := FindFilesBySourcePattern(*sourcePattern, matcher, progressHandler)
	expect.Equal(4, len(foundFiles))
	expect.Equal("../data/fixtures/file/WalkPathByPattern/dir/", foundFiles[0])
	expect.Equal("../data/fixtures/file/WalkPathByPattern/dir/dirfile.txt", foundFiles[1])
	expect.Equal("../data/fixtures/file/WalkPathByPattern/dir/subdir/", foundFiles[2])
	expect.Equal("../data/fixtures/file/WalkPathByPattern/dir/subdir/subdirfile.log", foundFiles[3])

	expect.Equal(int64(0), progressHandler.increaseItemsCalls)
	expect.Equal(int64(4), progressHandler.increaseMatchesCalls)
	expect.Equal(int64(1), progressHandler.finishCalls)

}

func TestFindFilesBySourcePatternWithGlob(t *testing.T) {
	expect := assert.New(t)

	sourcePattern := newpattern.NewSourcePattern("../data/fixtures/file/WalkPathByPattern/dir/*irfile*")
	compiledRegex, _ := sourcePattern.Compile()
	matcher := newmatcher.NewRegexMatcher(*compiledRegex)

	progressHandler := &FakeProgressHandler{}
	foundFiles, _ := FindFilesBySourcePattern(*sourcePattern, matcher, progressHandler)
	expect.Equal(2, len(foundFiles))
	expect.Equal("../data/fixtures/file/WalkPathByPattern/dir/dirfile.txt", foundFiles[0])
	expect.Equal("../data/fixtures/file/WalkPathByPattern/dir/subdir/subdirfile.log", foundFiles[1])

	expect.Equal(int64(2), progressHandler.increaseItemsCalls)
	expect.Equal(int64(2), progressHandler.increaseMatchesCalls)
	expect.Equal(int64(1), progressHandler.finishCalls)
}