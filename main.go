package main

import (
	"fmt"
	"github.com/sandreas/afero"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {

	app := &cli.App{
		Name:  "walk",
		Usage: "walk",
		Action: func(c *cli.Context) error {
			path := c.Args().Get(1)
			fmt.Printf("walk over <%s>", path)
			fs := afero.NewOsFs()
			err := afero.Walk(fs, path, func(path string, info os.FileInfo, err error) error {
				fmt.Println(path)
				return nil
			})

			if err != nil {
				fmt.Printf("error: %s", err)
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
