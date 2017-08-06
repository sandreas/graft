package newtransfer

import (
	"github.com/spf13/afero"
	"os"
)

type MoveStrategy struct {
	TransferStrategyInterface
	Fs              afero.Fs
}


func NewMoveStrategy() *MoveStrategy {
	return &MoveStrategy{
		Fs:              afero.NewOsFs(),
	}
}

func (c *MoveStrategy) Transfer(s, d string)  error {
	return os.Rename(s, d)
}


