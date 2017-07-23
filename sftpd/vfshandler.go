package sftpd

import (
	"os"
	"path/filepath"
	"strings"
	"github.com/pkg/sftp"
	"sort"
	"github.com/sandreas/graft/file"
	"io"
)

var debug bool

type vfs struct {
	files []string
	pathMap map[string][]string
}

func VfsHandler(matchingPaths []string, dbg bool) sftp.Handlers {
	debug = dbg

	virtualFileSystem := &vfs{}

	sort.Strings(matchingPaths)


	virtualFileSystem.files = matchingPaths
	virtualFileSystem.pathMap = file.MakePathMap(matchingPaths)


	return sftp.Handlers{
		virtualFileSystem,
		virtualFileSystem,
		virtualFileSystem,
		virtualFileSystem,
	}
}


func dumpSftpRequest(message string, r sftp.Request) {
	if debug {
		println(message)
		println("    Filepath: " , r.Filepath)
		println("    Target: " , r.Target)
		println("    Method: " , r.Method)
		println("    Attrs: " , r.Attrs)
		println("    Flags: " , r.Flags)
	}
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
	return f, nil
}

func (fs *vfs) Filewrite(r sftp.Request) (io.WriterAt, error) {
	dumpSftpRequest("Filewrite: ", r)
	return nil, os.ErrInvalid
}

func (fs *vfs) Filecmd(r sftp.Request) error {
	dumpSftpRequest("Filecmd: ", r)
	return os.ErrInvalid
}

func (fs *vfs) Fileinfo(r sftp.Request) ([]os.FileInfo, error) {
	dumpSftpRequest("Fileinfo: ", r)
	
	requestedPath := filepath.ToSlash(r.Filepath)
	
	switch r.Method {
	case "List":
		ordered_names, ok := fs.pathMap[requestedPath]
		if ! ok {
			println("did not find requested Filepath", requestedPath)
			return nil, os.ErrInvalid
		}

		list := make([]os.FileInfo, len(ordered_names))
		for i, fileName := range ordered_names {
			stat, _ := os.Stat(fileName)
			list[i] = stat
		}
		return list, nil
	case "Stat":
		println("Stat filepath: ", requestedPath)
		foundFile := fetch(fs, requestedPath)
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
	key := filepath.ToSlash(filepath.Dir(requestedPath))
	ordered_names, ok := fs.pathMap[key]
	if ok == false {
		println("did not find requested Filepath", requestedPath)
		return ""
	}

	var foundFile = ""

	for _, b := range ordered_names {
		if b == requestedPath || b == strings.TrimLeft(requestedPath, "/") {
			foundFile = filepath.ToSlash(b)
			break
		}
	}
	return foundFile
}