package action

import (
	"github.com/urfave/cli"
	"log"
	"io/ioutil"
	"os"
	"io"
	"os/user"
	"runtime"
	"strings"
	"errors"
)

const (
	ErrorPreventUsingSingleQuotesOnWindows = 1
	ErrorPositionalArgumentCount

)


func NewActionFactory(action string) CliActionInterface {
	switch action {
	case "find":
		return new(FindAction)
	}
	return nil
}


type GlobalArguments struct {
	Quiet         bool `arg:"help:do not show any output"`
	Force         bool `arg:"help:force the requested action - even if it might be not a good idea"`
	Debug         bool `arg:"-d,help:debug mode with logging to Stdout and into $HOME/.graft/application.log"`
	Regex         bool `arg:"help:use a real regex instead of glob patterns (e.g. src/.*\\.jpg)"`
	CaseSensitive bool `arg:"--case-sensitive,help:be case sensitive when matching files and folders"`
	MaxAge        string `arg:"--max-age,help:maximum age (e.g. 2d / 8w / 2016-12-24 / etc. )"`
	MinAge        string `arg:"--min-age,help:minimum age (e.g. 2d / 8w / 2016-12-24 / etc. )"`
	MaxSize       string `arg:"--max-size,help:maximum size in bytes or format string (e.g. 2G / 8M / 1000K etc. )"`
	MinSize       string `arg:"--min-size,help:minimum size in bytes or format string (e.g. 2G / 8M / 1000K etc. )"`
	ExportTo      string `arg:"--export-to,help:export found matches to a text file - one line per item"`
	FilesFrom     string `arg:"--files-from,help:import found matches from file - one line per item"`
}

type AbstractAction struct {
	Arguments *GlobalArguments
	Context *cli.Context
}

func (act *AbstractAction) PrepareExecution(c *cli.Context, positionalArgumentsCount int) error {
	act.ReadGlobalArguments(c)
	act.initLogging()
	if act.usedSingleQuotesAsQualifierOnWindows() {
		return cli.NewExitError("using single quotes as qualifier may lead to unexpected results - please use double quotes or --force", ErrorPreventUsingSingleQuotesOnWindows)
	}

	if err := act.assertPositionalArgumentsCount(positionalArgumentsCount); err != nil {
		return cli.NewExitError(err.Error(), ErrorPositionalArgumentCount)
	}

	return nil
}
func (act *AbstractAction) assertPositionalArgumentsCount(positionalArgumentsCount int) error {
	if len(act.Context.Args()) != 1 {
		return errors.New("find takes exactly one argument as search pattern")
	}
	return nil
}

func (act *AbstractAction) ReadGlobalArguments(c *cli.Context) {
	act.Arguments = &GlobalArguments{
		Debug: c.Bool("debug"),
		FilesFrom: c.String("files-from"),
		ExportTo: c.String("export-to"),
		MinAge: c.String("min-age"),
		MaxAge: c.String("max-age"),
		MinSize: c.String("min-size"),
		MaxSize: c.String("min-size"),
	}


}

func (act *AbstractAction) initLogging() {
	if !act.Arguments.Debug {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
		return
	}
	log.SetOutput(os.Stdout)

	homeDir, err := act.createHomeDirectoryIfNotExists()
	if err != nil {
		log.Println("could not create home directory: ", homeDir, err)
	}
	logFileName := homeDir + "/graft.log"
	os.Remove(logFileName)
	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println("could not open logfile: ", logFile, err)
	}
	defer logFile.Close()
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}

func (act *AbstractAction) createHomeDirectoryIfNotExists() (string, error) {
	u, _ := user.Current()
	homeDir := u.HomeDir + "/.graft"
	if _, err := os.Stat(homeDir); err != nil {
		if err := os.Mkdir(homeDir, os.FileMode(0755)); err != nil {
			return homeDir, err
		}
	}
	return homeDir, nil
}

func (act *AbstractAction) usedSingleQuotesAsQualifierOnWindows() bool {
	if runtime.GOOS != "windows" {
		return false
	}
	for _, arg := range act.Context.Args() {
		if strings.HasPrefix(arg, "'") {
			return true
		}
	}
	return false
}
