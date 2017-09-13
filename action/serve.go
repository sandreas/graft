package action

import (
	"log"
	"net"
	"os"
	"strings"

	"os/signal"
	"syscall"

	"strconv"

	"github.com/grandcat/zeroconf"
	"github.com/sandreas/graft/apputils"
	"github.com/sandreas/graft/sftpd"
	"github.com/urfave/cli"
	"fmt"
	"crypto/rand"
)

type ServeAction struct {
	AbstractAction
}

func (action *ServeAction) Execute(c *cli.Context) error {
	action.PrepareExecution(c, 1, "*")
	log.Printf("serve")
	if err := action.locateSourceFiles(); err != nil {
		return cli.NewExitError(err.Error(), ErrorLocateSourceFiles)
	}
	if err := action.ServeFoundFiles(); err != nil {
		return cli.NewExitError(err.Error(), ErrorStartingServer)
	}
	return nil
}

func pseudoUuid() (uuid string) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return
}

func (action *ServeAction) ServeFoundFiles() error {
	var err error
	var homeDir string
	var fi os.FileInfo

	if len(action.locator.SourceFiles) == 0 && !action.CliParameters.Force {
		action.suppressablePrintf("\nNo matching files found, server does not need to be started - use force to start server anyway\n")
		return nil
	}

	if homeDir, err = action.createHomeDirectoryIfNotExists(); err != nil {
		return err
	}

	if fi, err = action.sourceFs.Stat(action.sourcePattern.Path); err != nil {
		return err
	}
	basePath := action.sourcePattern.Path
	if fi.Mode().IsRegular() {
		basePath = strings.TrimSuffix(basePath, "/"+fi.Name())
		if basePath == fi.Name() {
			basePath = "."
		}
	}
	pathMapper := sftpd.NewPathMapper(action.locator.SourceFiles, basePath)
	listenAddress := "0.0.0.0"
	if action.CliContext.String("host") != "" {
		listenAddress = action.CliContext.String("host")
	}
	outboundIp, err := apputils.GetOutboundIpAsString("localhost", net.Dial)
	if err != nil {
		log.Printf("Error on GetOutboundIpAsString: %v", err)
	}

	username := action.CliContext.String("username")
	password := action.CliContext.String("password")
	port := action.CliContext.Int("port")

	if password == "" {
		password, err = action.promptPassword("Which password shall be used for user " + username + "?")
		if err != nil {
			return err
		}
	}

	if !action.CliContext.Bool("no-zeroconf") {
		action.suppressablePrintf("Publishing service via mdns: active\n")

		uuid := pseudoUuid()
		port := action.CliParameters.Port

		name := "graft-sftp-server-" + uuid + "_" + outboundIp + ":" + strconv.Itoa(port)
		service := "_graft._tcp"
		domain := "local."

		server, err := zeroconf.Register(name, service, domain, port, []string{"txtv=0.2", "domain=" + domain, "ip=" + outboundIp}, nil)
		if err != nil {
			panic(err)
		}
		defer server.Shutdown()
		log.Println("Published service:")
		log.Println("- Name:", name)
		log.Println("- Type:", service)
		log.Println("- Domain:", domain)
		log.Println("- Port:", port)

	}

	go func() {
		action.suppressablePrintf("Running sftp server, login as %s@%s:%d\nPress CTRL+C to stop\n", username, outboundIp, port)
		// sftpListener, err := sftpd.NewSimpleSftpServer(homeDir, listenAddress, port, username, password, pathMapper)
		_, err := sftpd.NewSimpleSftpServer(homeDir, listenAddress, port, username, password, pathMapper)
		if err != nil {
			log.Printf("Error starting sftp server: " + err.Error())
		}
	}()

	// Clean exit.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	// Timeout timer

	select {
	case <-sig:
		// Exit by user
	}

	// TODO Ctrl+C handling
	//handler := make(chan os.Signal, 1)
	//signal.Notify(handler, os.Interrupt)
	//for sig := range handler {
	//	if sig == os.Interrupt {
	//		//bonjourListener.Shutdown()
	//		//sftpListener.Close()
	//		wg.Done()
	//		time.Sleep(1e9)
	//		break
	//	}
	//}

	return nil
}
