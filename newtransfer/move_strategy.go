package newtransfer

import (
	"os"
	"sort"

	"github.com/spf13/afero"
)

type MoveStrategy struct {
	TransferStrategyInterface
	Fs           afero.Fs
	dirsToRemove []string
}

func NewMoveStrategy() *MoveStrategy {
	return &MoveStrategy{
		Fs:           afero.NewOsFs(),
		dirsToRemove: []string{},
	}
}

func (c *MoveStrategy) Transfer(s, d string) error {
	stat, err := c.Fs.Stat(s)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		c.dirsToRemove = append(c.dirsToRemove, s)
	}

	return os.Rename(s, d)
}

func (c *MoveStrategy) CleanUp() error {
	// sort and reverse iterate over dirs to remove
	sort.Strings(c.dirsToRemove)
	sliceLen := len(c.dirsToRemove)
	for i := sliceLen - 1; i >= 0; i-- {
		if err := c.Fs.Remove(c.dirsToRemove[i]); err != nil {
			return err
		}
	}
	return nil
}
