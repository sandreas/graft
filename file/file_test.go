package file_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/file"
	"regexp"
	"os"
)

func TestExists(t *testing.T) {
	expect := assert.New(t)
	expect.True(file.Exists("../data/fixtures/global/file.txt"))
	expect.False(file.Exists("../data/fixtures/global/not-exists.txt"))
}

func TestWalkPathByPattern(t *testing.T) {
	expect := assert.New(t)

	allPattern, _ := regexp.Compile("(.*)")
	txtPattern, _ := regexp.Compile("(.*)\\.txt")
	allFiles, _ := file.WalkPathByPattern("../data/fixtures/file/WalkPathByPattern", allPattern)
	txtFiles, _ := file.WalkPathByPattern("../data/fixtures/file/WalkPathByPattern", txtPattern)

	expect.Len(allFiles, 9)
	expect.Len(txtFiles, 4)
}

func TestCopy(t *testing.T) {
	expect := assert.New(t)

	srcStr := "../data/fixtures/global/textfile.txt"
	dstStr := "../data/tmp/dst.txt"

	if _, err := os.Stat(dstStr); os.IsExist(err) {
		os.Remove(dstStr)
	}
	file.Copy(srcStr, dstStr)
	equal, _ := file.ContentsEqual(srcStr, dstStr)
	expect.True(equal)
}

func TestReplace(t *testing.T) {
	expect := assert.New(t)

	srcStr := "../data/fixtures/global/file.txt"
	dstStr := "../data/tmp/dst.txt"

	file.Replace(srcStr, dstStr)
	equal1, _ := file.ContentsEqual(srcStr, dstStr);
	expect.True(equal1)

	srcStr = "../data/fixtures/global/textfile.txt"
	file.Replace(srcStr, dstStr)
	equal2, _ := file.ContentsEqual(srcStr, dstStr)

	expect.True(equal2)
}

func TestCopyResumed(t *testing.T) {
	expect := assert.New(t)

	src, dst := prepareCopyResumed("test1-src.txt", "test1-dst-not-exists.txt")
	file.CopyResumed(src, dst, func(bytesTransferred, size, chunkSize int64) (int64) {
		return chunkSize
	})
	expect.True(file.ContentsEqual((src).Name(), (dst).Name()))

	src, dst = prepareCopyResumed("test2-src.txt", "test2-dst-larger.txt")
	err := file.CopyResumed(src, dst, func(bytesTransferred, size, chunkSize int64) (int64) {
		return chunkSize
	})
	expect.Error(err)


	src, dst = prepareCopyResumed("test3-src.txt", "test3-dst-partial.txt")
	file.CopyResumed(src, dst, func(bytesTransferred, size, chunkSize int64) (int64) {
		return chunkSize
	})
	expect.True(file.ContentsEqual((src).Name(), (dst).Name()))

	src, dst = prepareCopyResumed("test4-src.txt", "test4-dst-exists.txt")
	file.CopyResumed(src, dst, func(bytesTransferred, size, chunkSize int64) (int64) {
		return chunkSize
	})
	expect.True(file.ContentsEqual((src).Name(), (dst).Name()))
}

func prepareCopyResumed(srcName, dstName string) (*os.File, *os.File) {
	srcFixture := "../data/fixtures/file/CopyResumed/" + srcName
	src := "../data/tmp/" + srcName

	dstFixture := "../data/fixtures/file/CopyResumed/" + dstName
	dst := "../data/tmp/" + dstName

	if file.Exists(dst) {
		os.Remove(dst)
	}

	if file.Exists(dstFixture) {
		file.Replace(dstFixture, dst)
	}

	file.Replace(srcFixture, src)

	stat, _ := os.Stat(src)

	srcPointer, _ := os.Open(src)
	dstPointer, _ := os.OpenFile(dst, os.O_RDWR | os.O_CREATE, stat.Mode())

	return srcPointer, dstPointer
}


func TestContentsEqual(t *testing.T) {
	expect := assert.New(t)
	file1 := "../data/fixtures/file/AreFilesEqual/equal1.txt"
	file2 := "../data/fixtures/file/AreFilesEqual/equal2.txt"
	file3 := "../data/fixtures/file/AreFilesEqual/not-equal.txt"
	equal12, _ := file.ContentsEqual(file1, file2)
	equal13, _ := file.ContentsEqual(file1, file3)

	expect.True(equal12)
	expect.False(equal13)
}

func TestFileContentsEqual(t *testing.T) {
	expect := assert.New(t)
	file1, _ := os.Open("../data/fixtures/file/AreFilesEqual/equal1.txt")
	file2, _ := os.Open("../data/fixtures/file/AreFilesEqual/equal2.txt")
	file3, _ := os.Open("../data/fixtures/file/AreFilesEqual/not-equal.txt")
	equal12, _ := file.FileContentsEqual(file1, file2)
	equal13, _ := file.FileContentsEqual(file1, file3)

	expect.True(equal12)
	expect.False(equal13)
}

func TestFilesEqualQuick(t *testing.T) {
	expect := assert.New(t)

	file1, _ := os.Open("../data/fixtures/file/AreFilesEqual/equal1.txt")
	file2, _ := os.Open("../data/fixtures/file/AreFilesEqual/equal2.txt")
	file3, _ := os.Open("../data/fixtures/file/AreFilesEqual/not-equal.txt")
	expect.True(file.FileContentsEqualQuick(file1, file2, 5))
	expect.False(file.FileContentsEqualQuick(file1, file3, 5))
}
