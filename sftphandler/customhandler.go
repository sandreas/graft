package sftphandler

// This serves as an example of how to implement the request server handler as
// well as a dummy backend for testing. It implements an in-memory backend that
// works as a very simple filesystem with simple flat key-value lookup system.

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"time"
	"github.com/pkg/sftp"
	"syscall"
)



// In memory file-system-y thing that the Hanlders live on
type root struct {
	*memFile
	files     map[string]*memFile
	filesLock sync.Mutex
}
// Implements os.FileInfo, Reader and Writer interfaces.
// These are the 3 interfaces necessary for the Handlers.
type memFile struct {
	name        string
	modtime     time.Time
	symlink     string
	isdir       bool
	content     []byte
	contentLock sync.RWMutex
	info        os.FileInfo
	reader      *os.File
}


// InMemHandler returns a Hanlders object with the test handlers
func CustomHandler(matchingPaths []string) sftp.Handlers {
	root := &root{
		files: make(map[string]*memFile),
	}
	root.memFile = newVirtualFile("/", true)

	fmt.Println("matchingPaths: ", matchingPaths)

	for _, element := range matchingPaths {

		stat, err := os.Stat(element)

		fmt.Println(element, err, stat.IsDir())

		// root.files["/" + element] = newMemFile(element, stat.IsDir())
		root.files["/" + element] = newRealFile(element)
	}

	fmt.Println("root.files: ", root.files)

	// root.files["/test.txt"] = newMemFile("test.txt", true)

	return sftp.Handlers{root, root, root, root}
}

// Handlers
func (fs *root) Fileread(r sftp.Request) (io.ReaderAt, error) {
	println("Fileread: ", r.Filepath)

	//fs.filesLock.Lock()
	//defer fs.filesLock.Unlock()

	file, err := fs.fetch(r.Filepath)
	if err != nil {
		return nil, err
	}
	if file.symlink != "" {
		file, err = fs.fetch(file.symlink)
		if err != nil {
			return nil, err
		}
	}
	return file.ReaderAt()
}

func (fs *root) Filewrite(r sftp.Request) (io.WriterAt, error) {
	println("Filewrite: ", r.Filepath)
	//fs.filesLock.Lock()
	//defer fs.filesLock.Unlock()
	file, err := fs.fetch(r.Filepath)
	if err == os.ErrNotExist {
		dir, err := fs.fetch(filepath.Dir(r.Filepath))
		if err != nil {
			return nil, err
		}
		if !dir.isdir {
			return nil, os.ErrInvalid
		}
		file = newRealFile(r.Filepath)
		fs.files[r.Filepath] = file
	}
	return file.WriterAt()

	//return nil, os.ErrInvalid
}

func (fs *root) Filecmd(r sftp.Request) error {
	println("Filecmd: ", r.Method)
	//fs.filesLock.Lock()
	//defer fs.filesLock.Unlock()
	switch r.Method {
	case "Setstat":
		return nil
	case "Rename":
		file, err := fs.fetch(r.Filepath)
		if err != nil {
			return err
		}
		if _, ok := fs.files[r.Target]; ok {
			return &os.LinkError{Op: "rename", Old: r.Filepath, New: r.Target,
				Err: fmt.Errorf("dest file exists")}
		}
		fs.files[r.Target] = file
		delete(fs.files, r.Filepath)
	case "Rmdir", "Remove":
		_, err := fs.fetch(filepath.Dir(r.Filepath))
		if err != nil {
			return err
		}
		delete(fs.files, r.Filepath)
	case "Mkdir":
		_, err := fs.fetch(filepath.Dir(r.Filepath))
		if err != nil {
			return err
		}
		fs.files[r.Filepath] = newVirtualFile(r.Filepath, true)
	case "Symlink":
		_, err := fs.fetch(r.Filepath)
		if err != nil {
			return err
		}
		link := newVirtualFile(r.Target, false)
		link.symlink = r.Filepath
		fs.files[r.Target] = link
	}
	return nil
}

func (fs *root) Fileinfo(r sftp.Request) ([]os.FileInfo, error) {
	println("Fileinfo: ", r.Method)

	//fs.filesLock.Lock()
	//defer fs.filesLock.Unlock()
	switch r.Method {
	case "List":
		var err error
		batch_size := 10
		current_offset := 0
		if token := r.LsNext(); token != "" {
			current_offset, err = strconv.Atoi(token)
			if err != nil {
				return nil, os.ErrInvalid
			}
		}
		ordered_names := []string{}
		for fn, _ := range fs.files {
			if filepath.Dir(fn) == r.Filepath {
				ordered_names = append(ordered_names, fn)
			}
		}
		sort.Sort(sort.StringSlice(ordered_names))
		list := make([]os.FileInfo, len(ordered_names))
		for i, fn := range ordered_names {
			list[i] = fs.files[fn]
		}
		if len(list) < current_offset {
			return nil, io.EOF
		}
		new_offset := current_offset + batch_size
		if new_offset > len(list) {
			new_offset = len(list)
		}
		r.LsSave(strconv.Itoa(new_offset))
		return list[current_offset:new_offset], nil
	case "Stat":
		file, err := fs.fetch(r.Filepath)
		if err != nil {
			return nil, err
		}

		tmp := []os.FileInfo{file}
		println("tmp", tmp)
		return tmp, nil
	case "Readlink":
		file, err := fs.fetch(r.Filepath)
		if err != nil {
			return nil, err
		}
		if file.symlink != "" {
			file, err = fs.fetch(file.symlink)
			if err != nil {
				return nil, err
			}
		}
		return []os.FileInfo{file}, nil
	}
	return nil, nil
}

func (fs *root) fetch(path string) (*memFile, error) {
	println("fetch: ", path)

	if path == "/" {
		return fs.memFile, nil
	}
	if file, ok := fs.files[path]; ok {
		println("file exists: ", file.name)
		return file, nil
	}
	println("file not exists: ", path)
	return nil, os.ErrNotExist
}



// factory to make sure modtime is set
func newVirtualFile(name string, isdir bool) *memFile {
	return &memFile{
		name:    name,
		modtime: time.Now(),
		isdir:   isdir,
		info: nil,
	}
}

func newRealFile(name string) *memFile {

	stat, _ := os.Stat(name)


	//name        string
	//modtime     time.Time
	//symlink     string
	//isdir       bool
	//content     []byte
	//contentLock sync.RWMutex
	//info        os.FileInfo

	return &memFile{
		name:    stat.Name(),
		modtime: stat.ModTime(),
		isdir:   stat.IsDir(),
		info: stat,
	}
}

// Have memFile fulfill os.FileInfo interface
func (f *memFile) Name() string {
	if f.info == nil {
		return filepath.Base(f.name)
	}
	return f.info.Name()
}
func (f *memFile) Size() int64 {
	if f.info == nil {
		return int64(len(f.content))
	}
	return f.info.Size()
}
func (f *memFile) Mode() os.FileMode {
	if f.info == nil {
		ret := os.FileMode(0644)
		if f.isdir {
			ret = os.FileMode(0755) | os.ModeDir
		}
		if f.symlink != "" {
			ret = os.FileMode(0777) | os.ModeSymlink
		}
		return ret
	}
	return f.info.Mode()
}
func (f *memFile) ModTime() time.Time {
	if f.info == nil {
		return f.modtime
	}
	return f.info.ModTime()
}

func (f *memFile) IsDir() bool {
	if f.info == nil {
		return f.isdir
	}
	return f.info.IsDir()
}
func (f *memFile) Sys() interface{} {
	if f.info == nil {
		return fakeFileInfoSys()
	}
	return f.info.Sys()
}

// Read/Write
func (f *memFile) ReaderAt() (io.ReaderAt, error) {
	if f.info == nil {
		if f.isdir {
			return nil, os.ErrInvalid
		}
		return bytes.NewReader(f.content), nil
	}

	var err error;
	if f.reader == nil {
		println("opened file: ", f.info.Name())
		f.reader, err = os.Open(f.info.Name())
		// defer f.reader.Close()
	} else {
		println("already opened file: ", f.info.Name())

	}

	//f.contentLock.Lock()
	// defer fp.Close()
	//defer f.contentLock.Unlock()

	if f.reader == nil || err != nil {
		println("error opening file: ", f.info.Name(), f.reader)
		return nil, os.ErrInvalid
	}


	return f.reader, nil
}

func (f *memFile) WriterAt() (io.WriterAt, error) {
	// we do not need to write
		//if f.isdir {
			return nil, os.ErrInvalid
		//}
		//return f, nil
}
func (f *memFile) WriteAt(p []byte, off int64) (int, error) {
	//return 0, os.ErrInvalid

	// fmt.Println(string(p), off)
	// mimic write delays, should be optional
	time.Sleep(time.Microsecond * time.Duration(len(p)))
	// f.contentLock.Lock()
	// defer f.contentLock.Unlock()
	plen := len(p) + int(off)
	if plen >= len(f.content) {
		nc := make([]byte, plen)
		copy(nc, f.content)
		f.content = nc
	}
	copy(f.content[off:], p)
	return len(p), nil
}

func fakeFileInfoSys() interface{} {
	return &syscall.Stat_t{Uid: 65534, Gid: 65534}
}
