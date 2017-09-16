package sftpd_test

import (
	"testing"

	"os"
	"path/filepath"

	"github.com/sandreas/graft/sftpd"
	"github.com/stretchr/testify/assert"
)

var files = []string{
	"data/fixtures/global/",
	"data/fixtures/global/dir/",
	"data/fixtures/global/dir/dirfile.txt",
	"data/fixtures/global/dir/subdir/",
	"data/fixtures/global/dir/subdir/subdirfile.log",
	"data/fixtures/global/documents (2010)/",
	"data/fixtures/global/documents (2010)/document (2010).txt",
	"data/fixtures/global/file.txt",
	"data/fixtures/global/textfile.txt",
}

var filesWithDot = []string{
	"../data/fixtures/global/",
	"../data/fixtures/global/dir/",
	"../data/fixtures/global/dir/dirfile.txt",
	"../data/fixtures/global/dir/subdir/",
	"../data/fixtures/global/dir/subdir/subdirfile.log",
	"../data/fixtures/global/documents (2010)/",
	"../data/fixtures/global/documents (2010)/document (2010).txt",
	"../data/fixtures/global/file.txt",
	"../data/fixtures/global/textfile.txt",
}

var filesOnly = []string{
	"data/fixtures/file/CopyResumed/test2-dst-larger.txt",
	"data/fixtures/file/CopyResumed/test3-dst-partial.txt",
	"data/fixtures/file/CopyResumed/test4-dst-exists.txt",
}

var filesMulti = []string{
	"data/fixtures/file/",
	"data/fixtures/file/AreFilesEqual/",
	"data/fixtures/file/AreFilesEqual/equal1.txt",
	"data/fixtures/file/AreFilesEqual/equal2.txt",
	"data/fixtures/file/AreFilesEqual/not-equal.txt",
	"data/fixtures/file/CopyResumed/",
	"data/fixtures/file/CopyResumed/test1-src.txt",
	"data/fixtures/file/CopyResumed/test2-dst-larger.txt",
	"data/fixtures/file/CopyResumed/test2-src.txt",
	"data/fixtures/file/CopyResumed/test3-dst-partial.txt",
	"data/fixtures/file/CopyResumed/test3-src.txt",
	"data/fixtures/file/CopyResumed/test4-dst-exists.txt",
	"data/fixtures/file/CopyResumed/test4-src.txt",
	"data/fixtures/file/ReadAllLines/",
	"data/fixtures/file/ReadAllLines/10-lines.txt",
	"data/fixtures/file/ReadAllLines/10-with-empty.txt",
	"data/fixtures/file/WalkPathByPattern/",
	"data/fixtures/file/WalkPathByPattern/dir/",
	"data/fixtures/file/WalkPathByPattern/dir/dirfile.txt",
	"data/fixtures/file/WalkPathByPattern/dir/subdir/",
	"data/fixtures/file/WalkPathByPattern/dir/subdir/subdirfile.log",
	"data/fixtures/file/WalkPathByPattern/documents (2010)/",
	"data/fixtures/file/WalkPathByPattern/documents (2010)/document (2010).txt",
	"data/fixtures/file/WalkPathByPattern/file.txt",
	"data/fixtures/file/WalkPathByPattern/test.part1.rar",
	"data/fixtures/file/WalkPathByPattern/test.part2.rar",
	"data/fixtures/file/WalkPathByPattern/test.part3.rar",
	"data/fixtures/file/WalkPathByPattern/textfile.txt",
	"data/fixtures/global/",
	"data/fixtures/global/dir/",
	"data/fixtures/global/dir/dirfile.txt",
	"data/fixtures/global/dir/subdir/",
	"data/fixtures/global/dir/subdir/subdirfile.log",
	"data/fixtures/global/documents (2010)/",
	"data/fixtures/global/documents (2010)/document (2010).txt",
	"data/fixtures/global/file.txt",
	"data/fixtures/global/textfile.txt",
}

var filesSpecial = []string{
	"data/fixtures/file/AreFilesEqual/equal2.txt",
	"data/fixtures/file/WalkPathByPattern/documents (2010)/document (2010).txt",
	"data/fixtures/global/documents (2010)/document (2010).txt",
}

var singleFile = []string{
	"./test.txt",
}

var singleFileServe = []string{
	"../data/fixtures/file/AreFilesEqual/equal2.txt",
}

func TestFiles(t *testing.T) {
	expect := assert.New(t)

	mapper := sftpd.NewPathMapper(files, "data/fixtures")

	result, ok := mapper.List("global")
	want := []string{"/global/dir", "/global/documents (2010)", "/global/file.txt", "/global/textfile.txt"}
	expect.True(ok)
	expect.Equal(want, result)

	resultWithLeadingLash, ok2 := mapper.List("/global")
	expect.True(ok2)
	expect.Equal(want, resultWithLeadingLash)
}

func TestFilesDoubleDotSlash(t *testing.T) {
	expect := assert.New(t)

	mapper := sftpd.NewPathMapper(filesWithDot, "../data/fixtures")

	result, ok := mapper.List("global")
	want := []string{"/global/dir", "/global/documents (2010)", "/global/file.txt", "/global/textfile.txt"}
	expect.True(ok)
	expect.Equal(want, result)

	resultWithLeadingLash, ok2 := mapper.List("/global")
	expect.True(ok2)
	expect.Equal(want, resultWithLeadingLash)
}

func TestDotSlash(t *testing.T) {
	expect := assert.New(t)

	mapper := sftpd.NewPathMapper(files, "./data/fixtures")

	result, ok := mapper.List("global")
	want := []string{"/global/dir", "/global/documents (2010)", "/global/file.txt", "/global/textfile.txt"}
	expect.True(ok)
	expect.Equal(want, result)

	resultWithLeadingLash, ok2 := mapper.List("/global")
	expect.True(ok2)
	expect.Equal(want, resultWithLeadingLash)
}

func TestFilesOnly(t *testing.T) {
	expect := assert.New(t)

	mapper := sftpd.NewPathMapper(filesOnly, "data/fixtures")
	result, ok := mapper.List(string(os.PathSeparator))
	want := []string{"/file"}
	expect.True(ok)
	expect.Equal(want, result)

}

func TestFilesMulti(t *testing.T) {
	expect := assert.New(t)

	mapper := sftpd.NewPathMapper(filesMulti, "data/fixtures")
	result, ok := mapper.List("/global")
	want := []string{"/global/dir", "/global/documents (2010)", "/global/file.txt", "/global/textfile.txt"}
	expect.True(ok)
	expect.Equal(want, result)

}

func TestFilesSpecial(t *testing.T) {
	expect := assert.New(t)

	mapper := sftpd.NewPathMapper(filesSpecial, "data/fixtures")
	result, ok := mapper.List("/")
	want := []string{"/file", "/global"}
	expect.True(ok)
	expect.Equal(want, result)

}

func TestPathTo(t *testing.T) {
	expect := assert.New(t)

	mapper := sftpd.NewPathMapper(files, "./data/fixtures")
	path, err := mapper.PathTo("global")
	expect.Equal(filepath.FromSlash("data/fixtures/global"), path)
	expect.Nil(err)

	path, err = mapper.PathTo("non-existing")
	expect.Equal("", path)
	expect.Error(err)

	mapper = sftpd.NewPathMapper(filesWithDot, "../data/fixtures")
	path, err = mapper.PathTo("global")
	expect.Equal(filepath.FromSlash("../data/fixtures/global"), path)
	expect.Nil(err)
}

func TestPathToSingleFile(t *testing.T) {
	expect := assert.New(t)

	mapper := sftpd.NewPathMapper(singleFile, ".")
	result, ok := mapper.PathTo("/test.txt")
	want := "test.txt"
	expect.NoError(ok)
	expect.Equal(want, result)
}

func TestPathToSingleFileDirectory(t *testing.T) {
	expect := assert.New(t)

	mapper := sftpd.NewPathMapper(singleFile, ".")
	result, ok := mapper.PathTo("/")
	want := "."
	expect.NoError(ok)
	expect.Equal(want, result)
}

func TestStat(t *testing.T) {
	expect := assert.New(t)

	mapper := sftpd.NewPathMapper(filesWithDot, "../data/fixtures")
	stat, err := mapper.Stat("global")
	expect.Nil(err)
	expect.True(stat.IsDir())

	stat, err = mapper.Stat("/")
	expect.Nil(err)
	expect.True(stat.IsDir())

	stat, err = mapper.Stat("non-existing")
	expect.Error(err)
	expect.Nil(stat)
}

func TestStatForServeFileOnly(t *testing.T) {
	expect := assert.New(t)

	mapper := sftpd.NewPathMapper(singleFileServe, "../data/fixtures/file/AreFilesEqual")
	stat, err := mapper.Stat("/")
	expect.Nil(err)
	expect.True(stat.IsDir())

	stat, err = mapper.Stat("/equal2.txt")
	expect.Nil(err)
	expect.False(stat.IsDir())
}
