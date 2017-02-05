package main

import (
	"os"
	"github.com/urfave/cli"
	"fmt"
	"errors"
)

func main() {
	app := cli.NewApp()
	app.Name = "graft"
	app.Usage = "find and copy files via command line"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "increase verbosity",
		},
	}

	app.Action = mainAction

	app.Run(os.Args)

}


func mainAction(c *cli.Context) error {
	sourcePattern := ""
	if c.NArg() < 1 {
		return errors.New("missing required parameter source-pattern, use --help parameter for usage instructions")
	}

	sourcePattern = c.Args().Get(0)
	destinationPattern := ""
	if c.NArg() > 1 {
		destinationPattern = c.Args().Get(1)
	}


	fmt.Println("sourcePattern: ", sourcePattern)
	fmt.Println("destinationPattern: ", destinationPattern)

	return nil
}