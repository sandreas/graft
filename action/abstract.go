package action

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/howeyc/gopass"
	"github.com/sandreas/graft/bitflag"
	"github.com/sandreas/graft/file"
	"github.com/sandreas/graft/filesystem"
	"github.com/sandreas/graft/matcher"
	"github.com/sandreas/graft/pattern"
	"github.com/spf13/afero"
	"github.com/urfave/cli"
)

const (
	ErrorPreventUsingSingleQuotesOnWindows = 1
	ErrorPositionalArgumentCount           = 2
	ErrorLocateSourceFiles                 = 3
	ErrorStartingServer                    = 4
	ErrorCopyFiles                         = 5
	ErrorMoveFiles                         = 6
	ErrorPrepareDestination                = 7
	//ErrorNoGraftServerAvailable            = 8
	//ErrorFailedToInitializeResolver        = 9
	ErrorDeleteFiles                       = 10
)

func NewActionFactory(action string) CliActionInterface {
	switch action {
	case "find":
		return new(FindAction)
	case "serve":
		return new(ServeAction)
	case "copy":
		return new(CopyAction)
	case "move":
		return new(MoveAction)
	case "receive":
		return new(ReceiveAction)
	case "delete":
		return new(DeleteAction)
	}

	return nil
}

type CliParameters struct {
	Quiet         bool   `arg:"help:do not show any output"`
	Force         bool   `arg:"help:force the requested action - even if it might be not a good idea"`
	Debug         bool   `arg:"-d,help:debug mode with logging to Stdout and into $HOME/.graft/application.log"`
	Regex         bool   `arg:"help:use a real regex instead of glob patterns (e.g. src/.*\\.jpg)"`
	CaseSensitive bool   `arg:"--case-sensitive,help:be case sensitive when matching files and folders"`
	MaxAge        string `arg:"--max-age,help:maximum age (e.g. 2d / 8w / 2016-12-24 / etc. )"`
	MinAge        string `arg:"--min-age,help:minimum age (e.g. 2d / 8w / 2016-12-24 / etc. )"`
	MaxSize       string `arg:"--max-size,help:maximum size in bytes or format string (e.g. 2G / 8M / 1000K etc. )"`
	MinSize       string `arg:"--min-size,help:minimum size in bytes or format string (e.g. 2G / 8M / 1000K etc. )"`
	ExportTo      string `arg:"--export-to,help:export found matches to a text file - one line per item"`
	FilesFrom     string `arg:"--files-from,help:import found matches from file - one line per item"`

	Client bool
	Host   string
	Port   int
}

type AbstractAction struct {
	CliParameters *CliParameters
	CliContext    *cli.Context

	PositionalArguments cli.Args
	sourceFs            afero.Fs
	sourcePattern       *pattern.SourcePattern
	compiledRegex       *regexp.Regexp
	locator             *file.Locator
}

func (action *AbstractAction) PrepareExecution(c *cli.Context, positionalArgumentsCount int, positionalDefaultsIfUnset ...string) error {

	action.ParseCliContext(c)
	action.initLogging()

	if action.usedSingleQuotesAsQualifierOnWindows() {
		return cli.NewExitError("using single quotes as qualifier may lead to unexpected results - please use double quotes or --force", ErrorPreventUsingSingleQuotesOnWindows)
	}

	if err := action.assertPositionalArgumentsCount(positionalArgumentsCount, positionalDefaultsIfUnset); err != nil {
		return cli.NewExitError(err.Error(), ErrorPositionalArgumentCount)
	}

	return nil
}
func (action *AbstractAction) assertPositionalArgumentsCount(expectedPositionalCount int, defaults []string) error {

	givenPositionalCount := len(action.CliContext.Args())

	var positionalStrings []string
	if givenPositionalCount != expectedPositionalCount {
		if len(defaults) == expectedPositionalCount {
			for i := 0; i < expectedPositionalCount; i++ {
				if i < givenPositionalCount {
					positionalStrings = append(positionalStrings, action.CliContext.Args().Get(i))
				} else {
					positionalStrings = append(positionalStrings, defaults[i])
				}
			}
			action.PositionalArguments = cli.Args(positionalStrings)
			return nil
		}

		suffix := "argument"
		if expectedPositionalCount != 1 {
			suffix += "s"
		}
		return errors.New(action.CliContext.Command.Name + " takes " + strconv.Itoa(expectedPositionalCount) + " " + suffix)
	} else {
		action.PositionalArguments = action.CliContext.Args()
	}
	return nil
}

func (action *AbstractAction) ParseCliContext(c *cli.Context) {
	action.CliContext = c
	action.CliParameters = &CliParameters{
		Debug:     c.Bool("debug"),
		Quiet:     c.Bool("quiet"),
		FilesFrom: c.String("files-from"),
		ExportTo:  c.String("export-to"),
		MinAge:    c.String("min-age"),
		MaxAge:    c.String("max-age"),
		MinSize:   c.String("min-size"),
		MaxSize:   c.String("min-size"),
		Client:    c.IsSet("client") && c.Bool("client"),
	}

	for _,name := range c.FlagNames() {
		if name == "host" {
			action.CliParameters.Host = c.String("host")
		}
		if name == "port" {
			action.CliParameters.Port = c.Int("port")
		}
	}
}

func (action *AbstractAction) initLogging() {
	if !action.CliParameters.Debug {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
		return
	}
	log.SetOutput(os.Stdout)

	homeDir, err := action.createHomeDirectoryIfNotExists()
	if err != nil {
		log.Println("could not create home directory: ", homeDir, err)
		return
	}
	logFileName := homeDir + "/graft.log"
	os.Remove(logFileName)
	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println("could not open logfile: ", logFile, err)
		return
	}
	defer logFile.Close()
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}

func (action *AbstractAction) createHomeDirectoryIfNotExists() (string, error) {
	u, _ := user.Current()
	homeDir := u.HomeDir + "/.graft"
	if _, err := os.Stat(homeDir); err != nil {
		if err := os.Mkdir(homeDir, os.FileMode(0755)); err != nil {
			return homeDir, err
		}
	}
	return homeDir, nil
}

func (action *AbstractAction) usedSingleQuotesAsQualifierOnWindows() bool {
	if runtime.GOOS != "windows" {
		return false
	}
	for _, arg := range action.CliContext.Args() {
		if strings.HasPrefix(arg, "'") {
			return true
		}
	}
	return false
}

func (action *AbstractAction) locateSourceFiles() error {
	if err := action.prepareSourcePattern(); err != nil {
		return err
	}

	if err := action.prepareLocator(); err != nil {
		return err
	}

	return nil
}

func (action *AbstractAction) prepareSourcePattern() error {
	var err error
	if err = action.prepareSourceFileSystem(); err != nil {
		return err
	}
	action.sourcePattern = pattern.NewSourcePattern(action.sourceFs, action.PositionalArguments.First(), action.parseSourcePatternBitFlags())
	return nil
}

func (action *AbstractAction) prepareSourceFileSystem() error {
	var err error
	if action.CliParameters.Client {
		host := action.CliParameters.Host
		port := action.CliParameters.Port
		username := action.CliContext.String("username")
		password := action.CliContext.String("password")
		action.sourceFs, err = filesystem.NewSftpFs(host, port, username, password)
		return err
	}
	action.sourceFs = filesystem.NewOsFs()
	return nil
}

func (action *AbstractAction) parseSourcePatternBitFlags() bitflag.Flag {
	var patternFlags bitflag.Flag
	if action.CliParameters.CaseSensitive {
		patternFlags |= pattern.CASE_SENSITIVE
	}
	if action.CliParameters.Regex {
		patternFlags |= pattern.USE_REAL_REGEX
	}
	return patternFlags
}

func (action *AbstractAction) prepareLocator() error {
	var err error
	locator := file.NewLocator(action.sourcePattern)
	locator.RegisterObserver(file.NewWalkObserver(action.suppressablePrintf))

	if action.compiledRegex, err = action.sourcePattern.Compile(); err != nil {
		return errors.New(`Could not compile source pattern - did you quote all special chars with backslash? (special chars: \ . + * ? ( ) | [ ] { } ^ $, error: ` + err.Error() + `)`)
	}

	if action.CliParameters.FilesFrom != "" {
		locatorCache := file.NewLocatorCache(action.CliParameters.FilesFrom)
		if err = locatorCache.Load(); err != nil {
			return err
		}
		locator.SourceFiles = locatorCache.Items
	} else {
		compositeMatcher := matcher.NewCompositeMatcher()
		compositeMatcher.Add(matcher.NewRegexMatcher(action.compiledRegex))

		minAge := time.Time{}
		maxAge := time.Time{}

		if action.CliParameters.MinAge != "" {
			if minAge, err = pattern.StrToAge(action.CliParameters.MinAge, time.Now()); err != nil {
				return err
			}
		}

		if action.CliParameters.MaxAge != "" {
			if maxAge, err = pattern.StrToAge(action.CliParameters.MaxAge, time.Now()); err != nil {
				return err
			}
		}

		if !minAge.IsZero() || !maxAge.IsZero() {
			compositeMatcher.Add(matcher.NewFileAgeMatcher(minAge, maxAge))
		}

		minSize := int64(-1)
		maxSize := int64(-1)
		if action.CliParameters.MinSize != "" {
			if minSize, err = pattern.StrToSize(action.CliParameters.MinSize); err != nil {
				return err
			}
		}

		if action.CliParameters.MaxSize != "" {
			if maxSize, err = pattern.StrToSize(action.CliParameters.MaxSize); err != nil {

			}
		}

		if minSize > -1 || maxSize > -1 {
			compositeMatcher.Add(matcher.NewFileSizeMatcher(minSize, maxSize))
		}

		locator.Find(compositeMatcher)
		if action.CliParameters.ExportTo != "" {
			locatorCache := file.NewLocatorCache(action.CliParameters.ExportTo)
			locatorCache.Items = locator.SourceFiles
			if err = locatorCache.Save(); err != nil {
				return err
			}
		}
	}

	action.locator = locator

	return nil
}

func (action *AbstractAction) suppressablePrintf(format string, a ...interface{}) (n int, err error) {
	if !action.CliParameters.Quiet {
		return fmt.Printf(format, a...)
	}
	log.Printf(format, a...)
	return 0, nil
}

func (action *AbstractAction) ShowMatchesForPath(path string) {
	elementMatches := pattern.BuildMatchList(action.compiledRegex, path)
	for i := 0; i < len(elementMatches); i++ {
		action.suppressablePrintf("    $" + strconv.Itoa(i+1) + ": " + elementMatches[i] + "\n")
	}
}

func (action *AbstractAction) promptPassword(message string) (string, error) {
	if message != "" {
		println(message)
	}

	if pass, err := gopass.GetPasswd(); err != nil {
		return "", err
	} else {
		return string(pass), nil
	}

}
