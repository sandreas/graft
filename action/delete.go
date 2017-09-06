package action

import (
	"log"
	"github.com/urfave/cli"
	"os"
)

type DeleteAction struct {
	AbstractAction
}

func (action *DeleteAction) Execute(c *cli.Context) error {
	log.Printf("delete")

	if err := action.PrepareExecution(c, 1); err != nil  {
		return cli.NewExitError(err.Error(), ErrorLocateSourceFiles)
	}
	if err := action.locateSourceFiles(); err != nil {
		return cli.NewExitError(err.Error(), ErrorLocateSourceFiles)
	}

	if err := action.DeleteFiles(); err != nil {
		return cli.NewExitError(err.Error(), ErrorDeleteFiles)
	}
	return nil
}

func (action *DeleteAction) DeleteFiles() error {
	var dirsToRemove = []string{}

	for _, path := range action.locator.SourceFiles {
		action.suppressablePrintf(path + "\n")
		// delete
		if !action.CliContext.Bool("dry-run") {
			stat, err := action.sourcePattern.Fs.Stat(path)
			if !os.IsNotExist(err) {
				if stat.Mode().IsRegular() {
					if err := action.sourcePattern.Fs.Remove(path); err != nil  {
						log.Printf("File %s could not be deleted: %s", path, err.Error())
					}
				} else if stat.Mode().IsDir() {
					dirsToRemove = append(dirsToRemove, path)
				}
			}
		}
	}


	for _, path := range dirsToRemove {
		if err := action.sourcePattern.Fs.Remove(path); err != nil  {
			log.Printf("Directory %s could not be deleted: %s", path, err.Error())
		}
	}
	return nil
}
