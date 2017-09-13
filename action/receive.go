package action

import (
	"log"
	"time"

	"github.com/sandreas/graft/transfer"
	"github.com/urfave/cli"

	"context"
	"errors"

	"fmt"

	"github.com/grandcat/zeroconf"
	"bufio"
	"os"
	"strings"
	"strconv"
	"net"
)

type ReceiveAction struct {
	AbstractTransferAction
	serverEntries []*MdnsServerEntry
}

type MdnsServerEntry struct {
	Host string
	Port int
}

func (action *ReceiveAction) Execute(c *cli.Context) error {
	action.serverEntries = []*MdnsServerEntry{}
	if err := action.PrepareExecution(c, 2, "*", "$1"); err != nil {
		return err
	}
	//if action.CliParameters.Host== "" {
	//	service := "_graft._tcp"
	//	domain := "local"
	//
	//	waitTime := 10
	//
	//	// Discover all services on the network (e.g. _workstation._tcp)
	//	resolver, err := zeroconf.NewResolver(nil)
	//	if err != nil {
	//		log.Fatalln("Failed to initialize resolver:", err.Error())
	//	}
	//
	//	entries := make(chan *zeroconf.ServiceEntry)
	//	go func(results <-chan *zeroconf.ServiceEntry) {
	//		for entry := range results {
	//			fmt.Println(entry)
	//		}
	//		fmt.Println("No more entries.")
	//	}(entries)
	//
	//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(waitTime))
	//	defer cancel()
	//	err = resolver.Browse(ctx, service, domain, entries)
	//	if err != nil {
	//		log.Fatalln("Failed to browse:", err.Error())
	//	}
	//
	//	<-ctx.Done()
	//	// Wait some additional time to see debug messages on go routine shutdown.
	//	time.Sleep(1 * time.Second)
	//	return nil
	//}
	//return nil

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
	var err error
	action.suppressablePrintf("hostname parameter is not set, trying to find graft servers...\n")

	service := "_graft._tcp"
	domain := ""

	waitTime := 2

	// Discover all services on the network (e.g. _workstation._tcp)
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err.Error())
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			fmt.Printf("%+v\n", entry)
			server := &MdnsServerEntry{
				Host: entry.HostName,
				Port: entry.Port,
			}
			action.serverEntries = append(action.serverEntries, server)
			println("found new server: " + fmt.Sprintf("%s:%d", server.Host, server.Port))
		}

	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(waitTime))
	defer cancel()
	err = resolver.Browse(ctx, service, domain, entries)
	if err != nil {
		return err
	}

	<-ctx.Done()
	// Wait some additional time to see debug messages on go routine shutdown.
	time.Sleep(1 * time.Second)
	action.chooseServerAndReceive()
	return nil
}
func (action *ReceiveAction) chooseServerAndReceive() error {
	serverCount := len(action.serverEntries)
	log.Printf("server entries found: %d", serverCount)

	if serverCount == 0 {
		return errors.New("graft did not find a server instance to receive from")
	}
	var selectedServer *MdnsServerEntry
	if serverCount == 1 {
		selectedServer = action.serverEntries[0]
	} else {
		action.suppressablePrintf("found %d servers, choose the one to receive from:\n", serverCount)

		for i := 0; i < serverCount; i++ {
			fmt.Printf("%d.) %s:%d\n", i+1, action.serverEntries[i].Host, action.serverEntries[i].Port)
		}

		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Choose server: ")
			text, _ := reader.ReadString('\n')

			chosenServerNum, err := strconv.Atoi(strings.Trim(text, "\r\n"))
			if err != nil || chosenServerNum < 1 || chosenServerNum > len(action.serverEntries) {
				fmt.Println("Invalid choice, please specify a valid number")
			} else {
				selectedServer = action.serverEntries[chosenServerNum-1]
				break
			}
		}
	}

	action.suppressablePrintf("selected server %s:%d\n", selectedServer.Host, selectedServer.Port)

	addr, err := net.LookupIP(selectedServer.Host)
	if err != nil {
		log.Printf("Could not lookup host %s\n", selectedServer.Host)
		action.CliParameters.Host = selectedServer.Host
	} else {
		for _, ip := range addr {
			ip := ip.To4()
			if ip == nil {
				ip = ip.To16()
			}
			if ip == nil {
				continue
			}

			log.Printf("Host lookup resolved to %s\n", ip)
			action.CliParameters.Host = ip.String()
			break
		}

	}

	action.CliParameters.Port = selectedServer.Port

	action.suppressablePrintf("connecting to %s:%d\n", action.CliParameters.Host, selectedServer.Port)


	action.receive()
	return nil

}
