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
	"net"
	"github.com/alexflint/go-arg"
	"github.com/sandreas/graft/action"
	"github.com/sandreas/graft/bitflag"
	"github.com/sandreas/graft/file"
	"github.com/sandreas/graft/matcher"
	"github.com/sandreas/graft/pattern"
	"github.com/sandreas/graft/sftpd"
	"github.com/sandreas/graft/transfer"
	"github.com/howeyc/gopass"
)

// TODO:
// fix graft *.wav G:/wav_backup
// xxx.wav => G:/wav_backupxx.wav
// Progressbar does not update
// fix concurrency problem with pathMapper
// => use default array implementation?
// - --verbose (ls -lah like output)
// - --files-only / --directories-only
// - javascript plugins? https://github.com/robertkrimen/otto
// - --hide-progress (for working like find)
// - copy strategy:  ResumeSkipDifferent=default, ResumeReplaceDifferent (ReplaceAll, ReplaceExisting, SkipExisting)
// - compare-strategy: quick, hash, full
// - improve progress-bar output (progress speed is not accurate enough)
// - sftp-server:
// 	    filezilla takes long and produces 0 byte files
// 		filesystem watcher for sftp server (https://godoc.org/github.com/fsnotify/fsnotify)
//		accept connections from specific ip: 		conn, e := listener.Accept() clientAddr := conn.RemoteAddr() if clientAddr
// - sftp client
// - --max-depth parameter (?)
// - limit-results when searching or moving
// - Input / Colors: https://github.com/dixonwille/wlog

const (
	ERROR_PARSING_SOURCE_PATTERN        = 2
	ERROR_PARSING_MIN_AGE               = 3
	ERROR_LOADING_FILES_FROM            = 4
	ERROR_EXPORT_TO                     = 5
	ERROR_CREATE_HOME_DIR               = 6
	ERROR_STAT_SOURCE_PATTERN_PATH      = 7
	ERROR_PREVENT_USING_SINGLE_QUOTES   = 8
	ERROR_SOURCE_PATTERN_SEEMS_UNWANTED = 9
	ERROR_READING_PASSWORD_FROM_INPUT   = 10
	ERROR_PARSING_MIN_SIZE              = 11
	ERROR_PARSING_MAX_SIZE              = 12
)

type PositionalArguments struct {
	Source      string `arg:"positional,required"`
	Destination string `arg:"positional"`
}

type TransferArguments struct {
	Move   bool `arg:"help:rename files instead of copy"`
	Delete bool `arg:"help:delete found files (be careful with this one - use --dry-run before execution)"`
	DryRun bool `arg:"--dry-run,help:simulation mode - shows output but files remain unaffected"`
	Times  bool `arg:"help:transfer source modify times to destination"`
	Force  bool `arg:"help:force the requested action - even if it might be not a good idea"`
}

type FilterArguments struct {
	MaxAge        string `arg:"--max-age,help:maximum age (e.g. 2d / 8w / 2016-12-24 / etc. )"`
	MinAge        string `arg:"--min-age,help:minimum age (e.g. 2d / 8w / 2016-12-24 / etc. )"`
	MaxSize       string `arg:"--max-size,help:maximum size in bytes or format string (e.g. 2G / 8M / 1000K etc. )"`
	MinSize       string `arg:"--min-size,help:minimum size in bytes or format string (e.g. 2G / 8M / 1000K etc. )"`
	Regex         bool `arg:"help:use a real regex instead of glob patterns (e.g. src/.*\\.jpg)"`
	CaseSensitive bool `arg:"--case-sensitive,help:be case sensitive when matching files and folders"`
}

type DisplayArguments struct {
	Debug       bool `arg:"-d,help:debug mode with logging to Stdout and into $HOME/.graft/application.log"`
	Quiet       bool `arg:"help:do not show any output"`
	ShowMatches bool `arg:"--show-matches,help:show pattern matches for each found file"`
}

type ImExportArguments struct {
	ExportTo  string `arg:"--export-to,help:export found matches to a text file - one line per item"`
	FilesFrom string `arg:"--files-from,help:import found matches from file - one line per item"`
}

type SftpArguments struct {
	Sftpd        bool `arg:"help:start sftp server providing only matching files and directories"`
	SftpPassword string `arg:"--sftp-password,help:Specify the password for the sftp server"`
	SftpUsername string `arg:"--sftp-username,help:Specify the username for the sftp server"`
	SftpPort     int `arg:"--sftp-port,help:Specifies the port on which the server listens for connections"`
}

var args struct {
	PositionalArguments
	TransferArguments
	DisplayArguments
	ImExportArguments
	FilterArguments
	SftpArguments
}

func main() {
	var err error

	if len(os.Args) == 1 {
		os.Args = append(os.Args, "--help")
	}

	args.SftpPort = 2022
	args.SftpUsername = "graft"
	args.SftpPassword = ""
	arg.MustParse(&args)

	args.SftpPassword = strings.TrimSpace(args.SftpPassword)

	if args.Sftpd && args.SftpPassword == "" {
		println("Enter password for sftp-server:")
		pass, err := gopass.GetPasswd()
		exitOnError(ERROR_READING_PASSWORD_FROM_INPUT, err)
		args.SftpPassword = string(pass)
	}

	initLogging()

	if runtime.GOOS == "windows" && (strings.HasPrefix(args.Source, "'") || strings.HasPrefix(args.Destination, "'")) {
		exitOnError(ERROR_PREVENT_USING_SINGLE_QUOTES, errors.New("prevent using single quotes as qualifier on windows - it can lead to unexpected results"))
	}

	sourcePattern := pattern.NewSourcePattern(args.Source, parseSourcePatternBitFlags())
	log.Printf("SourcePattern: %+v", sourcePattern)
	compiledRegex, err := sourcePattern.Compile()
	log.Printf("compiledRegex: %s", compiledRegex)
	exitOnError(ERROR_PARSING_SOURCE_PATTERN, err)

	if (sourcePattern.Path == "." || sourcePattern.Path == "/") && strings.Contains(sourcePattern.Pattern, "/") && !args.Force {
		exitOnError(ERROR_SOURCE_PATTERN_SEEMS_UNWANTED, errors.New("Your search might be incorrect, because parts of '" + args.Source+ "' do not exist.\n\nAll subdirectories in '"+sourcePattern.Path+"' will be recursively scanned for pattern '"+sourcePattern.Pattern+"'.\n\nIf this is really what you would like to do, use --force option"))
	}

	locator := file.NewLocator(*sourcePattern)
	locator.RegisterObserver(file.NewWalkObserver(suppressablePrintf))

	if args.FilesFrom != "" {
		locatorCache := file.NewLocatorCache(args.FilesFrom)
		err := locatorCache.Load()
		exitOnError(ERROR_LOADING_FILES_FROM, err)
		locator.SourceFiles = locatorCache.Items
	} else {
		compositeMatcher := matcher.NewCompositeMatcher()
		compositeMatcher.Add(matcher.NewRegexMatcher(*compiledRegex))

		if args.MinAge != "" {
			minAge, err := pattern.StrToAge(args.MinAge, time.Now())
			exitOnError(ERROR_PARSING_MIN_AGE, err)
			compositeMatcher.Add(matcher.NewMinAgeMatcher(minAge))
		}

		if args.MaxAge != "" {
			maxAge, err := pattern.StrToAge(args.MaxAge, time.Now())
			exitOnError(ERROR_PARSING_MIN_AGE, err)
			compositeMatcher.Add(matcher.NewMaxAgeMatcher(maxAge))
		}

		minSize := int64(-1)
		maxSize := int64(-1)
		if args.MinSize != "" {
			minSize, err = pattern.StrToSize(args.MinSize)
			exitOnError(ERROR_PARSING_MIN_SIZE, err)
		}

		if args.MaxSize != "" {
			maxSize, err = pattern.StrToSize(args.MaxSize)
			exitOnError(ERROR_PARSING_MAX_SIZE, err)
		}

		if minSize > -1 || maxSize > -1 {
			compositeMatcher.Add(matcher.NewSizeMatcher(minSize, maxSize))
		}

		locator.Find(compositeMatcher)
		if args.ExportTo != "" {
			locatorCache := file.NewLocatorCache(args.ExportTo)
			locatorCache.Items = locator.SourceFiles
			err := locatorCache.Save()
			exitOnError(ERROR_EXPORT_TO, err)
		}
	}

	if len(locator.SourceFiles) == 0 {
		suppressablePrintf("\nNo matches found!")
	}

	if args.Destination == "" {

		if args.Sftpd {
			homeDir, err := createHomeDirectoryIfNotExists()
			exitOnError(ERROR_CREATE_HOME_DIR, err)

			fi, err := os.Stat(sourcePattern.Path)
			exitOnError(ERROR_STAT_SOURCE_PATTERN_PATH, err)
			basePath := sourcePattern.Path
			if fi.Mode().IsRegular() {
				basePath = strings.TrimSuffix(basePath, "/"+fi.Name())
				if basePath == fi.Name() {
					basePath = "."
				}
			}

			pathMapper := sftpd.NewPathMapper(locator.SourceFiles, basePath)

			listenAddress := "0.0.0.0"
			outboundIp := GetOutboundIP()
			suppressablePrintf("Running sftp server, login as %s@%s:%d\nPress CTRL+C to stop\n", args.SftpUsername, outboundIp, args.SftpPort)
			sftpd.NewSimpleSftpServer(homeDir, listenAddress, args.SftpPort, args.SftpUsername, args.SftpPassword, pathMapper)
			return
		}

		for _, path := range locator.SourceFiles {
			suppressablePrintf(path + "\n")
			if args.ShowMatches {
				elementMatches := pattern.BuildMatchList(compiledRegex, path)
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

	destinationPattern := pattern.NewDestinationPattern(args.Destination)
	messagePrinter := transfer.NewMessagePrinterObserver(suppressablePrintf)
	actionBitFlags := parseActionBitFlags()

	if args.Move {
		moveStrategy := transfer.NewMoveStrategy()

		moveAction := action.NewTransferAction(locator.SourceFiles, moveStrategy, actionBitFlags)
		moveAction.RegisterObserver(messagePrinter)
		err = moveAction.Execute(sourcePattern, destinationPattern)
	} else {
		copyStrategy := transfer.NewCopyStrategy()
		copyStrategy.ProgressHandler = transfer.NewCopyProgressHandler(int64(32*1024), 1*time.Second)
		copyStrategy.RegisterObserver(messagePrinter)

		copyAction := action.NewTransferAction(locator.SourceFiles, copyStrategy, actionBitFlags)
		copyAction.RegisterObserver(messagePrinter)
		err = copyAction.Execute(sourcePattern, destinationPattern)
	}

	if err != nil {
		suppressablePrintf(err.Error())
	}
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
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

func parseSourcePatternBitFlags() bitflag.BitFlag {
	var patternFlags bitflag.BitFlag
	if args.CaseSensitive {
		patternFlags |= pattern.CASE_SENSITIVE
	}
	if args.Regex {
		patternFlags |= pattern.USE_REAL_REGEX
	}
	return patternFlags
}

func parseActionBitFlags() bitflag.BitFlag {
	var actionFlags bitflag.BitFlag
	if args.DryRun {
		actionFlags |= action.FLAG_DRY_RUN
	}

	if args.Times {
		actionFlags |= action.FLAG_TIMES
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
