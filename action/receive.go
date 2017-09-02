package action

import (
	"github.com/urfave/cli"
)

type ReceiveAction struct {
	*AbstractAction
}

func (action *ReceiveAction) Execute(c *cli.Context) error {

	// action.PrepareExecution(c, 1, defaultsForPositionalArguments)

	// TODO
	//action.PrepareExecution(c, 1, "*")
	//
	//if action.CliContext.String("host") == "" {
	//	action.lookupServiceAndReceive()
	//} else {
	//	action.receive()
	//}
	return nil
}
//func (action *ReceiveAction) receive() error {
//	if err := action.locateSourceFiles(); err != nil {
//		return cli.NewExitError(err.Error(), ErrorLocateSourceFiles)
//	}
//	return nil
//}
//func (action *ReceiveAction) lookupServiceAndReceive() {
//	resolver, err := bonjour.NewResolver(nil)
//	if err != nil {
//		log.Println("Failed to initialize resolver:", err.Error())
//		os.Exit(1)
//	}
//
//	results := make(chan *bonjour.ServiceEntry)
//
//	go func(results chan *bonjour.ServiceEntry, exitCh chan<- bool) {
//		for e := range results {
//			fmt.Printf("%v", e)
//			action.CliContext.Set("host", e.HostName)
//			action.CliContext.Set("port", string(e.Port))
//			action.receive()
//			exitCh <- true
//			time.Sleep(1e9)
//			os.Exit(0)
//		}
//	}(results, resolver.Exit)
//
//	err = resolver.Lookup("graft", "_graft._tcp", "", results)
//	if err != nil {
//		log.Println("Failed to browse:", err.Error())
//	}
//	select {}
//}
