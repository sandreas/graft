package main

import (
	"os"

	"github.com/sandreas/graft/action"
	"github.com/urfave/cli"
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

	networkFlags := []cli.Flag{
		cli.StringFlag{Name: "host", Usage: "Specify the hostname for the server (client mode only)"},
		cli.IntFlag{Name: "port", Usage: "Specifiy server port (used for server- and client mode)", Value: 2022},
		cli.StringFlag{Name: "username", Usage: "Specify server username (used in server- and client mode)", Value: "graft"},
		cli.StringFlag{Name: "password", Usage: "Specify server password (used for server- and client mode)"},
	}

	findFlags := []cli.Flag{
		cli.BoolFlag{Name: "hide-matches", Usage: "do not show matches for search pattern ($1=filename)"},
		cli.BoolFlag{Name: "client", Usage: "client mode - act as sftp client and search files remotely instead of local search"},
	}

	serveFlags := []cli.Flag{
		cli.BoolFlag{Name: "silent", Usage: "do not use mdns to publish multicast sftp server (graft receive will not work without parameters)"},
	}

	dryRunFlags := []cli.Flag{
		cli.BoolFlag{Name: "dry-run", Usage: "simulation mode - shows output but files remain unaffected"},
	}

	transferFlags := []cli.Flag{
		cli.BoolFlag{Name: "times", Usage: "transfer source modify times to destination"},
	}

	app := cli.NewApp()
	app.Name = "graft"
	app.Version = "0.2"
	app.Usage = "find, copy and serve files"

	app.Commands = []cli.Command{
		{
			Name: "find", Aliases: []string{"f"}, Action: action.NewActionFactory("find").Execute,
			Usage: "find files",
			Flags:mergeFlags(globalFlags, networkFlags, findFlags),
		},
		{
			Name: "serve", Aliases: []string{"s"}, Action: action.NewActionFactory("serve").Execute,
			Usage: "serve files",
			Flags: mergeFlags(globalFlags, networkFlags, serveFlags),
		},
		{
			Name: "copy", Aliases: []string{"c", "cp"}, Action: action.NewActionFactory("copy").Execute,
			Usage: "copy files from a source to a destination",
			Flags: mergeFlags(globalFlags, transferFlags, dryRunFlags),
		},
		{
			Name: "move", Aliases: []string{"m", "mv"}, Action: action.NewActionFactory("move").Execute,
			Usage: "move files from a source to a destination",
			Flags: mergeFlags(globalFlags, dryRunFlags, transferFlags),
		},
		{
			Name: "delete", Aliases: []string{"d", "rm"}, Action: action.NewActionFactory("delete").Execute,
			Usage: "delete files recursively",
			Flags: mergeFlags(globalFlags, dryRunFlags),
		},
		{
			Name: "receive", Aliases: []string{"r"}, Action: action.NewActionFactory("receive").Execute,
			Usage: "receive files from a graft server",
			Flags: mergeFlags(globalFlags, dryRunFlags, transferFlags, networkFlags),
		},
	}

	app.Run(os.Args)
}

func mergeFlags(flagsToMerge ...[]cli.Flag) []cli.Flag {
	mergedFlags := []cli.Flag{}
	for _,flags := range flagsToMerge  {
		mergedFlags = append(mergedFlags, flags...)
	}
	return mergedFlags
}
