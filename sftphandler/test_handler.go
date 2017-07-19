package sftphandler

import (
	"os"
	"io"
	"github.com/pkg/sftp"
	"path/filepath"
)

type vfs struct {
	files []string
}

func TestHandler(matchingPaths []string) sftp.Handlers {
	virtualFileSystem := &vfs{}
	virtualFileSystem.files = matchingPaths
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


	return nil, nil
}

func (fs *vfs) Filewrite(r sftp.Request) (io.WriterAt, error) {
	dumpSftpRequest("Filewrite: ", r)

	return nil, nil

	//return nil, os.ErrInvalid
}

func (fs *vfs) Filecmd(r sftp.Request) error {
	dumpSftpRequest("Filecmd: ", r)

	return nil
}

func (fs *vfs) Fileinfo(r sftp.Request) ([]os.FileInfo, error) {
	dumpSftpRequest("Fileinfo: ", r)

	switch r.Method {
	case "List":
		//var err error
		//batch_size := 10
		//current_offset := 0
		//if token := r.LsNext(); token != "" {
		//	current_offset, err = strconv.Atoi(token)
		//	if err != nil {
		//		return nil, os.ErrInvalid
		//	}
		//}
		ordered_names := []string{}
		for _, fn := range fs.files {
			println("fn:", fn)
			println("r.Filepath:", r.Filepath)
			println("dirname(" + fn +"): ", filepath.Dir(fn))

			if filepath.Dir(fn) == r.Filepath {
				println("   match!")
				ordered_names = append(ordered_names, fn)
			}
		}
		//println(ordered_names)
		//
		//sort.Sort(sort.StringSlice(ordered_names))
		//list := make([]os.FileInfo, len(ordered_names))
		//for i, fn := range ordered_names {
		//	stat, err := os.Stat(fs.files[fn])
		//	if err != nil {
		//		list[i] = stat
		//	}
		//}
		//if len(list) < current_offset {
		//	return nil, io.EOF
		//}
		//new_offset := current_offset + batch_size
		//if new_offset > len(list) {
		//	new_offset = len(list)
		//}
		//r.LsSave(strconv.Itoa(new_offset))
		//return list[current_offset:new_offset], nil
	//case "Stat":
	//	file, err := fs.fetch(r.Filepath)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	tmp := []os.FileInfo{file}
	//	println("tmp", tmp)
	//	return tmp, nil

	}
	return nil, nil

	return nil, nil
}

