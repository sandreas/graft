package filesystem

import (
	"os"
	"time"

	"github.com/pkg/sftp"
	"github.com/spf13/afero"
	"fmt"
)

// Fs is a afero.Fs implementation that uses functions provided by the sftp package.
//
// For details in any method, check the documentation of the sftp package
// (github.com/pkg/sftp).
type SftpFs struct {
	// *io.Closer
	afero.Fs
	client  *sftp.Client
	Context *SftpFsContext
}

func NewSftpFs(host string, port int, username, password string) (afero.Fs, error) {
	hostAndPort := fmt.Sprintf("%s:%d", host, port)
	ctx, err := NewSftpFsContext(username, password, hostAndPort)
	if err != nil {
		return nil, err
	}

	fs := &SftpFs{
		client:  ctx.SftpClient,
		Context: ctx,
	}
	_, err = fs.client.Getwd()
	if err != nil {
		return nil, err
	}
	return fs, nil

}

func NormalizeDir(name string, client *sftp.Client) (string, error) {
	if name == "." {
		return client.Getwd()
	}
	return name, nil
}

func (s SftpFs) Name() string { return "sftpfs" }

func (s SftpFs) Create(name string) (afero.File, error) {
	return FileCreate(s.client, name)
}

func (s SftpFs) Mkdir(name string, perm os.FileMode) error {
	err := s.client.Mkdir(name)
	if err != nil {
		return err
	}
	return s.client.Chmod(name, perm)
}

func (s SftpFs) MkdirAll(path string, perm os.FileMode) error {
	// Fast path: if we can tell whether path is a directory or file, stop with success or error.
	dir, err := s.Stat(path)
	if err == nil {
		if dir.IsDir() {
			return nil
		}
		return err
	}

	// Slow path: make sure parent exists and then call Mkdir for path.
	i := len(path)
	for i > 0 && os.IsPathSeparator(path[i-1]) { // Skip trailing path separator.
		i--
	}

	j := i
	for j > 0 && !os.IsPathSeparator(path[j-1]) { // Scan backward over element.
		j--
	}

	if j > 1 {
		// Create parent
		err = s.MkdirAll(path[0:j-1], perm)
		if err != nil {
			return err
		}
	}

	// Parent now exists; invoke Mkdir and use its result.
	err = s.Mkdir(path, perm)
	if err != nil {
		// Handle arguments like "foo/." by
		// double-checking that directory doesn't exist.
		dir, err1 := s.Lstat(path)
		if err1 == nil && dir.IsDir() {
			return nil
		}
		return err
	}
	return nil
}

func (s SftpFs) Open(name string) (afero.File, error) {
	return FileOpen(s.client, name)
}
// changed!!!
func (s SftpFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	// return s.client.OpenFile(name, flag)
 	return FileOpen(s.client, name)
}

func (s SftpFs) Remove(name string) error {
	return s.client.Remove(name)
}

func (s SftpFs) RemoveAll(path string) error {
	// TODO have a look at os.RemoveAll
	// https://github.com/golang/go/blob/master/src/os/path.go#L66
	return nil
}

func (s SftpFs) Rename(oldname, newname string) error {
	return s.client.Rename(oldname, newname)
}

func (s SftpFs) Stat(name string) (os.FileInfo, error) {
	name, err := NormalizeDir(name, s.client)
	if err != nil {
		return nil, err
	}
	stat, err := s.client.Stat(name)
	return stat, err
}

func (s SftpFs) Lstat(p string) (os.FileInfo, error) {
	return s.client.Lstat(p)
}

func (s SftpFs) Chmod(name string, mode os.FileMode) error {
	return s.client.Chmod(name, mode)
}

func (s SftpFs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return s.client.Chtimes(name, atime, mtime)
}

func (s SftpFs) Close() {
	s.Context.Close()
}
