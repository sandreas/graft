package action

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/sandreas/graft/designpattern/observer"
	"github.com/sandreas/graft/bitflag"
	"github.com/sandreas/graft/pattern"
	"github.com/sandreas/graft/transfer"
	"github.com/spf13/afero"
	"strings"
)

const (
	FLAG_DRY_RUN bitflag.Flag = 1 << iota
	FLAG_TIMES
)

type TransferAction struct {
	designpattern.Observable
	Fs               afero.Fs
	src              pattern.SourcePattern
	sourceFiles      []string
	transferStrategy transfer.TransferStrategyInterface
	dryRun           bool
	keepTimes        bool
	transferredDirs  []string
}

func NewTransferAction(sourceFiles []string, transferStrategy transfer.TransferStrategyInterface, params ...bitflag.Flag) *TransferAction {
	transferAction := &TransferAction{
		Fs:               afero.NewOsFs(),
		sourceFiles:      sourceFiles,
		transferStrategy: transferStrategy,
		transferredDirs:  []string{},
	}

	bitFlags := bitflag.NewParser(params...)
	transferAction.dryRun = bitFlags.HasFlag(FLAG_DRY_RUN)
	transferAction.keepTimes = bitFlags.HasFlag(FLAG_TIMES)

	return transferAction
}

func (act *TransferAction) Execute(srcPattern *pattern.SourcePattern, dstPattern *pattern.DestinationPattern) error {
	compiledPattern, err := srcPattern.Compile()
	if err != nil {
		return err
	}

	var loopErr error
	var dst string
	act.NotifyObservers("\n")
	for _, src := range act.sourceFiles {
		if dstPattern.Pattern == "" {
			dst = dstPattern.Path + src[len(srcPattern.Path):]
		} else {
			dst = compiledPattern.ReplaceAllString(src, dstPattern.Path+"/"+dstPattern.Pattern)
		}

		transferMessage := src + " => " + dst
		act.NotifyObservers(transferMessage + "\n")

		if !act.dryRun {
			loopErr = act.transfer(src, dst)
			if loopErr != nil {
				act.NotifyObservers("\n    - failed (" + loopErr.Error() + ")\n")
				err = errors.New("some files failed to transfer")
			}
		}

	}

	act.transferStrategy.CleanUp(act.transferredDirs)

	return err
}
func (act *TransferAction) transfer(src string, dst string) error {
	srcStat, err := act.Fs.Stat(src)
	if err != nil {
		return err
	}

	if srcStat.IsDir() {
		return act.transferDir(src, dst, srcStat, true)
	}

	_, err = act.Fs.Stat(dst)

	// Ensure directory of file exists
	if os.IsNotExist(err) || act.keepTimes {
		srcDirName := filepath.Dir(src)
		srcDirStat, err := act.Fs.Stat(srcDirName)
		if err != nil {
			return errors.New("Could not stat " + srcDirName + " of file " + src + ": " + err.Error())
		}

		dstDirName := filepath.Dir(dst)
		err = act.transferDir(srcDirName, dstDirName, srcDirStat, false)
		if err != nil {
			return err
		}
	}

	err = act.transferStrategy.Transfer(src, dst)

	if err == nil && act.keepTimes {
		return act.transferTimes(dst, srcStat)
	}
	return nil
}

func (act *TransferAction) transferDir(src, dst string, srcStat os.FileInfo, shouldRemoveAfterTransfer bool) error {
	err := act.Fs.MkdirAll(dst, srcStat.Mode())
	if err == nil && shouldRemoveAfterTransfer {
		act.transferredDirs = append(act.transferredDirs, strings.TrimRight(src, "/"))
	}
	if err == nil && act.keepTimes {
		err = act.transferTimes(dst, srcStat)
	}
	return err
}

func (act *TransferAction) transferTimes(dst string, inStats os.FileInfo) error {
	return act.Fs.Chtimes(dst, inStats.ModTime(), inStats.ModTime())

}
