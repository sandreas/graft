package newaction

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/sandreas/graft/newdesignpattern/observer"
	"github.com/sandreas/graft/newoptions"
	"github.com/sandreas/graft/newpattern"
	"github.com/sandreas/graft/newtransfer"
	"github.com/spf13/afero"
)


const (
	FLAG_DRY_RUN newoptions.BitFlag = 1 << iota
	FLAG_TIMES
)

type TransferAction struct {
	newdesignpattern.Observable
	Fs               afero.Fs
	src              newpattern.SourcePattern
	sourceFiles      []string
	transferStrategy newtransfer.TransferStrategyInterface
	dryRun           bool
	keepTimes        bool
}

func NewTransferAction(sourceFiles []string, transferStrategy newtransfer.TransferStrategyInterface, params ...newoptions.BitFlag) *TransferAction {
	transferAction := &TransferAction{
		Fs:               afero.NewOsFs(),
		sourceFiles:      sourceFiles,
		transferStrategy: transferStrategy,
	}

	bitFlags := newoptions.NewBitFlagParser(params...)
	transferAction.dryRun = bitFlags.HasFlag(FLAG_DRY_RUN)
	transferAction.keepTimes = bitFlags.HasFlag(FLAG_TIMES)

	return transferAction
}

func (act *TransferAction) Execute(srcPattern *newpattern.SourcePattern, dstPattern *newpattern.DestinationPattern) error {
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
	return err
}
func (act *TransferAction) transfer(src string, dst string) error {
	srcStat, srcStatErr := act.Fs.Stat(src)
	if srcStatErr != nil {
		return errors.New("Could not read source file " + src)
	}

	dstStat, dstStatErr := os.Stat(dst)
	if srcStat.IsDir() {
		if os.IsNotExist(dstStatErr) {
			e := act.Fs.MkdirAll(dst, srcStat.Mode())
			if e != nil {
				return e
			}

			if act.keepTimes {
				return act.transferTimes(dst, srcStat)
			}

			return nil
		}

		if !dstStat.IsDir() {
			return errors.New("transfer failed: " + src + " is a directory, " + dst + " exists and is not a directory")
		}

		if act.keepTimes {
			return act.transferTimes(dst, srcStat)
		}

		return nil
	}
	if !os.IsNotExist(dstStatErr) && dstStat.IsDir() {
		return errors.New("transfer failed: " + src + " is a file, " + dst + " is a directory")
	}

	// Ensure directory of file exists
	if os.IsNotExist(dstStatErr) {
		srcDirName := filepath.Dir(src)
		srcDirStat, srcDirStatErr := os.Stat(srcDirName)

		if srcDirStatErr != nil {
			return errors.New("Could not stat " + srcDirName + " of file " + src + ": " + srcDirStatErr.Error())
		}

		dstDirName := filepath.Dir(dst)
		mkdirErr := act.Fs.MkdirAll(dstDirName, srcDirStat.Mode())
		if mkdirErr != nil {
			return mkdirErr
		}

		if act.keepTimes {
			err := act.transferTimes(dstDirName, srcDirStat)
			if err != nil {
				return err
			}
		}
	}

	e := act.transferStrategy.Transfer(src, dst)

	if e != nil {
		return e
	}

	if act.keepTimes {
		return act.transferTimes(dst, srcStat)
	}
	return nil
}

func (act *TransferAction) transferTimes(dst string, inStats os.FileInfo) error {
	return act.Fs.Chtimes(dst, inStats.ModTime(), inStats.ModTime())

}