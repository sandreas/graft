package file_test

import (
	"os"
	"testing"

	"github.com/sandreas/graft/designpattern/observer"
	"github.com/sandreas/graft/file"
	"github.com/sandreas/graft/matcher"
	"github.com/sandreas/graft/pattern"
	"github.com/sandreas/graft/testhelpers"
	"github.com/stretchr/testify/assert"
)

const sep = string(os.PathSeparator)

type FakeObserver struct {
	designpattern.ObserverInterface
	increaseItemsCalls   int64
	increaseMatchesCalls int64
	finishCalls          int64
}

func (ph *FakeObserver) Notify(a ...interface{}) {
	if a[0] == file.LocatorIncreaseMatches {
		ph.increaseMatchesCalls++
		return
	}

	if a[0] == file.LocatorIncreaseItems {
		ph.increaseItemsCalls++
		return
	}

	if a[0] == file.LocatorFinish {
		ph.finishCalls++
		return
	}
}

func preparePattern(patternString string) (*file.Locator, *FakeObserver, *matcher.CompositeMatcher) {
	mockFs := testhelpers.MockFileSystem(map[string]string{
		"fixtures" + sep + "global" + sep + "":                                                                        "",
		"fixtures" + sep + "global" + sep + "file.txt":                                                                "",
		"fixtures" + sep + "file" + sep + "WalkPathByPattern" + sep + "dir" + sep + "":                                "",
		"fixtures" + sep + "file" + sep + "WalkPathByPattern" + sep + "dir" + sep + "dirfile.txt":                     "",
		"fixtures" + sep + "file" + sep + "WalkPathByPattern" + sep + "dir" + sep + "subdir" + sep + "":               "",
		"fixtures" + sep + "file" + sep + "WalkPathByPattern" + sep + "dir" + sep + "subdir" + sep + "subdirfile.log": "",
	})

	sourcePattern := pattern.NewSourcePattern(mockFs, patternString)
	compiledRegex, err := sourcePattern.Compile()
	if err != nil  {
		println("error compiling regex: ", err.Error())
	}
	m := matcher.NewRegexMatcher(compiledRegex)
	composite := matcher.NewCompositeMatcher()
	composite.Add(m)
	fakeObserver := &FakeObserver{}

	subject := file.NewLocator(sourcePattern)
	subject.RegisterObserver(fakeObserver)

	return subject, fakeObserver, composite
}

func TestFindWithFile(t *testing.T) {
	expect := assert.New(t)

	subject, fakeObserver, composite := preparePattern("fixtures/global/file.txt")

	subject.Find(composite)

	expect.Equal(1, len(subject.SourceFiles))
	expect.Equal("fixtures"+sep+"global"+sep+"file.txt", subject.SourceFiles[0])

	expect.Equal(int64(0), fakeObserver.increaseItemsCalls)
	expect.Equal(int64(1), fakeObserver.increaseMatchesCalls)
	expect.Equal(int64(1), fakeObserver.finishCalls)
}

func TestFindFilesWithDirectory(t *testing.T) {
	expect := assert.New(t)

	subject, fakeObserver, composite := preparePattern("fixtures/file/WalkPathByPattern/dir")

	subject.Find(composite)

	expect.Equal(4, len(subject.SourceFiles))
	expect.Equal("fixtures"+sep+"file"+sep+"WalkPathByPattern"+sep+"dir"+sep+"", subject.SourceFiles[0])
	expect.Equal("fixtures"+sep+"file"+sep+"WalkPathByPattern"+sep+"dir"+sep+"dirfile.txt", subject.SourceFiles[1])
	expect.Equal("fixtures"+sep+"file"+sep+"WalkPathByPattern"+sep+"dir"+sep+"subdir"+sep+"", subject.SourceFiles[2])
	expect.Equal("fixtures"+sep+"file"+sep+"WalkPathByPattern"+sep+"dir"+sep+"subdir"+sep+"subdirfile.log", subject.SourceFiles[3])

	expect.Equal(int64(0), fakeObserver.increaseItemsCalls)
	expect.Equal(int64(4), fakeObserver.increaseMatchesCalls)
	expect.Equal(int64(1), fakeObserver.finishCalls)
}

func TestFindFilesWithGlob(t *testing.T) {
	expect := assert.New(t)

	subject, fakeObserver, composite := preparePattern("fixtures/file/WalkPathByPattern/dir/*irfile*")

	subject.Find(composite)

	expect.Equal(2, len(subject.SourceFiles))
	expect.Equal("fixtures"+sep+"file"+sep+"WalkPathByPattern"+sep+"dir"+sep+"dirfile.txt", subject.SourceFiles[0])
	expect.Equal("fixtures"+sep+"file"+sep+"WalkPathByPattern"+sep+"dir"+sep+"subdir"+sep+"subdirfile.log", subject.SourceFiles[1])

	expect.Equal(int64(1), fakeObserver.increaseItemsCalls)
	expect.Equal(int64(2), fakeObserver.increaseMatchesCalls)
	expect.Equal(int64(1), fakeObserver.finishCalls)
}

func TestFind(t *testing.T) {
	expect := assert.New(t)

	subject, fakeObserver, composite := preparePattern("fixtures/file/WalkPathByPattern/dir/")

	subject.Find(composite)

	expect.Equal(4, len(subject.SourceFiles))
	expect.Equal("fixtures"+sep+"file"+sep+"WalkPathByPattern"+sep+"dir"+sep+"", subject.SourceFiles[0])
	expect.Equal("fixtures"+sep+"file"+sep+"WalkPathByPattern"+sep+"dir"+sep+"dirfile.txt", subject.SourceFiles[1])
	expect.Equal("fixtures"+sep+"file"+sep+"WalkPathByPattern"+sep+"dir"+sep+"subdir"+sep+"", subject.SourceFiles[2])
	expect.Equal("fixtures"+sep+"file"+sep+"WalkPathByPattern"+sep+"dir"+sep+"subdir"+sep+"subdirfile.log", subject.SourceFiles[3])
	//
	expect.Equal(int64(0), fakeObserver.increaseItemsCalls)
	expect.Equal(int64(4), fakeObserver.increaseMatchesCalls)
	expect.Equal(int64(1), fakeObserver.finishCalls)
}
