package sftphandler

import (
	"os"
	"io"
	"github.com/pkg/sftp"
	"sort"
	"github.com/sandreas/graft/file"
	"path/filepath"
	"strings"
)

type vfs struct {
	files []string
	pathMap map[string][]string
}


/*
graft.go
LICENSE
README.md
data
data/fixtures
data/fixtures/global/file.txt


vfs {
	path = /
	parent = nil
	chilren
}
parent = nil
path = /
children {
	graft.go

}


 */

func TestHandler(matchingPaths []string) sftp.Handlers {
	virtualFileSystem := &vfs{}

	sort.Strings(matchingPaths)


	virtualFileSystem.files = matchingPaths
	virtualFileSystem.pathMap = file.MakePathMap(matchingPaths)
	//for _, element := range matchingPaths {
	//	stat, err := os.Stat(element)
	//	if err == nil {
	//		virtualFileSystem.files = append(virtualFileSystem.files, &stat)
	//	}
	//}


	return sftp.Handlers{
		virtualFileSystem,
		virtualFileSystem,
		virtualFileSystem,
		virtualFileSystem,
	}
}

func dumpSftpRequest(message string, r sftp.Request) {
	println(message)
	println("    Filepath: " , r.Filepath)
	println("    Target: " , r.Target)
	println("    Method: " , r.Method)
	println("    Attrs: " , r.Attrs)
	println("    Flags: " , r.Flags)
}

func (fs *vfs) Fileread(r sftp.Request) (io.ReaderAt, error) {
	dumpSftpRequest("Fileread: ", r)

	foundFile := fetch(fs, r.Filepath)
	if(foundFile == "") {
		return nil, os.ErrInvalid
	}

	f, err := os.Open(foundFile)

	if err != nil {
		return nil, os.ErrInvalid
	}
	// defer f.Close()

	return f, nil
}

func (fs *vfs) Filewrite(r sftp.Request) (io.WriterAt, error) {
	dumpSftpRequest("Filewrite: ", r)

	return nil, os.ErrInvalid

	//return nil, os.ErrInvalid
}

func (fs *vfs) Filecmd(r sftp.Request) error {
	dumpSftpRequest("Filecmd: ", r)

	return os.ErrInvalid
}

func (fs *vfs) Fileinfo(r sftp.Request) ([]os.FileInfo, error) {
	dumpSftpRequest("Fileinfo: ", r)

	switch r.Method {
	case "List":
		ordered_names, ok := fs.pathMap[r.Filepath]
		if ! ok {
			println("did not find requested Filepath", r.Filepath)
			return nil, os.ErrInvalid
		}


		list := make([]os.FileInfo, len(ordered_names))
		for i, fileName := range ordered_names {
			stat, _ := os.Stat(fileName)
			list[i] = stat
		}
		return list, nil
	case "Stat":
		println("Stat filepath: ", r.Filepath)
		foundFile := fetch(fs, r.Filepath)
		if foundFile != "" {
			println("foundFile: ", foundFile)
			stat, _ := os.Stat(foundFile)
			return []os.FileInfo{stat}, nil
		}
		return nil, os.ErrInvalid
	}
	return nil, os.ErrInvalid
}


func fetch(fs *vfs, requestedPath string) string {
	key := filepath.Dir(requestedPath)
	ordered_names, ok := fs.pathMap[key]
	if ok == false {
		println("did not find requested Filepath", requestedPath)
		return ""
	}

	var foundFile = ""

	for _, b := range ordered_names {
		if b == requestedPath || b == strings.TrimLeft(requestedPath, "/") {
			foundFile = b
			break
		}
	}
	return foundFile
}