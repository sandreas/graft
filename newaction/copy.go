package newaction

import (
	"github.com/sandreas/graft/newpattern"
	"github.com/sandreas/graft/newtransfer"
	"path/filepath"
	"github.com/spf13/afero"
	"os"
	"errors"
	"github.com/sandreas/graft/newdesignpattern/observer"
	"github.com/sandreas/graft/newoptions"
)

const DRY_RUN = 1

type CopyAction struct {
	newdesignpattern.Observable
	Fs            afero.Fs
	src           newpattern.SourcePattern
	sourceFiles   []string
	copyStrategy  newtransfer.CopyStrategy
	currentSrc    os.FileInfo
	currentSrcErr error
	currentDst    os.FileInfo
	currentDstErr error
	dryRun        bool
}

func NewCopyAction(sourceFiles []string, copyStrategy newtransfer.CopyStrategy, params ...newoptions.BitFlag) *CopyAction {
	copyAction := &CopyAction{
		Fs:           afero.NewOsFs(),
		sourceFiles:  sourceFiles,
		copyStrategy: copyStrategy,
	}

	bitFlags := newoptions.NewBitFlagParser(params...)
	copyAction.dryRun = bitFlags.HasFlag(DRY_RUN)

	return copyAction
}

func (act *CopyAction) Copy(srcPattern *newpattern.SourcePattern, dstPattern *newpattern.DestinationPattern) error {
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

		if ! act.dryRun {
			loopErr = act.transfer(src, dst)
			if loopErr != nil {
				act.NotifyObservers("\n    - failed (" + loopErr.Error() + ")\n")
				err = errors.New("some files failed to transfer")
			}
		}

	}
	return err
}
func (act *CopyAction) transfer(src string, dst string) error {
	srcStat, srcStatErr := act.Fs.Stat(src)
	if srcStatErr != nil {
		return errors.New("Could not read source file " + src)
	}

	dstStat, dstStatErr := os.Stat(dst)
	if srcStat.IsDir() {
		if os.IsNotExist(dstStatErr) {
			return act.Fs.MkdirAll(dst, srcStat.Mode())
		}

		if !dstStat.IsDir() {
			return errors.New("transfer failed: " + src + " is a directory, " + dst + " exists and is not a directory")
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
	}

	return act.copyStrategy.Copy(src, dst)
}
