package transfer

import (
	"github.com/sandreas/graft/pattern"
	"regexp"
	"os"
	"path/filepath"
	"log"
	"github.com/sandreas/graft/designpattern/observer"
	"errors"
	"strings"
	"io"
	"time"
)

const (
	Copy = 1
	//Move
)

type AbstractStrategy struct {
	designpattern.Observable

	SourcePattern          *pattern.SourcePattern
	DestinationPattern     *pattern.DestinationPattern
	CompiledSourcePattern  *regexp.Regexp
	TransferredDirectories []string
	KeepTimes              bool
	DryRun                 bool

	transferMode    int
	ProgressHandler *CopyProgressHandler
	bufferSize      int64
}

func NewTransferStrategy(transferMode int, src *pattern.SourcePattern, dst *pattern.DestinationPattern) (*AbstractStrategy, error) {
	strategy := &AbstractStrategy{
		ProgressHandler: nil,
		bufferSize:      1024 * 32,
		transferMode:    transferMode,
	}
	strategy.SourcePattern = src
	strategy.DestinationPattern = dst
	var err error

	strategy.CompiledSourcePattern, err = strategy.SourcePattern.Compile()
	return strategy, err
}

func (strategy *AbstractStrategy) PerformFileTransfer(src string, dst string, srcStat os.FileInfo) error {
	if strategy.transferMode == Copy {
		return strategy.CopyResumed(src, dst, srcStat)
	}
	return nil
}

func (strategy *AbstractStrategy) CopyResumed(s, d string, srcStats os.FileInfo) error {

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

func (strategy *AbstractStrategy) handleProgress(bytesTransferred, srcSize, bufferSize int64) int64 {
	if strategy.ProgressHandler == nil {
		return bufferSize
	}
	newBufferSize, message := strategy.ProgressHandler.Update(bytesTransferred, srcSize, bufferSize, time.Now())
	strategy.NotifyObservers(message)
	return newBufferSize
}

func (strategy *AbstractStrategy) Cleanup() error {
	if strategy.transferMode == Copy {
		return nil
	}

	return nil
}

func (strategy *AbstractStrategy) Perform(strings []string) error {
	var err, returnError error

	strategy.NotifyObservers("\n")
	for _, src := range strings {
		err = strategy.PerformSingleTransfer(src)
		if err != nil {
			strategy.NotifyObservers("\n    - failed (" + err.Error() + ")\n")
			returnError = errors.New("some files failed to transfer")
		}
	}
	strategy.Cleanup()
	return returnError
}

func (strategy *AbstractStrategy) DestinationFor(src string) string {

	if strategy.SourcePattern.IsFile() && strategy.DestinationPattern.IsFile() {
		return strategy.DestinationPattern.Path
	}

	// source pattern points to an existing file or directory
	if strategy.SourcePattern.Pattern == "" {
		sourceParentDir := filepath.ToSlash(filepath.Dir(strategy.SourcePattern.Path))
		destinationPathParts := []string{
			strategy.DestinationPattern.Path,
		}

		if strategy.DestinationPattern.Pattern != "" {
			destinationPathParts = append(destinationPathParts, strings.TrimRight(strategy.DestinationPattern.Pattern, "\\/"))
		}

		sourcePartAppendToDestination := strings.Trim(strings.TrimPrefix(src, sourceParentDir), "\\/")
		destinationPathParts = append(destinationPathParts, sourcePartAppendToDestination)

		return strings.Join(destinationPathParts, "/")
	}

	// destination pattern points to an existing file or directory
	if strategy.DestinationPattern.Pattern == "" {
		return strategy.DestinationPattern.Path + src[len(strategy.SourcePattern.Path):]
	}

	return strategy.CompiledSourcePattern.ReplaceAllString(src, strategy.DestinationPattern.Path+"/"+strategy.DestinationPattern.Pattern)
}

func (strategy *AbstractStrategy) PerformSingleTransfer(src string) error {
	srcStat, err := strategy.SourcePattern.Fs.Stat(src)
	if err != nil {
		return err
	}

	dst := strategy.DestinationFor(src)

	strategy.NotifyObservers(src + " => " + dst + "\n")

	if strategy.DryRun {
		return nil
	}

	if srcStat.IsDir() {
		return strategy.PerformDirectoryTransfer(src, dst, srcStat, true)
	}

	if err := strategy.EnsureDirectoryOfFileExists(src, dst); err != nil {
		return err
	}

	if err := strategy.PerformFileTransfer(src, dst, srcStat); err != nil {
		return err
	}

	return nil
}

func (strategy *AbstractStrategy) EnsureDirectoryOfFileExists(src, dst string) error {
	_, err := strategy.DestinationPattern.Fs.Stat(dst)
	if os.IsNotExist(err) || strategy.KeepTimes {
		srcDirName := filepath.Dir(src)
		srcDirStat, err := strategy.DestinationPattern.Fs.Stat(srcDirName)
		if err != nil {
			log.Printf("Could not stat directory %s of file %s", srcDirName, src)
			return err
		}

		dstDirName := filepath.Dir(dst)
		return strategy.PerformDirectoryTransfer(srcDirName, dstDirName, srcDirStat, false)
	}
	return nil

}
func (strategy *AbstractStrategy) PerformDirectoryTransfer(src, dst string, srcStat os.FileInfo, shouldRemoveAfterTransfer bool) error {
	err := strategy.DestinationPattern.Fs.MkdirAll(dst, srcStat.Mode())
	if err == nil && shouldRemoveAfterTransfer {
		strategy.TransferredDirectories = append(strategy.TransferredDirectories, dst)
	}

	if err == nil && strategy.KeepTimes {
		err = strategy.keepTimes(dst, srcStat)
	}
	return err
}

func (strategy *AbstractStrategy) keepTimes(dst string, inStats os.FileInfo) error {
	return strategy.DestinationPattern.Fs.Chtimes(dst, inStats.ModTime(), inStats.ModTime())
}
