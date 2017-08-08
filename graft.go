package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"time"

	"errors"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/sandreas/graft/newaction"
	"github.com/sandreas/graft/newfile"
	"github.com/sandreas/graft/newmatcher"
	"github.com/sandreas/graft/newoptions"
	"github.com/sandreas/graft/newpattern"
	"github.com/sandreas/graft/newtransfer"
	"github.com/sandreas/graft/sftpd"
)

// TODO:
// - sftp-server
// - limit-results when searching or moving
// - graft xxx/* yyy/$1 when xxx does not exist could result in:
// 		(xxx/.*$) => yyy/$1 which may not be intended
// 		if pattern contains unmasked slash, suggest not searching, because directory does not exist
// - Input / Colors: https://github.com/dixonwille/wlog
//)

const (
	ERROR_PASSWORD_CANNOT_BE_EMPTY = 1
	ERROR_PARSING_SOURCE_PATTERN   = 2
	ERROR_PARSING_MIN_AGE          = 3
	ERROR_LOADING_FILES_FROM       = 4
	ERROR_EXPORT_TO                = 5
	ERROR_CREATE_HOME_DIR          = 6
	ERROR_STAT_SOURCE_PATTERN_PATH = 7
)

type PositionalArguments struct {
	Source      string `arg:"positional,required"`
	Destination string `arg:"positional"`
}

type BooleanFlags struct {
	CaseSensitive bool `arg:"--case-sensitive,help:be case sensitive when matching files and folders"`
	Regex         bool `arg:"help:use a real regex instead of glob patterns (e.g. src/.*\\.jpg)"`
	Debug         bool `arg:"-d,help:debug mode with logging to Stdout and into $HOME/.graft/application.log"`
	Quiet         bool `arg:"help:do not show any output"`
	ShowMatches   bool `arg:"--show-matches,help:show pattern matches for each found file"`
	DryRun        bool `arg:"--dry-run,help:simulation mode output only files remain unaffected"`
	Times         bool `arg:"help:transfer source modify times to destination"`
	Move          bool `arg:"help:move / rename files - do not make a copy"`
	Delete        bool `arg:"help:delete found files"`

	Serve bool `arg:"help:provide matching files via sftp"`
	//Verbose bool `arg:"-v,help:be verbose"`
}

type IntParameters struct {
	Port int `arg:"help:Specifies the port on which the server listens for connections"`
}

type StringParameters struct {
	MinAge    string `arg:"--min-age,help:minimum age (e.g. 2d / 8w / 2016-12-24 / etc. )"`
	MaxAge    string `arg:"--max-age,help:maximum age (e.g. 2d / 8w / 2016-12-24 / etc. )"`
	FilesFrom string `arg:"--files-from,help:import source listing from file - one line per item"`
	ExportTo  string `arg:"--export-to,help:export source listing to file - one line per item"`
	User      string `arg:"help:Specify the username for the sftp server"`
	Password  string `arg:"help:Specify the password for the sftp server"`
}

var args struct {
	PositionalArguments
	BooleanFlags
	IntParameters
	StringParameters
}

func main() {
	var err error

	if len(os.Args) == 1 {
		os.Args = append(os.Args, "--help")
	}

	args.Port = 2022
	args.User = "graft"
	args.Password = ""
	arg.MustParse(&args)

	args.Password = strings.TrimSpace(args.Password)

	if args.Serve && args.Password == "" {
		exitOnError(ERROR_PASSWORD_CANNOT_BE_EMPTY, errors.New("Password cannot be empty!"))
	}

	initLogging()

	sourcePattern := newpattern.NewSourcePattern(args.Source, parseSourcePatternBitFlags())
	log.Printf("SourcePattern: %+v", sourcePattern)
	compiledRegex, err := sourcePattern.Compile()
	log.Printf("compiledRegex: %s", compiledRegex)
	exitOnError(ERROR_PARSING_SOURCE_PATTERN, err)

	locator := newfile.NewLocator(*sourcePattern)
	locator.RegisterObserver(newfile.NewWalkObserver(suppressablePrintf))

	if args.FilesFrom != "" {
		locatorCache := newfile.NewLocatorCache(args.FilesFrom)
		err := locatorCache.Load()
		exitOnError(ERROR_LOADING_FILES_FROM, err)
		locator.SourceFiles = locatorCache.Items
	} else {
		compositeMatcher := newmatcher.NewCompositeMatcher()
		compositeMatcher.Add(newmatcher.NewRegexMatcher(*compiledRegex))

		if args.MinAge != "" {
			minAge, err := newpattern.StrToAge(args.MinAge, time.Now())
			exitOnError(ERROR_PARSING_MIN_AGE, err)
			compositeMatcher.Add(newmatcher.NewMinAgeMatcher(minAge))
		}

		if args.MaxAge != "" {
			maxAge, err := newpattern.StrToAge(args.MaxAge, time.Now())
			exitOnError(ERROR_PARSING_MIN_AGE, err)
			compositeMatcher.Add(newmatcher.NewMaxAgeMatcher(maxAge))
		}

		locator.Find(compositeMatcher)
		if args.ExportTo != "" {
			locatorCache := newfile.NewLocatorCache(args.ExportTo)
			locatorCache.Items = locator.SourceFiles
			err := locatorCache.Save()
			exitOnError(ERROR_EXPORT_TO, err)
		}
	}

	if args.Destination == "" {

		if args.Serve {
			homeDir, err := createHomeDirectoryIfNotExists()
			exitOnError(ERROR_CREATE_HOME_DIR, err)

			fi, err := os.Stat(sourcePattern.Path)
			exitOnError(ERROR_STAT_SOURCE_PATTERN_PATH, err)
			basePath := sourcePattern.Path
			if fi.Mode().IsRegular() {
				basePath = strings.TrimSuffix(basePath, "/" + fi.Name())
			}

			pathMapper := sftpd.NewPathMapper(locator.SourceFiles, basePath)

			sftpd.NewGraftServer(homeDir, "0.0.0.0", args.Port, args.User, args.Password, pathMapper)
			return
		}

		for _, path := range locator.SourceFiles {
			suppressablePrintf(path + "\n")
			if args.ShowMatches {
				elementMatches := newpattern.BuildMatchList(compiledRegex, path)
				for i := 0; i < len(elementMatches); i++ {
					suppressablePrintf("    $" + strconv.Itoa(i+1) + ": " + elementMatches[i] + "\n")
				}
			}

			// delete
			if args.Delete && !args.DryRun {
				var dirsToRemove = []string{}
				stat, err := os.Stat(path)

				if !os.IsNotExist(err) {
					if stat.Mode().IsRegular() {
						os.Remove(path)
					} else if stat.Mode().IsDir() {
						dirsToRemove = append(dirsToRemove, path)
					}
				}

				for _, path := range dirsToRemove {
					os.Remove(path)
				}
			}
		}

		return
	}

	destinationPattern := newpattern.NewDestinationPattern(args.Destination)
	messagePrinter := newtransfer.NewMessagePrinterObserver(suppressablePrintf)
	actionBitFlags := parseActionBitFlags()

	if args.Move {
		moveStrategy := newtransfer.NewMoveStrategy()

		moveAction := newaction.NewTransferAction(locator.SourceFiles, moveStrategy, actionBitFlags)
		moveAction.RegisterObserver(messagePrinter)
		err = moveAction.Execute(sourcePattern, destinationPattern)
	} else {
		copyStrategy := newtransfer.NewCopyStrategy()
		copyStrategy.ProgressHandler = newtransfer.NewCopyProgressHandler(int64(32*1024), 2*time.Second)
		copyStrategy.RegisterObserver(messagePrinter)

		copyAction := newaction.NewTransferAction(locator.SourceFiles, copyStrategy, actionBitFlags)
		copyAction.RegisterObserver(messagePrinter)
		err = copyAction.Execute(sourcePattern, destinationPattern)
	}

	if err != nil {
		suppressablePrintf(err.Error())
	}
}

func (PositionalArguments) Description() string {
	return "graft 0.2 - a command line application to search for and transfer files\n"
}

func initLogging() {
	if !args.Debug {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
		return
	}
	log.SetOutput(os.Stdout)

	homeDir, err := createHomeDirectoryIfNotExists()
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

func suppressablePrintf(format string, a ...interface{}) (n int, err error) {
	if !args.Quiet {
		return fmt.Printf(format, a...)
	}
	log.Printf(format, a...)
	return 0, nil
}

func parseSourcePatternBitFlags() newoptions.BitFlag {
	var patternFlags newoptions.BitFlag
	if args.CaseSensitive {
		patternFlags |= newpattern.CASE_SENSITIVE
	}
	if args.Regex {
		patternFlags |= newpattern.USE_REAL_REGEX
	}
	return patternFlags
}

func parseActionBitFlags() newoptions.BitFlag {
	var actionFlags newoptions.BitFlag
	if args.DryRun {
		actionFlags |= newaction.FLAG_DRY_RUN
	}

	if args.Times {
		actionFlags |= newaction.FLAG_TIMES
	}
	return actionFlags
}

func createHomeDirectoryIfNotExists() (string, error) {
	u, _ := user.Current()
	homeDir := u.HomeDir + "/.graft"
	if _, err := os.Stat(homeDir); err != nil {
		if err := os.Mkdir(homeDir, os.FileMode(0755)); err != nil {
			return homeDir, err
		}
	}
	return homeDir, nil
}

func exitOnError(exitCode int, err error) {
	if err == nil {
		return
	}

	_, fn, line, _ := runtime.Caller(1)
	if args.Debug {
		log.Printf("[error] %s:%d %v (Code: %d)", fn, line, err, exitCode)
	} else {
		suppressablePrintf(err.Error()+" (Code: %d)", exitCode)
	}
	os.Exit(exitCode)
}
