package action

import (
	"github.com/urfave/cli"
	"github.com/oleksandr/bonjour"
	"log"
	"time"
	"github.com/sandreas/graft/transfer"

	"errors"
)

type ReceiveAction struct {
	AbstractTransferAction
}

type MdnsServerEntry struct {
	Host string
	Port int
}

func (action *ReceiveAction) Execute(c *cli.Context) error {
	if err := action.PrepareExecution(c, 2, "*", "$1"); err != nil {
		return err
	}

	if action.shouldLookup() {
		return action.lookupServiceAndReceive()
	} else {
		return action.receive()
	}
	return nil
}
func (action *ReceiveAction) shouldLookup() bool {
	return action.CliContext.String("host") == ""
}
func (action *ReceiveAction) receive() error {

	action.CliParameters.Client = true

	action.suppressablePrintf("receive from %s@%s:%d", action.CliContext.String("username"), action.CliParameters.Host, action.CliParameters.Port)

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

	if err := action.prepareDestination(); err != nil {
		return cli.NewExitError(err.Error(), ErrorPrepareDestination)
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

func (action *ReceiveAction) lookupServiceAndReceive() error {
	action.suppressablePrintf("hostname parameter is not set, trying to find graft servers...\n")

	resolver, err := bonjour.NewResolver(nil)
	if err != nil {
		log.Println("Failed to initialize resolver:", err.Error())
		return cli.NewExitError("Failed to initialize resolver: "+err.Error(), ErrorFailedToInitializeResolver)
	}
	results := make(chan *bonjour.ServiceEntry)

	err = resolver.Lookup("graft", "_graft._tcp.", "", results)
	if err != nil {
		return cli.NewExitError("Could not find graft server: "+err.Error(), ErrorNoGraftServerAvailable)
	}
	serverEntries := []*MdnsServerEntry{}
	retriesWithoutNewEntry := 0
	receiveTriggered := false
	for {
		select {
		case nextResult := <-results:
			server := &MdnsServerEntry{
				Host: nextResult.HostName,
				Port: nextResult.Port,
			}
			serverEntries = append(serverEntries, server)
			log.Printf("found new server %s:%d\n", server.Host, server.Port)
		default:
			retriesWithoutNewEntry++
			log.Printf("try %d\n", retriesWithoutNewEntry)
			if !receiveTriggered && retriesWithoutNewEntry > 20 {
				receiveTriggered = true
				return action.chooseServerAndReceive(serverEntries, resolver.Exit)
			}

		}
		time.Sleep(200 * time.Millisecond)
	}
}
func (action *ReceiveAction) chooseServerAndReceive(serverEntries []*MdnsServerEntry, exitCh chan<- bool) error {
	serverCount := len(serverEntries)
	log.Printf("server entries found: %d", serverCount)

	if serverCount == 0 {
		//action.suppressablePrintf("graft did not find a server instance to receive from, exiting\n")
		exitCh <- true
		time.Sleep(1e9)
		return errors.New("graft did not find a server instance to receive from, exiting")
	}
	var selectedServer *MdnsServerEntry
	if serverCount == 1 {
		selectedServer = serverEntries[0]
	} else {
		//println("found multiple servers to receive from - please provide hostname, port, username and password", serverCount)
		exitCh <- true
		time.Sleep(1e9)
		return errors.New("found multiple servers to receive from - please provide hostname, port, username and password")

		// Todo: Handle multiple found servers
		//action.suppressablePrintf("found %d servers, choose the one to receive from:\n", serverCount)
		//for i:=0;i<serverCount;i++ {
		//	fmt.Printf("%d.)  %s:%d\n", i+1, serverEntries[i].Host, serverEntries[i].Port)
		//}
		//
		//reader := bufio.NewReader(os.Stdin)
		//fmt.Print("Enter text: ")
		//text, _ := reader.ReadString('\n')
		//fmt.Println(text)

	}
	action.CliParameters.Host = selectedServer.Host
	action.CliParameters.Port = selectedServer.Port
	action.receive()
	exitCh <- true
	time.Sleep(1e9)
	return nil

}
