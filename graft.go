package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"github.com/urfave/cli"
	"github.com/spf13/afero"
	"fmt"
	"github.com/sandreas/graft/sftpd"
	"github.com/sandreas/graft/sftpfs"
	"github.com/howeyc/gopass"
	"strings"
	"runtime"
	"github.com/sandreas/graft/pattern"
	"github.com/sandreas/graft/bitflag"
	"github.com/sandreas/graft/file"
	"github.com/sandreas/graft/matcher"
	"regexp"
	"time"
	"strconv"
)

const (
	//ERROR_PARSING_SOURCE_PATTERN        = 2
	ERROR_PARSING_MIN_AGE               = 3
	ERROR_PARSING_MAX_AGE               = 128
	ERROR_LOADING_FILES_FROM            = 4
	//ERROR_EXPORT_TO                     = 5
	//ERROR_CREATE_HOME_DIR               = 6
	//ERROR_STAT_SOURCE_PATTERN_PATH      = 7
	ERROR_PREVENT_USING_SINGLE_QUOTES   = 8
	//ERROR_SOURCE_PATTERN_SEEMS_UNWANTED = 9
	//ERROR_READING_PASSWORD_FROM_INPUT   = 10
	ERROR_PARSING_MIN_SIZE              = 11
	ERROR_PARSING_MAX_SIZE              = 12
	//ERROR_CONNECT_TO_SERVER             = 13
	ERROR_TOO_MANY_ARGUMENTS			= 14
)


//type PositionalArguments struct {
//	Source      string `arg:"positional"`
//	Destination string `arg:"positional"`
//}
//
//type TransferArguments struct {
//	DryRun bool `arg:"--dry-run,help:simulation mode - shows output but files remain unaffected"`
//	Times  bool `arg:"help:transfer source modify times to destination"`
//}

type GlobalArguments struct {
	Quiet       bool `arg:"help:do not show any output"`
	Force  bool `arg:"help:force the requested action - even if it might be not a good idea"`
	Debug       bool `arg:"-d,help:debug mode with logging to Stdout and into $HOME/.graft/application.log"`
	Regex         bool `arg:"help:use a real regex instead of glob patterns (e.g. src/.*\\.jpg)"`
	CaseSensitive bool `arg:"--case-sensitive,help:be case sensitive when matching files and folders"`
	// ShowMatches bool `arg:"--show-matches,help:show pattern matches for each found file"`
	MaxAge        string `arg:"--max-age,help:maximum age (e.g. 2d / 8w / 2016-12-24 / etc. )"`
	MinAge        string `arg:"--min-age,help:minimum age (e.g. 2d / 8w / 2016-12-24 / etc. )"`
	MaxSize       string `arg:"--max-size,help:maximum size in bytes or format string (e.g. 2G / 8M / 1000K etc. )"`
	MinSize       string `arg:"--min-size,help:minimum size in bytes or format string (e.g. 2G / 8M / 1000K etc. )"`
	ExportTo  string `arg:"--export-to,help:export found matches to a text file - one line per item"`
	FilesFrom string `arg:"--files-from,help:import found matches from file - one line per item"`
}

//type SftpArguments struct {
//	Server   bool `arg:"help:server mode - act as sftp server and provide only files and directories matching the source pattern"`
//	Client   bool `arg:"help:client mode - act as sftp client and download files instead of local search"`
//	Host     string `arg:"help:Specify the hostname for the server (client mode only)"`
//	Username string `arg:"help:Specify server username (used in server- and client mode)"`
//	Password string `arg:"help:Specify server password (used for server- and client mode)"`
//	Port     int `arg:"help:Specifiy server port (used for server- and client mode)"`
//}

var args struct {
	GlobalArguments
}


func main() {
	app := cli.NewApp()
	app.Name = "graft"
	app.Version = "0.2"
	app.Usage = "find, copy and serve files"
	app.Flags = []cli.Flag{
		cli.BoolFlag{Name: "quiet, q", Usage: "do not show any output"},                                           // does quiet make sense in find?
		cli.BoolFlag{Name: "force, f", Usage: "force the requested action - even if it might be not a good idea"}, // does force make sense in find?
		cli.BoolFlag{Name: "debug, d", Usage: "debug mode with logging to Stdout and into $HOME/.graft/application.log"},
		cli.BoolFlag{Name: "regex", Usage: "use a real regex instead of glob patterns (e.g. src/.*\\.jpg)"},
		cli.BoolFlag{Name: "case-sensitive", Usage: "be case sensitive when matching files and folders"},
		cli.StringFlag{Name: "max-age", Usage: "maximum age (e.g. 2d / 8w / 2016-12-24 / etc. )"},
		cli.StringFlag{Name: "min-age", Usage: "minimum age (e.g. 2d / 8w / 2016-12-24 / etc. )"},
		cli.StringFlag{Name: "max-size", Usage: "maximum size in bytes or format string (e.g. 2G / 8M / 1000K etc. )"},
		cli.StringFlag{Name: "min-size", Usage: "minimum size in bytes or format string (e.g. 2G / 8M / 1000K etc. )"},
		cli.StringFlag{Name: "export-to", Usage: "export found matches to a text file - one line per item (can also be used as save cache for large scans)"},
		cli.StringFlag{Name: "files-from", Usage: "import found matches from file - one line per item (can also be used as load cache for large scans)"},
	}

	app.Commands = []cli.Command{
		{
			Name:  "find", Aliases: []string{"f"}, Action: findAction,
			Usage: "find files",
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "hide-matches", Usage: "do not show matches for search pattern ($1=filename)"},
			},
		},
	}

	app.Run(os.Args)
}

func findAction(c *cli.Context) error {
	readGlobalArguments(c)
	initLogging(c)
	if usedSingleQuotesAsQualifierOnWindows(c) {
		return cli.NewExitError("using single quotes as qualifier may lead to unexpected results - please use double quotes or --force", ERROR_PREVENT_USING_SINGLE_QUOTES)
	}

	if len(c.Args()) != 1 {
		return cli.NewExitError("find takes exactly one argument as search pattern", ERROR_TOO_MANY_ARGUMENTS)
	}

	sourcePatternString := c.Args().First()


	sourceFileSystem, ctx, err := prepareSourceFileSystem(c)
	if err != nil {
		// todo: make exit error
		return err
	} else if ctx != nil {
		defer ctx.Disconnect()
	}

	sourcePattern := pattern.NewSourcePattern(sourceFileSystem, sourcePatternString, parseSourcePatternBitFlags(c))
	log.Printf("SourcePattern: %+v", sourcePattern)
	compiledRegex, err := sourcePattern.Compile()
	log.Printf("compiledRegex: %s", compiledRegex)


	locator, err := prepareLocator(sourcePattern, compiledRegex)
	if len(locator.SourceFiles) == 0 {
		suppressablePrintf("\nNo matching files found!")
	}

	for _, path := range locator.SourceFiles {
		suppressablePrintf(path + "\n")
		if ! c.Bool("hide-matches") {
			elementMatches := pattern.BuildMatchList(compiledRegex, path)
			for i := 0; i < len(elementMatches); i++ {
				suppressablePrintf("    $" + strconv.Itoa(i+1) + ": " + elementMatches[i] + "\n")
			}
		}

	}

	return nil
}
func prepareLocator(sourcePattern *pattern.SourcePattern, compiledRegex *regexp.Regexp) (*file.Locator, error) {
	locator := file.NewLocator(sourcePattern)
	locator.RegisterObserver(file.NewWalkObserver(suppressablePrintf))

	if args.FilesFrom != "" {
		locatorCache := file.NewLocatorCache(args.FilesFrom)
		err := locatorCache.Load()
		if err != nil {
			return nil, cli.NewExitError(err, ERROR_LOADING_FILES_FROM)
		}
		locator.SourceFiles = locatorCache.Items
		return locator, nil
	}

	compositeMatcher := matcher.NewCompositeMatcher()
	compositeMatcher.Add(matcher.NewRegexMatcher(*compiledRegex))

	var err error
	minAge := time.Time{}
	maxAge := time.Time{}

	if args.MinAge != "" {
		minAge, err = pattern.StrToAge(args.MinAge, time.Now())
		if err != nil  {
			return nil, cli.NewExitError(err, ERROR_PARSING_MIN_AGE)
		}
	}

	if args.MaxAge != "" {
		maxAge, err = pattern.StrToAge(args.MaxAge, time.Now())
		if err != nil  {
			return nil, cli.NewExitError(err, ERROR_PARSING_MAX_AGE)
		}
	}

	if !minAge.IsZero() || !maxAge.IsZero() {
		compositeMatcher.Add(matcher.NewFileAgeMatcher(minAge, maxAge))
	}

	minSize := int64(-1)
	maxSize := int64(-1)
	if args.MinSize != "" {
		minSize, err = pattern.StrToSize(args.MinSize)
		if err != nil  {
			return nil, cli.NewExitError(err, ERROR_PARSING_MIN_SIZE)
		}
	}

	if args.MaxSize != "" {
		maxSize, err = pattern.StrToSize(args.MaxSize)
		if err != nil  {
			return nil, cli.NewExitError(err, ERROR_PARSING_MAX_SIZE)
		}	}

	if minSize > -1 || maxSize > -1 {
		compositeMatcher.Add(matcher.NewFileSizeMatcher(minSize, maxSize))
	}

	locator.Find(compositeMatcher)
	if args.ExportTo != "" {
		locatorCache := file.NewLocatorCache(args.ExportTo)
		locatorCache.Items = locator.SourceFiles
		err := locatorCache.Save()
		if err != nil {
			return nil, cli.NewExitError(err, ERROR_PARSING_MIN_SIZE)
		}
	}

	return locator, nil
}

func readGlobalArguments(context *cli.Context) {
	args.Debug = context.Bool("debug")
	args.FilesFrom = context.String("files-from")
	args.ExportTo = context.String("export-to")
	args.MinAge = context.String("min-age")
	args.MaxAge = context.String("max-age")
	args.MinSize = context.String("min-size")
	args.MaxSize = context.String("min-size")
}
func suppressablePrintf(format string, a ...interface{}) (n int, err error) {
	if !args.Quiet {
		return fmt.Printf(format, a...)
	}
	log.Printf(format, a...)
	return 0, nil
}
func parseSourcePatternBitFlags(c *cli.Context) bitflag.Flag {
	var patternFlags bitflag.Flag
	if c.Bool("case-sensitive") {
		patternFlags |= pattern.CASE_SENSITIVE
	}
	if c.Bool("regex") {
		patternFlags |= pattern.USE_REAL_REGEX
	}
	return patternFlags
}

func usedSingleQuotesAsQualifierOnWindows(c *cli.Context) bool {
	if runtime.GOOS != "windows" {
		return false
	}
	for _, arg := range c.Args() {
		if strings.HasPrefix(arg, "'") {
			return true
		}
	}
	return false
}

func initLogging(c *cli.Context) {
	if ! c.Bool("debug") {
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

func prepareSourceFileSystem(c *cli.Context)(afero.Fs, *sftpd.SftpFsContext, error) {
	if c.Bool("client") {
		username := c.String("username")
		host := fmt.Sprintf("%s:%d", c.String("host"), c.Int("port"))

		password, err := promptPasswordIfEmpty(c, fmt.Sprintf("Enter password for %s@%s:", username, host))
		if err != nil {
			return nil, nil, err
		}
		ctx, err := sftpd.NewSftpFsContext(username, password, host)
		if err != nil {
			return nil, nil, err
		}
		return sftpfs.New(ctx.Sftpc),ctx, nil
	}
	return afero.NewOsFs(), nil, nil
}

func promptPasswordIfEmpty(c *cli.Context, message string) (string, error) {
	if c.String("password") != "" {
		return c.String("password"), nil
	}
	println(message)
	pass, err := gopass.GetPasswd()
	return string(pass), err
}