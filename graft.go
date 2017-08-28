package main

import (
	"os"
	"github.com/urfave/cli"
	"github.com/sandreas/graft/action"
)

func main() {

	globalFlags := []cli.Flag{
		cli.BoolFlag{Name: "quiet, q", Usage: "do not show any output"},                                           // does quiet make sense in find?
		cli.BoolFlag{Name: "force, f", Usage: "force the requested action - even if it might be not a good idea"}, // does force make sense in find?
		cli.BoolFlag{Name: "debug", Usage: "debug mode with logging to Stdout and into $HOME/.graft/application.log"},
		cli.BoolFlag{Name: "regex", Usage: "use a real regex instead of glob patterns (e.g. src/.*\\.jpg)"},
		cli.BoolFlag{Name: "case-sensitive", Usage: "be case sensitive when matching files and folders"},
		cli.StringFlag{Name: "max-age", Usage: "maximum age (e.g. 2d / 8w / 2016-12-24 / etc. )"},
		cli.StringFlag{Name: "min-age", Usage: "minimum age (e.g. 2d / 8w / 2016-12-24 / etc. )"},
		cli.StringFlag{Name: "max-size", Usage: "maximum size in bytes or format string (e.g. 2G / 8M / 1000K etc. )"},
		cli.StringFlag{Name: "min-size", Usage: "minimum size in bytes or format string (e.g. 2G / 8M / 1000K etc. )"},
		cli.StringFlag{Name: "export-to", Usage: "export found matches to a text file - one line per item (can also be used as save cache for large scans)"},
		cli.StringFlag{Name: "files-from", Usage: "import found matches from file - one line per item (can also be used as load cache for large scans)"},
	}

	app := cli.NewApp()
	app.Name = "graft"
	app.Version = "0.2"
	app.Usage = "find, copy and serve files"
	app.Commands = []cli.Command{
		{
			Name: "find", Aliases: []string{"f"}, Action: action.NewActionFactory("find").Execute,
			Usage: "find files",
			Flags: append(globalFlags, []cli.Flag{
				cli.BoolFlag{Name: "hide-matches", Usage: "do not show matches for search pattern ($1=filename)"},
			}...),
		},
	}

	app.Run(os.Args)
}

