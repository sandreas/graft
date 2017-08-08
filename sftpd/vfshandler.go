package sftpd

import (
	"io"
	"log"
	"os"
	"github.com/pkg/sftp"
)

type vfs struct {
	pathMap PathMapper
}

func VfsHandler(mapper *PathMapper) sftp.Handlers {
	virtualFileSystem := &vfs{
		pathMap: *mapper,
	}
	return sftp.Handlers{
		FileGet: virtualFileSystem,
		FilePut: virtualFileSystem,
		FileCmd: virtualFileSystem,
		FileInfo: virtualFileSystem,
	}
}

func dumpSftpRequest(message string, r sftp.Request) {
	log.Println(message, "Filepath: ", r.Filepath, ", Target: ", r.Target, ", Method: ", r.Method)
}

func (fs *vfs) Fileread(r sftp.Request) (io.ReaderAt, error) {
	dumpSftpRequest("Fileread: ", r)

	filePath, err := fs.pathMap.PathTo(r.Filepath)
	if err == nil {
		return os.Open(filePath)
	}
	return nil, err
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
	switch r.Method {
	case "List":
		listing, ok := fs.pathMap.List(r.Filepath)
		if ! ok {
			return nil, os.ErrInvalid
		}

		statList := make([]os.FileInfo, len(listing))
		for i, fileName := range listing {
			stat, err := fs.pathMap.Stat(fileName)
			if err != nil {
				log.Println("Could not stat file", fileName, err)
				continue
			}

			statList[i] = stat
			log.Println("Stat for file "+fileName+": isDir=>", stat.IsDir(), "size=>", stat.Size())
		}
		return statList, nil
	case "Stat":
		stat, err := fs.pathMap.Stat(r.Filepath)
		if err != nil {
			return nil, err
		}
		return []os.FileInfo{stat}, nil
	}
	return nil, os.ErrInvalid
}