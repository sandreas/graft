package newfile

import (
	"github.com/sandreas/graft/newdesignpattern/observer"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/newpattern"
	"github.com/sandreas/graft/newmatcher"
)

type FakeObserver struct {
	newdesignpattern.ObserverInterface
	increaseItemsCalls int64
	increaseMatchesCalls int64
	finishCalls int64
}

func(ph *FakeObserver) Notify(a...interface{}){
	if a[0] == LOCATOR_INCREASE_MATCHES {
		ph.increaseMatchesCalls++
		return
	}

	if a[0] == LOCATOR_INCREASE_ITEMS {
		ph.increaseItemsCalls++
		return
	}

	if a[0] == LOCATOR_FINISH {
		ph.finishCalls++
		return
	}
}


func preparePattern(patternString string) (*Locator, *FakeObserver, *newmatcher.CompositeMatcher) {
	sourcePattern := newpattern.NewSourcePattern(patternString)
	compiledRegex, _ := sourcePattern.Compile()
	matcher := newmatcher.NewRegexMatcher(*compiledRegex)
	composite := newmatcher.NewCompositeMatcher()
	composite.Add(matcher)
	fakeObserver := &FakeObserver{}

	subject := NewLocator(*sourcePattern)
	subject.RegisterObserver(fakeObserver)

	return subject, fakeObserver, composite
}

func TestFindWithFile(t *testing.T) {
	expect := assert.New(t)

	subject, fakeObserver, composite := preparePattern("../data/fixtures/global/file.txt")

	subject.Find(composite)

	expect.Equal(1, len(subject.SourceFiles))
	expect.Equal("../data/fixtures/global/file.txt", subject.SourceFiles[0])

	expect.Equal(int64(0), fakeObserver.increaseItemsCalls)
	expect.Equal(int64(1), fakeObserver.increaseMatchesCalls)
	expect.Equal(int64(1), fakeObserver.finishCalls)
}


func TestFindFilesWithDirectory(t *testing.T) {
	expect := assert.New(t)

	subject, fakeObserver, composite := preparePattern("../data/fixtures/file/WalkPathByPattern/dir")

	subject.Find(composite)

	expect.Equal(4, len(subject.SourceFiles))
	expect.Equal("../data/fixtures/file/WalkPathByPattern/dir/", subject.SourceFiles[0])
	expect.Equal("../data/fixtures/file/WalkPathByPattern/dir/dirfile.txt", subject.SourceFiles[1])
	expect.Equal("../data/fixtures/file/WalkPathByPattern/dir/subdir/", subject.SourceFiles[2])
	expect.Equal("../data/fixtures/file/WalkPathByPattern/dir/subdir/subdirfile.log", subject.SourceFiles[3])

	expect.Equal(int64(0), fakeObserver.increaseItemsCalls)
	expect.Equal(int64(4), fakeObserver.increaseMatchesCalls)
	expect.Equal(int64(1), fakeObserver.finishCalls)
}

func TestFindFilesWithGlob(t *testing.T) {
	expect := assert.New(t)

	subject, fakeObserver, composite := preparePattern("../data/fixtures/file/WalkPathByPattern/dir/*irfile*")

	subject.Find(composite)


	expect.Equal(2, len(subject.SourceFiles))
	expect.Equal("../data/fixtures/file/WalkPathByPattern/dir/dirfile.txt", subject.SourceFiles[0])
	expect.Equal("../data/fixtures/file/WalkPathByPattern/dir/subdir/subdirfile.log", subject.SourceFiles[1])

	expect.Equal(int64(2), fakeObserver.increaseItemsCalls)
	expect.Equal(int64(2), fakeObserver.increaseMatchesCalls)
	expect.Equal(int64(1), fakeObserver.finishCalls)
}