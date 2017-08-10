package transfer

import (
	"os"
	"sort"

	"github.com/spf13/afero"
)

type MoveStrategy struct {
	TransferStrategyInterface
	Fs           afero.Fs}

func NewMoveStrategy() *MoveStrategy {
	return &MoveStrategy{
		Fs:           afero.NewOsFs(),
	}
}

func (c *MoveStrategy) Transfer(s, d string) error {
	_, err := c.Fs.Stat(s)
	if err != nil {
		return err
	}

	return os.Rename(s, d)
}

func (c *MoveStrategy) CleanUp(dirsToRemove []string) error {
	// sort and reverse iterate over dirs to remove
	sort.Strings(dirsToRemove)
	sliceLen := len(dirsToRemove)
	lastDir := ""
	for i := sliceLen - 1; i >= 0; i-- {
		if dirsToRemove[i] == lastDir {
			continue
		}
		err := c.Fs.Remove(dirsToRemove[i])
		lastDir = dirsToRemove[i]
		if err != nil {
			str := err.Error()
			println(str)
			return err
		}
	}
	return nil
}
