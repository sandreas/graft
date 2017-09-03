package action

import (
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/oleksandr/bonjour"
	"github.com/sandreas/graft/apputils"
	"github.com/sandreas/graft/sftpd"
	"github.com/urfave/cli"
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
func (action *ServeAction) ServeFoundFiles() error {
	var err error
	var homeDir string
	var fi os.FileInfo

	if len(action.locator.SourceFiles) == 0 && !action.CliGlobalParameters.Force {
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

	wg := &sync.WaitGroup{}
	//var sftpListener *net.Listener
	//var bonjourListener *bonjour.Server
	wg.Add(1)
	go func() {
		action.suppressablePrintf("Running sftp server, login as %s@%s:%d\nPress CTRL+C to stop\n", username, outboundIp, port)
		// sftpListener, err := sftpd.NewSimpleSftpServer(homeDir, listenAddress, port, username, password, pathMapper)
		_, err := sftpd.NewSimpleSftpServer(homeDir, listenAddress, port, username, password, pathMapper)
		if err != nil {
			log.Printf("Error starting sftp server: " + err.Error())
		}
		wg.Done()
	}()

	if !action.CliContext.Bool("silent") {
		wg.Add(1)
		go func() {
			action.suppressablePrintf("Publishing service via mdns: active\n")

			// Run registration (blocking call)
			//bonjourListener, err := bonjour.Register("graft", "_graft._tcp", "", 9999, []string{"txtv=1", "app=graft"}, nil)
			_, err := bonjour.Register("graft", "_graft._tcp", "", 9999, []string{"txtv=1", "app=graft"}, nil)

			if err != nil {
				log.Printf("Error starting mdns: %v", err.Error())
			}
			wg.Done()
		}()

	}

	wg.Wait()

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
