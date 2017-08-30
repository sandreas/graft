package action

import (
	"github.com/urfave/cli"
	"log"
)

type FindAction struct {
	AbstractAction
}

func (act *FindAction) Execute(c *cli.Context) error {
	act.PrepareExecution(c, 1)
	log.Printf("find")
	// act.LocateSourceFiles()
	// act.WalkSourceFiles(func ...)
	// act.DisconnectSourceFileSystem()
	return nil
}
