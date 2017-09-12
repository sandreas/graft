package transfer

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/sandreas/graft/designpattern/observer"
	"github.com/sandreas/graft/pattern"
	"github.com/sandreas/graft/filesystem"
)

const (
	CopyResumed = 1
	Move        = 2
)

type Strategy struct {
	designpattern.Observable

	SourcePattern          *pattern.SourcePattern
	DestinationPattern     *pattern.DestinationPattern
	CompiledSourcePattern  *regexp.Regexp
	TransferredDirectories []string
	KeepTimes              bool
	DryRun                 bool

	transferMethod  int
	ProgressHandler *CopyProgressHandler
	bufferSize      int64
}

func NewTransferStrategy(transferMethod int, src *pattern.SourcePattern, dst *pattern.DestinationPattern) (*Strategy, error) {
	var err error
	if transferMethod < CopyResumed || transferMethod > Move {
		return nil, errors.New("invalid transfer method" + string(transferMethod))
	}
	strategy := &Strategy{
		ProgressHandler: nil,
		bufferSize:      1024 * 32,
		transferMethod:  transferMethod,
	}
	strategy.SourcePattern = src
	strategy.DestinationPattern = dst

	strategy.CompiledSourcePattern, err = strategy.SourcePattern.Compile()
	return strategy, err
}

func (strategy *Strategy) PerformFileTransfer(src string, dst string, srcStat os.FileInfo) error {
	if strategy.transferMethod == CopyResumed {
		return strategy.CopyResumed(src, dst, srcStat)
	}

	if strategy.transferMethod == Move {
		return strategy.Move(src, dst, srcStat)
	}

	return nil
}

func (strategy *Strategy) CopyResumed(s, d string, srcStats os.FileInfo) error {
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

func (strategy *Strategy) Move(src string, dst string, srcStat os.FileInfo) error {
	return strategy.SourcePattern.Fs.Rename(src, dst)
}

func (strategy *Strategy) handleProgress(bytesTransferred, srcSize, bufferSize int64) int64 {
	if strategy.ProgressHandler == nil {
		return bufferSize
	}
	newBufferSize, message := strategy.ProgressHandler.Update(bytesTransferred, srcSize, bufferSize, time.Now())
	strategy.NotifyObservers(message)
	return newBufferSize
}

func (strategy *Strategy) Cleanup() error {
	if strategy.transferMethod == CopyResumed {
		return nil
	}

	if strategy.transferMethod == Move {
		sort.Strings(strategy.TransferredDirectories)
		sliceLen := len(strategy.TransferredDirectories)
		lastDir := ""
		for i := sliceLen - 1; i >= 0; i-- {
			if strategy.TransferredDirectories[i] == lastDir {
				continue
			}
			err := strategy.SourcePattern.Fs.Remove(strategy.TransferredDirectories[i])
			lastDir = strategy.TransferredDirectories[i]
			if err != nil {
				str := err.Error()
				println(str)
				return err
			}
		}
		return nil
	}
	return nil
}

func (strategy *Strategy) Perform(strings []string) error {
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

func (strategy *Strategy) DestinationFor(src string) string {
	cleanedSrc := filesystem.CleanPath(strategy.SourcePattern.Fs, src)



	srcPatternPathLen := len(strategy.SourcePattern.Path)
	if strategy.SourcePattern.Path == "." {
		srcPatternPathLen = 0
	}


	if strategy.SourcePattern.IsFile() {
		if strategy.DestinationPattern.IsFile() {
			return strategy.DestinationPattern.Path
		}

		if strategy.DestinationPattern.Pattern == "" {
			return strategy.DestinationPattern.Path + string(os.PathSeparator) + filepath.Base(cleanedSrc)
		}

		l:=len(strategy.DestinationPattern.Pattern)
		if !os.IsPathSeparator(strategy.DestinationPattern.Pattern[l-1]) {
			return strategy.DestinationPattern.Path + string(os.PathSeparator) + strategy.DestinationPattern.Pattern
		}
		cleanedPattern := strings.TrimRight(strategy.DestinationPattern.Pattern, "\\/")
		return strategy.DestinationPattern.Path+string(os.PathSeparator)+cleanedPattern+string(os.PathSeparator)+filepath.Base(cleanedSrc)
	}



	// source pattern points to an existing file or directory
	if strategy.SourcePattern.Pattern == "" {
		sourceParentDir := filepath.Dir(strategy.SourcePattern.Path)
		destinationPathParts := []string{
			strategy.DestinationPattern.Path,
		}

		if strategy.DestinationPattern.Pattern != "" {
			destinationPathParts = append(destinationPathParts, strings.TrimRight(strategy.DestinationPattern.Pattern, "\\/"))
		}

		sourcePartAppendToDestination := strings.Trim(strings.TrimPrefix(cleanedSrc, sourceParentDir), "\\/")
		destinationPathParts = append(destinationPathParts, sourcePartAppendToDestination)

		return strings.Join(destinationPathParts, string(os.PathSeparator))
	}

	// destination pattern points to an existing file or directory
	if strategy.DestinationPattern.Pattern == "" {
		return strategy.DestinationPattern.Path + string(os.PathSeparator) + strings.TrimLeft(cleanedSrc[srcPatternPathLen:], "\\/")
	}

	return strategy.CompiledSourcePattern.ReplaceAllString(cleanedSrc, strategy.DestinationPattern.Path+string(os.PathSeparator)+strings.TrimLeft(strategy.DestinationPattern.Pattern, "\\/"))
}

func (strategy *Strategy) PerformSingleTransfer(src string) error {

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

func (strategy *Strategy) EnsureDirectoryOfFileExists(src, dst string) error {
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
func (strategy *Strategy) PerformDirectoryTransfer(src, dst string, srcStat os.FileInfo, shouldRemoveAfterTransfer bool) error {
	err := strategy.DestinationPattern.Fs.MkdirAll(dst, srcStat.Mode())
	if err == nil && shouldRemoveAfterTransfer {
		strategy.TransferredDirectories = append(strategy.TransferredDirectories, dst)
	}

	if err == nil && strategy.KeepTimes {
		err = strategy.keepTimes(dst, srcStat)
	}
	return err
}

func (strategy *Strategy) keepTimes(dst string, inStats os.FileInfo) error {
	return strategy.DestinationPattern.Fs.Chtimes(dst, inStats.ModTime(), inStats.ModTime())
}
