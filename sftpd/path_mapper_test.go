package sftpd

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewSourcePattern(t *testing.T) {
	expect := assert.New(t)

	files := []string{
		//"data/fixtures/file/",
		//"data/fixtures/file/AreFilesEqual/",
		//"data/fixtures/file/AreFilesEqual/equal1.txt",
		//"data/fixtures/file/AreFilesEqual/equal2.txt",
		//"data/fixtures/file/AreFilesEqual/not-equal.txt",
		//"data/fixtures/file/CopyResumed/",
		//"data/fixtures/file/CopyResumed/test1-src.txt",
		//"data/fixtures/file/CopyResumed/test2-dst-larger.txt",
		//"data/fixtures/file/CopyResumed/test2-src.txt",
		//"data/fixtures/file/CopyResumed/test3-dst-partial.txt",
		//"data/fixtures/file/CopyResumed/test3-src.txt",
		//"data/fixtures/file/CopyResumed/test4-dst-exists.txt",
		//"data/fixtures/file/CopyResumed/test4-src.txt",
		//"data/fixtures/file/ReadAllLines/",
		//"data/fixtures/file/ReadAllLines/10-lines.txt",
		//"data/fixtures/file/ReadAllLines/10-with-empty.txt",
		//"data/fixtures/file/WalkPathByPattern/",
		//"data/fixtures/file/WalkPathByPattern/dir/",
		//"data/fixtures/file/WalkPathByPattern/dir/dirfile.txt",
		//"data/fixtures/file/WalkPathByPattern/dir/subdir/",
		//"data/fixtures/file/WalkPathByPattern/dir/subdir/subdirfile.log",
		//"data/fixtures/file/WalkPathByPattern/documents (2010)/",
		//"data/fixtures/file/WalkPathByPattern/documents (2010)/document (2010).txt",
		//"data/fixtures/file/WalkPathByPattern/file.txt",
		//"data/fixtures/file/WalkPathByPattern/test.part1.rar",
		//"data/fixtures/file/WalkPathByPattern/test.part2.rar",
		//"data/fixtures/file/WalkPathByPattern/test.part3.rar",
		//"data/fixtures/file/WalkPathByPattern/textfile.txt",
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

	mapper := NewPathMapper(files, "data/fixtures")

	result, ok := mapper.Get("global")
	want := []string{"/global/dir", "/global/documents (2010)", "/global/file.txt", "/global/textfile.txt"}
	expect.True(ok)
	expect.Equal(want, result)


	//resultWithLeadingLash, ok2 := mapper.Get("/global")
	//expect.True(ok2)
	//expect.Equal(want, resultWithLeadingLash)
}
