package file_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/file"
	"regexp"
	"os"
	"github.com/sandreas/graft/pattern"
	"path/filepath"
	"sort"
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
	rarPattern, _ := regexp.Compile(pattern.GlobToRegex("*.rar"))



	allFiles, _ := file.WalkPathByPattern("../data/fixtures/file/WalkPathByPattern", allPattern, progressHandlerWalkPathByPattern)
	txtFiles, _ := file.WalkPathByPattern("../data/fixtures/file/WalkPathByPattern", txtPattern, progressHandlerWalkPathByPattern)
	rarFiles, _ := file.WalkPathByPattern("../data/fixtures/file/WalkPathByPattern", rarPattern, progressHandlerWalkPathByPattern)
	expect.Len(allFiles, 11)
	expect.Len(txtFiles, 4)
	expect.Len(rarFiles, 3)
	expect.Equal("../data/fixtures/file/WalkPathByPattern/test.part1.rar", filepath.ToSlash(rarFiles[0]))
	expect.Equal("../data/fixtures/file/WalkPathByPattern/test.part2.rar", filepath.ToSlash(rarFiles[1]))
	expect.Equal("../data/fixtures/file/WalkPathByPattern/test.part3.rar", filepath.ToSlash(rarFiles[2]))
	//

	wd, _ := os.Getwd()
	os.Chdir("../data/fixtures/file/WalkPathByPattern")
	rarFiles, _ = file.WalkPathByPattern("", rarPattern, progressHandlerWalkPathByPattern)
	os.Chdir(wd)

	expect.Len(rarFiles, 3)
	expect.Equal("test.part1.rar", rarFiles[0])
	expect.Equal("test.part2.rar", rarFiles[1])
	expect.Equal("test.part3.rar", rarFiles[2])
}

func TestWalkPathFiltered(t *testing.T) {
	expect := assert.New(t)

	rarFiles, _ := file.WalkPathFiltered("../data/fixtures/file/WalkPathByPattern", WalkPathFilteredFilterFunc, progressHandlerWalkPathByPattern)

	expect.Len(rarFiles, 3)
	expect.Equal("../data/fixtures/file/WalkPathByPattern/test.part1.rar", filepath.ToSlash(rarFiles[0].Path))
	expect.Equal("../data/fixtures/file/WalkPathByPattern/test.part2.rar", filepath.ToSlash(rarFiles[1].Path))
	expect.Equal("../data/fixtures/file/WalkPathByPattern/test.part3.rar", filepath.ToSlash(rarFiles[2].Path))

}

func WalkPathFilteredFilterFunc(f file.File, err error) bool {
	if err == nil {
		rarPattern, _ := regexp.Compile(pattern.GlobToRegex("*.rar"))
		return rarPattern.MatchString(f.Path)
	}
	return false
}

func progressHandlerWalkPathByPattern(entriesWalked, entriesMatched int64, finished bool) (int64) {
	if(finished) {
		entriesWalked = entriesMatched
	}
	return 5
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


func TestReadAllLinesFunc(t *testing.T) {
	expect := assert.New(t)


	fileWithEmptySkipped, _ := file.ReadAllLinesFunc("../data/fixtures/file/ReadAllLines/10-with-empty.txt", file.SkipEmptyLines)
	expect.Len(fileWithEmptySkipped, 7)
	expect.Equal(fileWithEmptySkipped[0], "line 1")
	expect.Equal(fileWithEmptySkipped[5], "line 7")
	expect.Equal(fileWithEmptySkipped[6], "line 10")
}

func TestIsFile(t *testing.T) {
	expect := assert.New(t)

	isFile, _, err := file.IsFile("../data/fixtures/global/not-exists.txt")
	expect.False(isFile)
	expect.NotNil(err)

	isFile, _, err = file.IsFile("../data/fixtures/global/textfile.txt")
	expect.True(isFile)
	expect.Nil(err)

	isFile, _, err = file.IsFile("../data/fixtures/global/")
	expect.False(isFile)
	expect.Nil(err)

}

func TestMakePathMap(t *testing.T){
	expect := assert.New(t)

	var matchingPaths []string
	matchingPaths = append(matchingPaths, "graft.go")
	matchingPaths = append(matchingPaths, "LICENSE")
	matchingPaths = append(matchingPaths, "README.md")
	matchingPaths = append(matchingPaths, "data/fixtures/global/file.txt")
	matchingPaths = append(matchingPaths, "data/fixtures")
	matchingPaths = append(matchingPaths, "data")

	pathMap := file.MakePathMap(matchingPaths)
	expected := []string{ "data", "graft.go", "LICENSE", "README.md"}
	sort.Strings(expected)
	expect.Equal(pathMap["/"], expected)

	expected = []string{"data/fixtures"}
	expect.Equal(pathMap["/data"], expected)


	expected = []string{"data/fixtures/global"}
	expect.Equal(pathMap["/data/fixtures"], expected)

	expected = []string{"data/fixtures/global/file.txt"}
	expect.Equal(pathMap["/data/fixtures/global"], expected)

}

//func TestMkdirAll(t *testing.T) {
//	expect := assert.New(t)
//	srcStat, _ := os.Stat("../data")
//	dst := "../data/tmp/mkdirall/recursive/directory"
//
//	if file.Exists(dst) {
//		os.RemoveAll("../data/tmp/mkdirall")
//	}
//
//	expect.False(file.Exists(dst))
//
//	file.MkdirAll(dst, srcStat.Mode())
//	dstStat, _ := os.Stat(dst)
//
//	expect.True(dstStat.IsDir())
//	expect.Equal(srcStat.Mode(), dstStat.Mode())
//}

//func TestHumanReadableSize(t *testing.T) {
//	expect := assert.New(t)
//
//
//	expect.Equal("0 B", file.HumanReadableSize(0, true))
//	//SI     BINARY
//	//
//	//0:        0 B        0 B
//	//27:       27 B       27 B
//	//999:      999 B      999 B
//	//1000:     1.0 kB     1000 B
//	//1023:     1.0 kB     1023 B
//	//1024:     1.0 kB    1.0 KiB
//	//1728:     1.7 kB    1.7 KiB
//	//110592:   110.6 kB  108.0 KiB
//	//7077888:     7.1 MB    6.8 MiB
//	//452984832:   453.0 MB  432.0 MiB
//	//28991029248:    29.0 GB   27.0 GiB
//	//1855425871872:     1.9 TB    1.7 TiB
//	//9223372036854775807:     9.2 EB    8.0 EiB   (Long.MAX_VALUE)
//}