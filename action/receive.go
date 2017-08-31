package action

import (
	"github.com/urfave/cli"
	"github.com/hashicorp/mdns"
	"fmt"
)

type ReceiveAction struct {
	AbstractAction
}

func (action *ReceiveAction) Execute(c *cli.Context) error {

	// action.PrepareExecution(c, 1, defaultsForPositionalArguments)

	action.PrepareExecution(c, 1, "*")

	if action.CliContext.String("host") == "" {
		action.lookupService()
	} else {
		fmt.Println("kaputt")
	}


	//
	//log.Printf("serve")
	//if err := action.locateSourceFiles(); err != nil {
	//	return cli.NewExitError(err.Error(), ErrorLocateSourceFiles)
	//}
	//if err := action.ServeFoundFiles(); err != nil {
	//	return cli.NewExitError(err.Error(), ErrorStartingServer)
	//}
	return nil
}
func (action *ReceiveAction) lookupService() {
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	go func() {
		for entry := range entriesCh {
			fmt.Printf("Got new entry: %v\n", entry)
		}
	}()

	// Start the lookup
	mdns.Lookup("_graft._tcp", entriesCh)
	close(entriesCh)
}
