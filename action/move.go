package action

import (
	"log"
	"github.com/urfave/cli"
	"github.com/sandreas/graft/transfer"
	"time"
)

type MoveAction struct {
	AbstractTransferAction
}

func (action *MoveAction) Execute(c *cli.Context) error {
	log.Printf("move")

	if err := action.prepareTransferAction(c, 2); err != nil {
		return err
	}

	if err := action.MoveFiles(); err != nil {
		return cli.NewExitError(err.Error(), ErrorMoveFiles)
	}

	return nil
}

func (action *MoveAction) MoveFiles() error {
	messagePrinter := transfer.NewMessagePrinterObserver(action.suppressablePrintf)
	transferStrategy, err := transfer.NewTransferStrategy(transfer.Move, action.sourcePattern, action.destinationPattern)
	if err != nil {
		return err
	}
	transferStrategy.ProgressHandler = transfer.NewCopyProgressHandler(int64(32*1024), 1*time.Second)
	transferStrategy.RegisterObserver(messagePrinter)
	transferStrategy.DryRun = action.CliContext.Bool("dry-run")
	transferStrategy.KeepTimes = action.CliContext.Bool("times")
	return transferStrategy.Perform(action.locator.SourceFiles)

}
