package sftpd

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/sftp"
)

type vfs struct {
	files   []string
	pathMap map[string][]string
}

func VfsHandler(matchingPaths []string) sftp.Handlers {
	virtualFileSystem := &vfs{}

	sort.Strings(matchingPaths)

	virtualFileSystem.files = matchingPaths
	virtualFileSystem.pathMap = MakePathMap(matchingPaths)

	return sftp.Handlers{
		virtualFileSystem,
		virtualFileSystem,
		virtualFileSystem,
		virtualFileSystem,
	}
}

func dumpSftpRequest(message string, r sftp.Request) {
	log.Println(message, "Filepath: ", r.Filepath, ", Target: ", r.Target, ", Method: ", r.Method)
}

func (fs *vfs) Fileread(r sftp.Request) (io.ReaderAt, error) {
	dumpSftpRequest("Fileread: ", r)

	foundFile := fetch(fs, r.Filepath)
	log.Println("foundFile: ", foundFile)
	if foundFile == "" {
		return nil, os.ErrInvalid
	}

	f, err := os.Open(foundFile)

	if err != nil {
		log.Println("Could not open file", foundFile, err)
		return nil, os.ErrInvalid
	}
	return f, nil
}

func (fs *vfs) Filewrite(r sftp.Request) (io.WriterAt, error) {
	dumpSftpRequest("Filewrite (disabled): ", r)
	return nil, os.ErrInvalid
}

func (fs *vfs) Filecmd(r sftp.Request) error {
	dumpSftpRequest("Filecmd (disabled): ", r)
	return os.ErrInvalid
}

func (fs *vfs) Fileinfo(r sftp.Request) ([]os.FileInfo, error) {
	dumpSftpRequest("Fileinfo: ", r)

	requestedPath := filepath.ToSlash(r.Filepath)
	log.Println("requestedPath: ", requestedPath)

	switch r.Method {
	case "List":
		ordered_names, ok := fs.pathMap[requestedPath]
		if !ok {
			log.Println("did not find pathMapping for requestedPath", requestedPath)
			return nil, os.ErrInvalid
		}
		log.Println("pathMapping for "+requestedPath+" contains: ", len(ordered_names))

		list := make([]os.FileInfo, len(ordered_names))
		for i, fileName := range ordered_names {
			stat, err := os.Stat(fileName)
			if err != nil {
				log.Println("Could not stat file", fileName, err)
				continue
			}

			list[i] = stat
			log.Println("Stat for file "+fileName+": isDir=>", stat.IsDir(), "size=>", stat.Size())
		}
		return list, nil
	case "Stat":
		log.Println("Stat filepath: ", requestedPath)
		foundFile := fetch(fs, requestedPath)
		if foundFile != "" {
			log.Println("foundFile: ", foundFile)
			stat, err := os.Stat(foundFile)
			if err != nil {
				log.Println("Could not stat file", foundFile, err)
				return nil, os.ErrInvalid
			}
			return []os.FileInfo{stat}, nil
		}
		log.Println("Could not 'fetch' file for " + requestedPath)
		return nil, os.ErrInvalid
	}
	return nil, os.ErrInvalid
}

func fetch(fs *vfs, requestedPath string) string {
	log.Println("fetch requestedPath:  ", requestedPath)

	key := filepath.ToSlash(filepath.Dir(requestedPath))
	log.Println("mapping key:  ", key)

	ordered_names, ok := fs.pathMap[key]

	if ok == false {
		log.Println("did not find key in pathMap", key)
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



func MakePathMap(matchingPaths []string) map[string][]string {
	pathMap := make(map[string][]string)

	sort.Strings(matchingPaths)

	//if val, ok := dict["foo"]; ok {
	//	//do something here
	//}

	for _, path := range matchingPaths {
		key, parentPath := normalizePathMapItem(path)

		for  {
			// println("append: ", key, " => ", path)
			pathMap[key] = append(pathMap[key], path)
			path = parentPath
			//println("before => key:", key, "parentPath:", parentPath)
			key, parentPath = normalizePathMapItem(parentPath)
			//println("after  => key:", key, "parentPath:", parentPath)
			_, ok := pathMap[key]

			//println("is present?", key, ok)
			if ok {
				break
			}
		}
	}



	return pathMap
}

func normalizePathMapItem(path string) (string, string) {
	parentPath := filepath.ToSlash(filepath.Dir(path))
	key := parentPath
	if parentPath == "." {
		key = "/"
	}

	firstChar := string([]rune(key)[0])
	if firstChar != "/" {
		key = "/" + key
	}
	return key, parentPath
}
