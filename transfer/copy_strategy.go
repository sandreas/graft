package transfer

import (
	"errors"
	"io"
	"os"
	"time"
	"github.com/sandreas/graft/pattern"
)

type CopyStrategy struct {
	AbstractStrategy
	ProgressHandler *CopyProgressHandler
	bufferSize      int64
}

func NewCopyStrategy(src *pattern.SourcePattern, dst *pattern.DestinationPattern) (*CopyStrategy, error) {
	strategy := &CopyStrategy{
		ProgressHandler: nil,
		bufferSize:      1024 * 32,
	}
	strategy.SourcePattern = src
	strategy.DestinationPattern = dst

	var err error

	strategy.CompiledSourcePattern, err = strategy.SourcePattern.Compile()
	return strategy, err
}


func (strategy *CopyStrategy) PerformFileTransfer(s string, d string, srcStats os.FileInfo) error {
	return strategy.Copy(s, d, srcStats)
}

func (strategy *CopyStrategy) CleanUp() error {
	return nil
}


func (strategy *CopyStrategy) Copy(s, d string, srcStats os.FileInfo) error {

	srcSize := srcStats.Size()
	dstSize := int64(0)
	dstStats, err := strategy.DestinationPattern.Fs.Stat(d)

	dstExists := true
	if err == nil {
		dstSize = dstStats.Size()
	} else if !os.IsNotExist(err) {
		return err
	} else {
		dstExists = false
	}

	if dstSize > srcSize {
		return errors.New("File cannot be resumed, destination is larger than source")
	}

	strategy.handleProgress(dstSize, srcSize, strategy.bufferSize)

	if dstExists && srcSize == dstSize {
		return nil
	}

	src, err := strategy.SourcePattern.Fs.OpenFile(s, os.O_RDONLY, srcStats.Mode())
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := strategy.DestinationPattern.Fs.OpenFile(d, os.O_RDWR|os.O_CREATE, srcStats.Mode())
	if err != nil {
		return err
	}
	defer dst.Close()

	src.Seek(dstSize, 0)
	dst.Seek(dstSize, 0)

	buf := make([]byte, strategy.bufferSize)
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
		bytesTransferred += int64(n)
		newBufferSize := strategy.handleProgress(bytesTransferred, srcSize, strategy.bufferSize)
		if newBufferSize != strategy.bufferSize {
			strategy.bufferSize = newBufferSize
			buf = make([]byte, strategy.bufferSize)
		}
	}
	dst.Sync()

	return nil
}

func (strategy *CopyStrategy) handleProgress(bytesTransferred, srcSize, bufferSize int64) int64 {
	if strategy.ProgressHandler == nil {
		return bufferSize
	}
	newBufferSize, message := strategy.ProgressHandler.Update(bytesTransferred, srcSize, bufferSize, time.Now())
	strategy.NotifyObservers(message)
	return newBufferSize
}

