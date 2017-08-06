package newtransfer

import (
	"github.com/spf13/afero"
	"io"
	"os"
	"errors"
	"time"
	"github.com/sandreas/graft/newdesignpattern/observer"
)

type CopyStrategy struct {
	newdesignpattern.Observable
	Fs              afero.Fs
	progressHandler *CopyProgressHandler
	bufferSize      int64
}


func NewCopyStrategy() *CopyStrategy {
	copier := &CopyStrategy{
		Fs:              afero.NewOsFs(),
		progressHandler: nil,
		bufferSize:      1024 * 32,
	}
	return copier
}

func (c *CopyStrategy) SetProgressHandler(h *CopyProgressHandler) {
	c.progressHandler = h
}

func (c *CopyStrategy) Transfer(s, d string)  error {

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

func(c *CopyStrategy) handleProgress(bytesTransferred, srcSize, bufferSize int64) (int64) {
	if c.progressHandler == nil {
		return bufferSize
	}
	newBufferSize, message := c.progressHandler.Update(bytesTransferred, srcSize, bufferSize, time.Now())
	c.NotifyObservers(message)
	return newBufferSize
}