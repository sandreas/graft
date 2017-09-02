package transfer

import (
	"github.com/sandreas/graft/pattern"
	"regexp"
	"os"
	"path/filepath"
	"log"
	"github.com/sandreas/graft/designpattern/observer"
	"errors"
)

type AbstractStrategy struct {
	designpattern.Observable

	SourcePattern          *pattern.SourcePattern
	DestinationPattern     *pattern.DestinationPattern
	CompiledSourcePattern  *regexp.Regexp
	TransferredDirectories []string
	KeepTimes              bool
	DryRun					bool

}

func (strategy *AbstractStrategy) PerformFileTransfer(src string, dst string, srcStat os.FileInfo) error {
	return errors.New("method PerformFileTransfer is abstract and must be overridden in strategy")
}

func (strategy *AbstractStrategy) Cleanup() error {
	return errors.New("method Cleanup is abstract and must be overridden in strategy")
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
	if strategy.DestinationPattern.Pattern == "" {
		return strategy.DestinationPattern.Path + src[len(strategy.SourcePattern.Path):]
	} else {
		return strategy.CompiledSourcePattern.ReplaceAllString(src, strategy.DestinationPattern.Path+"/"+strategy.DestinationPattern.Pattern)
	}
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
