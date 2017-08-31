package action

import (
	"github.com/urfave/cli"
	"log"
	"strings"
	"github.com/sandreas/graft/apputils"
	"net"
	"os"
	"github.com/sandreas/graft/sftpd"
)



type ServeAction struct {
	AbstractAction
}

func (action *ServeAction) Execute(c *cli.Context) error {
	action.PrepareExecution(c, 1)
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

	if len(action.locator.SourceFiles) == 0 && !action.CliGlobalParameters.Force{
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
		password, err = action.promptPassword("Which password shall be used for user " + username+ "?")
		if err != nil {
			return err
		}
	}

	action.suppressablePrintf("Running sftp server, login as %s@%s:%d\nPress CTRL+C to stop\n", username, outboundIp, port)
	sftpd.NewSimpleSftpServer(homeDir, listenAddress, port, username, password, pathMapper)
	return nil
}
