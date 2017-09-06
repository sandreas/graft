package action

import (
	"log"
	"github.com/urfave/cli"
	"github.com/sandreas/graft/transfer"
	"time"
)

type CopyAction struct {
	AbstractTransferAction
}

func (action *CopyAction) Execute(c *cli.Context) error {
	log.Printf("copy")

	if err := action.prepareTransferAction(c, 2); err != nil {
		return err
	}

	if err := action.CopyFiles(); err != nil {
		return cli.NewExitError(err.Error(), ErrorCopyFiles)
	}

	return nil
}

func (action *CopyAction) CopyFiles() error {
	messagePrinter := transfer.NewMessagePrinterObserver(action.suppressablePrintf)
	transferStrategy, err := transfer.NewTransferStrategy(transfer.CopyResumed, action.sourcePattern, action.destinationPattern)
	if err != nil {
		return err
	}
	transferStrategy.ProgressHandler = transfer.NewCopyProgressHandler(int64(32*1024), 1*time.Second)
	transferStrategy.RegisterObserver(messagePrinter)
	transferStrategy.DryRun = action.CliContext.Bool("dry-run")
	transferStrategy.KeepTimes = action.CliContext.Bool("times")
	return transferStrategy.Perform(action.locator.SourceFiles)

}
