package action

import (
	"fmt"
	"log"
	"github.com/urfave/cli"
)

type FindAction struct {
	AbstractAction
}

func (action *FindAction) Execute(c *cli.Context) error {
	//fs := filesystem.NewOsFs()
	//fs.Stat("graft.go")
	//os.Exit(0)

	action.PrepareExecution(c, 1)
	log.Printf("find")
	if err := action.locateSourceFiles(); err != nil {
		return cli.NewExitError(err.Error(), ErrorLocateSourceFiles)
	}
	action.ShowFoundFiles()
	return nil
}
func (action *FindAction) ShowFoundFiles() {

	if len(action.locator.SourceFiles) == 0 {
		action.suppressablePrintf("\nNo matching files found!\n")
		return
	}

	showMatches := action.CliContext.Bool("show-matches") && !action.CliParameters.Quiet
	for _, path := range action.locator.SourceFiles {
		fmt.Println(path) // quiet does not influence the output of the file listing, since this is the only sense of this action
		if showMatches {
			action.ShowMatchesForPath(path)
		}
	}
}
