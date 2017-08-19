// Copyright © 2015 Jerry Jacobs <jerry.jacobs@xor-gate.org>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sftpfs

import (
	"github.com/pkg/sftp"
	"os"
	"io"
)

type File struct {
	file   *sftp.File
	client *sftp.Client
	name   string
}

func FileOpen(s *sftp.Client, name string) (*File, error) {
	fd, err := s.Open(name)
	if err != nil {
		return &File{}, err
	}
	return &File{
		file:   fd,
		client: s,
		name:   name,
	}, nil
}

func FileCreate(s *sftp.Client, name string) (*File, error) {
	fd, err := s.Create(name)
	if err != nil {
		return &File{}, err
	}
	return &File{
		file:   fd,
		client: s,
		name:   name,
	}, nil
}

func (f *File) Close() error {
	return f.file.Close()
}

func (f *File) Name() string {
	return f.file.Name()
}

func (f *File) Stat() (os.FileInfo, error) {
	return f.file.Stat()
}

func (f *File) Sync() error {
	return nil
}

func (f *File) Truncate(size int64) error {
	return f.file.Truncate(size)
}

func (f *File) Read(b []byte) (n int, err error) {
	return f.file.Read(b)
}

// TODO
func (f *File) ReadAt(b []byte, off int64) (n int, err error) {
	f.file.Seek(off, io.SeekStart)
	return f.file.Read(b)
}

// TODO
func (f *File) Readdir(count int) (res []os.FileInfo, err error) {
	name, err := NormalizeDir(f.name, f.client)
	if err != nil {
		return []os.FileInfo{}, err

	}

	dirs, err := f.client.ReadDir(name)
	if len(dirs) > count && count > 0 {
		return dirs[0:count-1], err
	}
	return dirs, err
}

// TODO
func (f *File) Readdirnames(n int) (names []string, err error) {
	dirs, err := f.Readdir(n)
	dirNames := []string{}

	if err != nil {
		return dirNames, err
	}

	for _, dir := range dirs {
		dirNames = append(dirNames, dir.Name())
	}
	return dirNames, nil
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	return f.file.Seek(offset, whence)
}

func (f *File) Write(b []byte) (n int, err error) {
	return f.file.Write(b)
}

// TODO
func (f *File) WriteAt(b []byte, off int64) (n int, err error) {
	return 0, nil
}

func (f *File) WriteString(s string) (ret int, err error) {
	return f.file.Write([]byte(s))
}
