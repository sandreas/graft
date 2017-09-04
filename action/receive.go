package action

import (
	"github.com/urfave/cli"
	"github.com/oleksandr/bonjour"
	"log"
	"os"
	"time"
	"github.com/sandreas/graft/transfer"
)

type ReceiveAction struct {
	AbstractTransferAction
}

func (action *ReceiveAction) Execute(c *cli.Context) error {
	if err := action.prepareTransferAction(c, 2, "*", "$1"); err != nil {
		return err
	}

	if action.shouldLookup() {
		action.lookupServiceAndReceive()
	} else {
		action.receive()
	}
	return nil
}
func (action *ReceiveAction) shouldLookup() bool {
	return action.CliContext.String("host") == ""
}
func (action *ReceiveAction) receive() error {
	action.Settings.Client = true
	action.suppressablePrintf("receive from %s@%s:%d", action.CliContext.String("username"), action.Settings.Host, action.Settings.Port)
	if action.CliContext.String("password") == "" {
		password, err := action.promptPassword("Enter password:")
		if err != nil {
			return err
		}
		action.CliContext.Set("password", password)
	}

	if err := action.locateSourceFiles(); err != nil {
		return cli.NewExitError(err.Error(), ErrorLocateSourceFiles)
	}

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

func (action *ReceiveAction) lookupServiceAndReceive() {

	resolver, err := bonjour.NewResolver(nil)
	if err != nil {
		log.Println("Failed to initialize resolver:", err.Error())
		os.Exit(1)
	}
	results := make(chan *bonjour.ServiceEntry)

	go func(results chan *bonjour.ServiceEntry, exitCh chan<- bool) {
		for e := range results {
			action.Settings.Host = e.HostName
			action.Settings.Port = e.Port

			action.receive()
			exitCh <- true
			time.Sleep(1e9)
			os.Exit(0)
		}
	}(results, resolver.Exit)

	err = resolver.Lookup("graft", "_graft._tcp.", "", results)
	if err != nil {
		log.Println("Failed to browse:", err.Error())
	}
	select {}
}
