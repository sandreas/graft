package sftpd

import (
	"io"
	"log"
	"os"
	"github.com/sandreas/sftp"
)

type vfs struct {
	pathMap PathMapper
}

func VfsHandler(mapper *PathMapper) sftp.Handlers {
	virtualFileSystem := &vfs{
		pathMap: *mapper,
	}
	return sftp.Handlers{
		FileGet:  virtualFileSystem,
		FilePut:  virtualFileSystem,
		FileCmd:  virtualFileSystem,
		FileList: virtualFileSystem,
	}
}

func dumpSftpRequest(message string, r sftp.Request) {
	log.Println(message, "Filepath: ", r.Filepath, ", Target: ", r.Target, ", Method: ", r.Method)
}

func (fs *vfs) Fileread(r sftp.Request) (io.ReaderAt, error) {
	dumpSftpRequest("Fileread: ", r)

	filePath, err := fs.pathMap.PathTo(r.Filepath)

	if err == nil {
		f, err := os.Open(filePath)
		if err != nil {
			defer f.Close()
		}
		return f, err
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

type listerAt []os.FileInfo

// Modeled after strings.Reader's ReadAt() implementation
func (l listerAt) ListAt(ls []os.FileInfo, offset int64) (int, error) {
	var n int
	if offset >= int64(len(l)) {
		return 0, io.EOF
	}
	n = copy(ls, l[offset:])
	if n < len(ls) {
		return n, io.EOF
	}
	return n, nil
}


func (fs *vfs) Filelist(r sftp.Request) (sftp.ListerAt, error) {
	dumpSftpRequest("Fileinfo: ", r)
	switch r.Method {
	case "List":
		listing, ok := fs.pathMap.List(r.Filepath)
		if !ok {
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
		return listerAt(statList), nil
	case "Stat":
		stat, err := fs.pathMap.Stat(r.Filepath)
		if err != nil {
			return nil, err
		}
		return listerAt([]os.FileInfo{stat}), nil
	}
	return nil, os.ErrInvalid
}
