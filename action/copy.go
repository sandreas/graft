package action

import (
	"log"
	"github.com/urfave/cli"
	"github.com/sandreas/graft/transfer"
	"time"
)

type CopyAction struct {
	*AbstractTransferAction
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
	copyStrategy := transfer.NewCopyStrategy()
	copyStrategy.ProgressHandler = transfer.NewCopyProgressHandler(int64(32*1024), 1*time.Second)
	copyStrategy.RegisterObserver(messagePrinter)
	copyStrategy.DryRun = action.CliContext.Bool("dry-run")
	copyStrategy.KeepTimes = action.CliContext.Bool("times")
	return copyStrategy.Perform(action.locator.SourceFiles)

}
