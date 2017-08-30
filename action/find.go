package action

import (
	"github.com/urfave/cli"
	"log"
)

const (
	ErrorLocateSourceFiles = 1
)

type FindAction struct {
	AbstractAction
}

func (action *FindAction) Execute(c *cli.Context) error {
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
		action.suppressablePrintf("\nNo matches found!")
		return
	}

	hideMatches := action.CliContext.Bool("hide-matches")
	for _, path := range action.locator.SourceFiles {
		// todo: Is quiet useful here?
		action.suppressablePrintf(path + "\n")
		if !hideMatches {
			action.ShowMatchesForPath(path)
		}
	}
}
