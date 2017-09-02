package transfer

import (
	"github.com/sandreas/graft/pattern"
	"regexp"
	"os"
	"path/filepath"
	"log"
)

type AbstractStrategy struct {
	SourcePattern          *pattern.SourcePattern
	DestinationPattern     *pattern.DestinationPattern
	CompiledSourcePattern  *regexp.Regexp
	TransferredDirectories []string
	KeepTimes bool
}

func (strategy *AbstractStrategy) DestinationFor(src string) string {
	if strategy.DestinationPattern.Pattern == "" {
		return strategy.DestinationPattern.Path + src[len(strategy.SourcePattern.Path):]
	} else {
		return strategy.CompiledSourcePattern.ReplaceAllString(src, strategy.DestinationPattern.Path+"/"+strategy.DestinationPattern.Pattern)
	}
}

func (strategy *AbstractStrategy) PerformTransfer(src string) error {
	srcStat, err := strategy.SourcePattern.Fs.Stat(src)
	if err != nil {
		return err
	}

	dst := strategy.DestinationFor(src)

	if srcStat.IsDir() {
		return strategy.PerformDirectoryTransfer(src, dst, srcStat, true)
	}

	if err := strategy.EnsureDirectoryOfFileExists(src, dst); err != nil {
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
	t := inStats.ModTime()
	return strategy.DestinationPattern.Fs.Chtimes(dst, t, t)
}

