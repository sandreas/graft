package newtransfer

import (
	"github.com/spf13/afero"
	"io"
	"os"
	"errors"
	"time"
	"github.com/sandreas/graft/newdesignpattern/observer"
)

type FileCopier struct {
	newdesignpattern.Observable
	Fs afero.Fs
	ProgressHandler *CopyProgressHandler
	bufferSize int64
}


func NewFileCopier() *FileCopier {
	copier := &FileCopier{
		Fs: afero.NewOsFs(),
		ProgressHandler: nil,
		bufferSize: 1024 * 32,
	}
	return copier
}


func (c *FileCopier) Copy(s, d string)  error {

	srcStats, err := c.Fs.Stat(s)
	if err != nil {
		return err
	}

	srcSize := srcStats.Size()
	dstSize := int64(0)
	dstStats, err := c.Fs.Stat(d)
	if err == nil {
		dstSize = dstStats.Size()
	} else if !os.IsNotExist(err) {
		return err
	}


	if dstSize > srcSize {
		return errors.New("File cannot be resumed, destination is larger than source")
	}

	if srcSize == dstSize {
		c.handleProgress(dstSize, srcSize, c.bufferSize)
		return nil
	}

	src, err := c.Fs.OpenFile(s, os.O_RDONLY, srcStats.Mode())
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := c.Fs.OpenFile(d, os.O_RDWR | os.O_CREATE, srcStats.Mode())
	if err != nil {
		return err
	}
	defer dst.Close()


	src.Seek(dstSize, 0)
	dst.Seek(dstSize, 0)


	buf := make([]byte, c.bufferSize)
	bytesTransferred := dstSize
	for {
		n, err := src.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := dst.Write(buf[:n]); err != nil {
			return err
		}
		bytesTransferred += int64(n);
		newBufferSize := c.handleProgress(bytesTransferred, srcSize, c.bufferSize)
		if newBufferSize != c.bufferSize {
			c.bufferSize = newBufferSize
			buf = make([]byte, c.bufferSize)
		}
	}
	dst.Sync()

	return nil
}

func(c *FileCopier) handleProgress(bytesTransferred, srcSize, bufferSize int64) (int64) {
	if c.ProgressHandler == nil {
		return bufferSize
	}
	newBufferSize, message := c.ProgressHandler.Update(bytesTransferred, srcSize, bufferSize, time.Now())
	c.NotifyObservers(message)
	return newBufferSize
}

//// CopyFile copies a file from src to dst. If src and dst files exist, and are
//// the same, then return success. Otherise, attempt to create a hard link
//// between the two files. If that fail, copy the file contents from src to dst.
//func CopyFile(src, dst string) (err error) {
//	sfi, err := os.Stat(src)
//	if err != nil {
//		return
//	}
//	if !sfi.Mode().IsRegular() {
//		// cannot copy non-regular files (e.g., directories,
//		// symlinks, devices, etc.)
//		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
//	}
//	dfi, err := os.Stat(dst)
//	if err != nil {
//		if !os.IsNotExist(err) {
//			return
//		}
//	} else {
//		if !(dfi.Mode().IsRegular()) {
//			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
//		}
//		if os.SameFile(sfi, dfi) {
//			return
//		}
//	}
//	if err = os.Link(src, dst); err == nil {
//		return
//	}
//	err = copyFileContents(src, dst)
//	return
//}
//
//// copyFileContents copies the contents of the file named src to the file named
//// by dst. The file will be created if it does not already exist. If the
//// destination file exists, all it's contents will be replaced by the contents
//// of the source file.
//func copyFileContents(src, dst string) (err error) {
//	in, err := os.Open(src)
//	if err != nil {
//		return
//	}
//	defer in.Close()
//	out, err := os.Create(dst)
//	if err != nil {
//		return
//	}
//	defer func() {
//		cerr := out.Close()
//		if err == nil {
//			err = cerr
//		}
//	}()
//	if _, err = io.Copy(out, in); err != nil {
//		return
//	}
//	err = out.Sync()
//	return
//}
