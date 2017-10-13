package action

import (
	"log"
	"time"

	"github.com/sandreas/graft/transfer"
	"github.com/urfave/cli"

	"context"
	"errors"

	"fmt"

	"bufio"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/grandcat/zeroconf"
)

type ReceiveAction struct {
	AbstractTransferAction
	serverEntries []*zeroconf.ServiceEntry
}

func (action *ReceiveAction) Execute(c *cli.Context) error {
	action.serverEntries = []*zeroconf.ServiceEntry{}
	if err := action.PrepareExecution(c, 2, "*", "."); err != nil {
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

	action.suppressablePrintf("receive from %s@%s:%d\n", action.CliContext.String("username"), action.CliParameters.Host, action.CliParameters.Port)

	if action.CliContext.String("password") == "" {
		password, err := action.promptPassword("\nEnter password:")
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

	waitTime := 1

	// Discover all services on the network (e.g. _workstation._tcp)
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err.Error())
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			action.serverEntries = append(action.serverEntries, entry)
			println("found new server: " + fmt.Sprintf("%s:%d", entry.HostName, entry.Port))
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
	var selectedServer *zeroconf.ServiceEntry
	if serverCount == 1 {
		selectedServer = action.serverEntries[0]
	} else {
		action.suppressablePrintf("found %d servers, choose the one to receive from:\n", serverCount)

		for i := 0; i < serverCount; i++ {
			fmt.Printf("%d.) %s:%d\n", i+1, action.serverEntries[i].HostName, action.serverEntries[i].Port)
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

	action.suppressablePrintf("selected server %s:%d\n", selectedServer.HostName, selectedServer.Port)

	addr, err := net.LookupIP(selectedServer.HostName)
	if err != nil {
		log.Printf("Could not lookup host %s\n", selectedServer.HostName)
		action.CliParameters.Host = selectedServer.HostName
	} else {
		lookupSuccess := false
		for _, ip := range addr {
			if resolvedIp := action.resolveIpConnection(ip, selectedServer.Port); resolvedIp != "" {
				action.CliParameters.Host = resolvedIp
				lookupSuccess = true
				break
			}
		}

		if !lookupSuccess {
			log.Printf("Initial lookup for host %s failed\n", selectedServer.HostName)
			for _, ip := range selectedServer.AddrIPv4 {
				if resolvedIp := action.resolveIpConnection(ip, selectedServer.Port); resolvedIp != "" {
					action.CliParameters.Host = resolvedIp
					lookupSuccess = true
					break
				}
			}
		}
		if !lookupSuccess {
			log.Printf("IPv4 lookup for host %s failed\n", selectedServer.HostName)
			for _, ip := range selectedServer.AddrIPv6 {
				if resolvedIp := action.resolveIpConnection(ip, selectedServer.Port); resolvedIp != "" {
					action.CliParameters.Host = resolvedIp
					lookupSuccess = true
					break
				}
			}
		}

		if lookupSuccess {
			log.Printf("Lookup for host %s successful: %s\n", selectedServer.HostName, action.CliParameters.Host)
		} else {
			log.Printf("Lookup for host %s failed", selectedServer.HostName)
		}
	}

	action.CliParameters.Port = selectedServer.Port

	action.suppressablePrintf("connecting to %s:%d\n", action.CliParameters.Host, selectedServer.Port)

	action.receive()
	return nil

}

func (action *ReceiveAction) resolveIpConnection(netIp net.IP, port int) string {
	ip := netIp.To4()
	if ip == nil {
		ip = netIp.To16()
	}
	if ip == nil {
		return ""
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip.String(), port))
	if err != nil {
		return ""
	}
	defer conn.Close()
	return ip.String()
}
