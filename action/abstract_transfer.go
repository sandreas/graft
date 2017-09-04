package action

import (
	"github.com/spf13/afero"
	"github.com/sandreas/graft/pattern"
	"github.com/urfave/cli"
)

type AbstractTransferAction struct {
	AbstractAction
	destinationPattern *pattern.DestinationPattern
}


func (action *AbstractTransferAction) prepareTransferAction(c *cli.Context, positionalArgumentsCount int, positionalDefaultsIfUnset ...string) error {
	if err := action.PrepareExecution(c, positionalArgumentsCount, positionalDefaultsIfUnset...); err != nil {
		return err
	}
	if err := action.locateSourceFiles(); err != nil {
		return cli.NewExitError(err.Error(), ErrorLocateSourceFiles)
	}

	if err := action.prepareDestination(); err != nil {
		return cli.NewExitError(err.Error(), ErrorPrepareDestination)
	}
	return nil
}


func (action *AbstractTransferAction) prepareDestination() error {
	destinationFs := afero.NewOsFs()
	action.destinationPattern = pattern.NewDestinationPattern(destinationFs, action.PositionalArguments.Get(1))
	return nil
}