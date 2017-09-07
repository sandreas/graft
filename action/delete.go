package action

import (
	"log"
	"github.com/urfave/cli"
	"os"
	"github.com/sandreas/graft/filesystem"
	"bufio"
	"fmt"
	"strings"
	"errors"
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
	dryRun := action.CliContext.Bool("dry-run")
	fileCount := len(action.locator.SourceFiles)
	if !dryRun && fileCount > 0 && !action.CliParameters.Quiet && !action.CliParameters.Force {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("%d files will be deleted. proceed (y/N)?:", fileCount)
		text, _ := reader.ReadString('\n')

		if strings.ToLower(strings.TrimSpace(text)) != "y" {
			return errors.New("Deletion aborted by user")
		}
	}

	for _, path := range action.locator.SourceFiles {
		action.suppressablePrintf(path + "\n")
		// delete
		if !action.CliContext.Bool("dry-run") {
			absPath,err  := filesystem.ToAbsIfWindowsOsFs(action.sourcePattern.Fs, path)
			if err != nil {
				log.Printf("File %s could not be converted to absolute path: %s", path, err.Error())
			}
			stat, err := action.sourcePattern.Fs.Stat(absPath)
			if !os.IsNotExist(err) {
				if stat.Mode().IsRegular() {
					if err := action.sourcePattern.Fs.Remove(absPath); err != nil  {
						log.Printf("File %s could not be deleted: %s", absPath, err.Error())
					}
				} else if stat.Mode().IsDir() {
					dirsToRemove = append(dirsToRemove, absPath)
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
