package newaction

import (
	"github.com/sandreas/graft/newpattern"
	"github.com/sandreas/graft/newtransfer"
	"path/filepath"
	"github.com/spf13/afero"
	"os"
	"errors"
)

type CopyAction struct {
	Fs          afero.Fs
	src         newpattern.SourcePattern
	sourceFiles []string
	copyStrategy newtransfer.FileCopier
	currentSrc 	os.FileInfo
	currentSrcErr error
	currentDst os.FileInfo
	currentDstErr error

}

func NewCopyAction(sourceFiles []string, srcPattern newpattern.SourcePattern) *CopyAction {
	copyAction := &CopyAction{
		Fs: afero.NewOsFs(),
		src: srcPattern,
		sourceFiles: sourceFiles,
	}
	return copyAction
}

func (act *CopyAction) Copy(destination *newpattern.DestinationPattern, copyStrategy *newtransfer.FileCopier) error {
	compiledPattern, err := act.src.Compile()
	if err != nil {
		return err
	}

	var dst string
	var failedTransfers map[string]string
	for _, src := range act.sourceFiles {
		if destination.Pattern == "" {
			dst = filepath.ToSlash(destination.Path + src[len(act.src.Path) + 1:])
		} else {
			dst = compiledPattern.ReplaceAllString(filepath.ToSlash(src), filepath.ToSlash(destination.Pattern))
		}

		err := act.transfer(src, dst)

		if err != nil  {
			failedTransfers[src] = dst
		}
	}

	return err
}
func (act *CopyAction) transfer(src string, dst string) error {
	act.currentSrc, act.currentSrcErr = act.Fs.Stat(src)
	if act.currentSrcErr != nil {
		return errors.New("Could not read source file " + src)
	}

	act.currentDst, act.currentDstErr = os.Stat(dst)

	if act.currentSrc.IsDir() {
		return act.handleSourceIsDirectory(src, dst)
	}

	return act.handleSourceIsFile(src, dst)
}
func (act *CopyAction) handleSourceIsFile(src, dst string) error {
	err := act.copyStrategy.Copy(src, dst)
	return err

}
func (act *CopyAction) handleSourceIsDirectory(src string, dst string) error {

	return nil
}
