package main

import "fmt"
import (
	"flag"
	"os"
	"strings"
	"regexp"
	"path/filepath"
	"bytes"
	"io"
	"github.com/urfave/cli"
)

type TransferSource struct {
	path string
}

const ERR_MISSING_PARAMS = 1
const ERR_COULD_NOT_COMPILE_SOURCE_PATTERN = 2

//var wordPtr = flag.String("word", "foo", "a string")
//var numbPtr = flag.Int("numb", 42, "an int")
//var svar string
var debug = flag.Bool("debug", false, "enable debug messages")
var help = flag.Bool("help", false, "show help")
var useRegex = flag.Bool("use-regex", false, "use real regex instead of glob patterns")
var simulate = flag.Bool("simulate", false, "simulation - just show preview, do not really transfer")
var times = flag.Bool("times", false, "keep times")
var move = flag.Bool("move", false, "move files instead of copying")

func dbg(a ...interface{}) {
	if (*debug) {
		fmt.Println(a...)
	}
}

func printlnWrapper(a ...interface{}) {
	fmt.Println(a...)
}

func exitWithError(message string, code int) {
	fmt.Println(message)
	os.Exit(code)
}

func exitWithHelp(message string) {
	fmt.Println(message)
	flag.Usage();
	os.Exit(ERR_MISSING_PARAMS)
}

func parseSourcePattern(sourcePattern string) (string, string) {
	path := sourcePattern;
	pattern := ""
	pathExists := false
	for {
		if _, err := os.Stat(path); err == nil {
			pattern = strings.Replace(sourcePattern, path, "", 1)
			if len(pattern) > 0 {
				pattern = pattern[1:]
			}
			pathExists = true
			break
		}

		lastSlashIndex := strings.LastIndex(normalizeDirSep(path), "/")
		if lastSlashIndex == -1 {
			break
		}
		path = path[0:lastSlashIndex]
	}

	if ! pathExists {
		path = ""
		pattern = sourcePattern
	}
	return path, pattern
}

func compilePattern(path string, pattern string) (*regexp.Regexp, error) {
	preparedPath := strings.Replace(path, "\\", "/", -1)
	preparedPattern := pattern //strings.Replace(pattern, "*", ".*", -1)
	preparedPatternToCompile := regexp.QuoteMeta(preparedPath) + "/" + preparedPattern
	dbg("pattern to compile:", preparedPatternToCompile)

	return regexp.Compile(preparedPatternToCompile)
}

func normalizeDirSep(path string) (string) {
	return strings.Replace(path, "\\", "/", -1)
}
func showFindResults(paths []string, sourcePattern *regexp.Regexp) {
	for i := 0; i < len(paths); i++ {
		printlnWrapper(paths[i])
		normalizedPath := normalizeDirSep(paths[i])
		sourcePattern.ReplaceAllStringFunc(normalizedPath, func(m string) string {
			parts := sourcePattern.FindStringSubmatch(m)
			i := 1
			for range parts[1:] {
				printlnWrapper("  $1: " + parts[i])
				i++

			}
			return m
		})

	}
}

// https://blog.gopheracademy.com/advent-2014/parsers-lexers/
// https://gist.github.com/yangls06/5464683
func GlobToRegex(glob string) (string) {
	var buffer bytes.Buffer
	r := strings.NewReader(glob)

	escape := false
	braceOpen := 0
	for {
		r, _, err := r.ReadRune()
		if (err != nil) {
			break;
		}

		if (escape) {
			buffer.WriteRune(r)
			escape = false
			continue
		}

		if (r == '\\') {
			buffer.WriteRune(r)
			escape = true;
			continue
		}

		if (r == '*') {
			buffer.WriteString(".*")
			continue
		}

		if (r == '.') {
			buffer.WriteString("\\.")
			continue
		}

		if (r == '{') {
			buffer.WriteString("(")
			braceOpen++
			continue
		}

		if (r == '}') {
			buffer.WriteString(")")
			braceOpen--
			continue
		}

		if (r == ',' && braceOpen > 0) {
			buffer.WriteString("|")
			continue
		}

		buffer.WriteRune(r)
	}

	return buffer.String()
}



//var debug = flag.Bool("debug", false, "enable debug messages")
//var help = flag.Bool("help", false, "show help")
//var useRegex = flag.Bool("use-regex", false, "use real regex instead of glob patterns")
//var simulate = flag.Bool("simulate", false, "simulation - just show preview, do not really transfer")
//var times = flag.Bool("times", false, "keep times")
//var move = flag.Bool("move", false, "move files instead of copying")
//


func main() {

	// default values for args
	app := cli.NewApp()
	app.Name = "graft"
	app.Usage = "find and copy files via command line"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "dry-run",
			Usage: "perform a dry-run without transferring files",
		},
	}
	//app.Flags = []cli.Flag{
	//	cli.StringFlag{
	//		Name:  "p, project-name",
	//		Usage: "Specify an alternate project name (default: directory name)",
	//	},
	//}



	//app.PositionalArgs = []cli.StringArg{
	//	Name: "branch",
	//	Optional: true,
	//	Value: &x
	//}
	// global level flags
	//app.Flags = []gangstaCli.Flag{
	//	gangstaCli.BoolFlag{
	//		Name:  "verbose",
	//		Usage: "Show more output",
	//	},
	//	gangstaCli.StringFlag{
	//		Name:  "f, file",
	//		Usage: "Specify an alternate fig file (default: fig.yml)",
	//	},
	//	gangstaCli.StringFlag{
	//		Name:  "p, project-name",
	//		Usage: "Specify an alternate project name (default: directory name)",
	//	},
	//}
	//
	//// Commands
	//app.Commands = []gangstaCli.Command{
	//	{
	//		Name: "build",
	//		Flags: []gangstaCli.Flag{
	//			gangstaCli.BoolFlag{
	//				Name:  "no-cache",
	//				Usage: "Do not use cache when building the image.",
	//			},
	//		},
	//		Usage:  "Build or rebuild services",
	//		Action: CmdBuild,
	//	},
	//	// etc...
	//	{
	//		Name: "run",
	//		Flags: []gangstaCli.Flag{
	//			gangstaCli.BoolFlag{
	//				Name:  "d",
	//				Usage: "Detached mode: Run container in the background, print new container name.",
	//			},
	//			gangstaCli.BoolFlag{
	//				Name:  "T",
	//				Usage: "Disables psuedo-tty allocation. By default `fig run` allocates a TTY.",
	//			},
	//			gangstaCli.BoolFlag{
	//				Name:  "rm",
	//				Usage: "Remove container after run.  Ignored in detached mode.",
	//			},
	//			gangstaCli.BoolFlag{
	//				Name:  "no-deps",
	//				Usage: "Don't start linked services.",
	//			},
	//		},
	//		Usage:  "Run a one-off command",
	//		Action: CmdRm,
	//	},
	//
	//	{
	//		Name: "up",
	//		Flags: []gangstaCli.Flag{
	//			gangstaCli.BoolFlag{
	//				Name:  "watch",
	//				Usage: "Watch build directory for changes and auto-rebuild/restart",
	//			},
	//			gangstaCli.BoolFlag{
	//				Name:  "d",
	//				Usage: "Detached mode: Run containers in the background, print new container names.",
	//			},
	//			gangstaCli.BoolFlag{
	//				Name:  "k,kill",
	//				Usage: "Kill instead of stop on terminal stignal",
	//			},
	//			gangstaCli.BoolFlag{
	//				Name:  "no-clean",
	//				Usage: "Don't remove containers after termination signal interrupt (CTRL+C)",
	//			},
	//			gangstaCli.BoolFlag{
	//				Name:  "no-deps",
	//				Usage: "Don't start linked services.",
	//			},
	//			gangstaCli.BoolFlag{
	//				Name:  "no-recreate",
	//				Usage: "If containers already exist, don't recreate them.",
	//			},
	//		},
	//		Usage:  "Create and start containers",
	//		Action: CmdUp,
	//	},
	//}


	app.Action = func(c *cli.Context) error {


		appCommand(c)


		//fmt.Println("src: ", sourcePattern)
		//fmt.Println("dst: ", destinationPattern)

		//if language == "spanish" {
		//	fmt.Println("Hola", name)
		//} else {
		//	fmt.Println("Hello", name)
		//}
		return nil
	}

	app.Run(os.Args)
	os.Exit(0)


}

func appCommand(c *cli.Context) {
	//patt, err := compilePattern("fixtures", "(?i)(.*)")
	//x := patt.ReplaceAllString("fixtures/global/textfile.txt", "test/$1")
	//dbg(x)
	//os.Exit(0)
	sourcePattern := ""
	if c.NArg() < 1 {
		fmt.Println("missing required parameter source-pattern, use --help parameter for usage instructions")
		return nil
	}

	sourcePattern = c.Args().Get(0)
	destinationPattern := ""
	if c.NArg() > 1{
		destinationPattern = c.Args().Get(1)
	}





	path, pattern := parseSourcePattern(sourcePattern)

	dbg("src - parameter:", sourcePattern)
	dbg("src - parsedPath: ", path)
	dbg("src - pattern: ", pattern)
	dbg("dst - parameter:", destinationPattern)

	dbg("regex preparation - before: " + pattern)
	var replacedPattern string
	if (*useRegex) {
		replacedPattern = pattern
	} else {
		replacedPattern = GlobToRegex(pattern)
	}
	dbg("regex preparation - after: " + replacedPattern)

	compiledPattern, err := compilePattern(path, "(?i)" + replacedPattern)
	if compiledPattern.NumSubexp() < 1 {
		compiledPattern, err = compilePattern(path, "(?i)(" + replacedPattern + ")")
	}

	if (err != nil) {
		exitWithError("could not compile source pattern: " + err.Error(), ERR_COULD_NOT_COMPILE_SOURCE_PATTERN)
	}

	dbg("=============================================");
	if destinationPattern == "" {
		dbg("search in path " + path + ", pattern: " + pattern)
	} else {
		dbg("replace in path " + path + ", pattern: " + pattern + ", replacement: " + destinationPattern)
	}
	dbg("=============================================");


	list := make([]string, 0)
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		dbg("===================================")
		dbg("path: " + path)

		normalizedPath := normalizeDirSep(path)
		// normalizedPath := strings.Replace(path, "\\", "/", -1)

		dbg("normalized: " + normalizedPath)

		if ! compiledPattern.MatchString(normalizedPath) {
			dbg("match: no")
			return nil
		}
		dbg("match: yes, appending to list")

		//compiledPattern.ReplaceAllStringFunc(normalizedPath, func(m string) string {
		//	parts := compiledPattern.FindStringSubmatch(m)
		//	i:=1
		//	for range parts[1:] {
		//		dbg("    match: " + parts[i])
		//		i++
		//
		//	}
		//	return m
		//})
		list = append(list, path)

		//if info.IsDir() {
		//	return nil
		//}
		//if filepath.Ext(path) == ".sh" {
		//	list = append(list, path)
		//}
		return nil
	})

	if destinationPattern == "" {
		showFindResults(list, compiledPattern)
		os.Exit(0)
	}

	printlnWrapper("===================================")
	if *move {
		printlnWrapper("move files: " + sourcePattern + " => " + destinationPattern)
	} else {
		printlnWrapper("copy files: " + sourcePattern + " => " + destinationPattern)
	}

	printlnWrapper("===================================")
	transferFiles(list, compiledPattern, destinationPattern)

	//if err != nil {
	//	exitWithError("walk failed: " + err.Error(), ERR_WALK_FAILED)
	//}

	//if compiledPattern.MatchString(foundFile) {
	//	fmt.Println("match: " + foundFile + " => " + preparedPatternToCompile)
	//} else {
	//	fmt.Println("no match: " + foundFile + " => " + preparedPatternToCompile)
	//}

	//res := compiledPattern.FindStringSubmatch(foundFile)
	//fmt.Printf("%v", res)

	// input := `bla bla b:foo="hop" blablabla b:bar="hu?"`
	//r := regexp.MustCompile(`(\bb:\w+=")([^"]+)`)
	//fmt.Println(compiledPattern.ReplaceAllStringFunc(foundFile, func(m string) string {
	//	parts := compiledPattern.FindStringSubmatch(m)
	//	// return parts[1] + complexFunc(parts[2])
	//	// fmt.Println(m)
	//	fmt.Println("0: " + parts[0])
	//	fmt.Println("1: " + parts[1])
	//	//fmt.Println("2: " + parts[2])
	//	return m
	//}))


	// var sourcePattern = [0]
	// var sourcePattern = flag.Args()[0]

	//err := filepath.Walk(sourcePattern, fileWalkCallback)
	//
	//for e := mapping.Front(); e != nil; e = e.Next() {
	//	fmt.Println(e.Value)
	//}


	//fmt.Printf("filepath.Walk() returned %v\n", err)
	//
	//fmt.Println("sourcePattern:", sourcePattern)
	//fmt.Println("destinationPattern:", destinationPattern)
	//if len(flag.Args()) == 0 {
	//	flag.Usage()
	//	os.Exit(1)
	//}
	//fmt.Println(*sourcePattern);

	// reader := io.ReaderAt(sourcePattern)
	// reader.ReadAt(),
}
func transferFiles(paths []string, sourcePattern *regexp.Regexp, replacement string) {
	for i := 0; i < len(paths); i++ {
		dbg("path: " + paths[i])
		dbg("patt: ", sourcePattern)
		dbg("repl: " + replacement)
		transferFile(paths[i], sourcePattern.ReplaceAllString(paths[i], replacement))
	}
}
func transferFile(src string, dst string) {
	printlnWrapper(src + " => " + dst)
	if *simulate {
		return
	}

	var inDirStats os.FileInfo
	inStats, err := os.Stat(src)
	inDirStats = inStats
	var srcSize int64 = 0

	var srcDir string
	if !inStats.IsDir() {
		srcDir = filepath.Dir(src)
		inDirStats, err = os.Stat(srcDir)
		srcSize = inStats.Size()
	}

	dbg("srcSize: ", srcSize)

	if err != nil {
		printlnWrapper("could not determine attributes for " + src + ": " + err.Error())
		return
	}

	var dstStats os.FileInfo
	dstStats, err = os.Stat(dst)
	dstExists := false
	var dstSize int64 = 0
	if !os.IsNotExist(err) {
		dstExists = true
		dstSize = dstStats.Size()
	}

	dbg("dstSize: ", dstSize)
	dbg("dstExists: ", dstExists)

	if inStats.IsDir() {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dst, inDirStats.Mode())
		}

		if err != nil {
			printlnWrapper("could not create destination directory " + dst + ": " + err.Error())
		}
		return
	}

	if *move && !dstExists {
		renameErr := os.Rename(src, dst)
		if renameErr != nil {
			printlnWrapper("could not rename " + src + ": " + renameErr.Error())
		} else {
			os.Remove(srcDir)
			return
		}
	}

	fi, inError := os.Open(src)
	defer fi.Close()
	if inError != nil {
		printlnWrapper("could not open source file " + src + ": " + err.Error())
		return
	}


	//flags := os.O_RDWR | os.O_CREATE
	//if dstExists {
	//	flags = os.O_RDWR | os.O_APPEND
	//}
	fo, outError := os.OpenFile(dst, os.O_RDWR | os.O_CREATE | os.O_APPEND, inStats.Mode())
	//var fo os.File
	//var outError error
	//if dstExists {
	//	fo, outError := os.Open(dst)
	//} else {
	//	fo, outError := os.Create(dst)
	//}

	defer fo.Close()
	if outError != nil {
		printlnWrapper("could not open destination file " + dst + ": " + err.Error())
		return
	}

	if srcSize == 0 {
		return
	}

	if dstExists {

		if (!areFilesEqual(fi, fo, srcSize, dstSize)) {
			printlnWrapper("source and destination are not equal " + src + " != " + dst)
			return
		}

		if *move {
			removeErr := os.Remove(dst)
			if removeErr != nil {
				printlnWrapper("Could not remove existing file before moving")
				return
			}

			renameErr := os.Rename(src, dst)
			if renameErr != nil {
				printlnWrapper("could not rename " + src + ": " + renameErr.Error())
			} else {
				os.Remove(srcDir)
				return
			}
		}

		_, fiErr := fi.Seek(dstSize, 0)
		if fiErr != nil {
			printlnWrapper("could not seek source file " + src + ": " + fiErr.Error())
			return
		}

		_, foErr := fo.Seek(dstSize, 0)
		if foErr != nil {
			printlnWrapper("could not seek destination file " + dst + ": " + foErr.Error())
			return
		}
	}

	buf := make([]byte, 1024)
	for {
		// read a chunk
		n, err := fi.Read(buf)
		if err != nil && err != io.EOF {
			printlnWrapper("reading file chunk failed: " + err.Error())
			return
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := fo.Write(buf[:n]); err != nil {
			printlnWrapper("writing file chunk failed: " + err.Error())
		}
	}

	if *times {
		os.Chtimes(dst, inStats.ModTime(), inStats.ModTime())
	}

	if (*move) {
		os.Remove(src)
		os.Remove(srcDir)
	}

	//var fo os.File
	//
	//if os.IsNotExist(err) {
	//	fo, err = os.CreateFile(dst)
	//} else {
	//	fo, err = os.OpenFile()
	//}


	//fi, inErr := os.Open(src)
	//if inErr != nil {
	//	printlnWrapper("could not open source file " + src + ": " + inErr.Error())
	//	return
	//}
	//fo, outErr := createOrOpenFile(dst)
	//if outErr != nil {
	//	printlnWrapper("could not open destination file " + dst + ": " + outErr.Error())
	//	return
	//}


	// os.Chtimes()
	//os.Chown()
	//os.Chmod()
}

func areFilesEqual(fi *os.File, fo *os.File, inSize int64, outSize int64) (bool) {

	if (outSize > inSize) {
		return false
	}

	var bufSize int64
	bufSize = 1024 * 1024 * 1024
	backBufSize := bufSize
	if bufSize > outSize {
		bufSize = outSize
		backBufSize = 0
	} else if outSize < bufSize * 2 {
		backBufSize = outSize - bufSize
	}

	fiBuf := make([]byte, bufSize)
	_, err := fi.ReadAt(fiBuf, 0)

	if err != nil {
		printlnWrapper("comparing files failed reading in buffer: " + err.Error())
	}

	foBuf := make([]byte, bufSize)
	_, err = fo.ReadAt(foBuf, 0)

	if err != nil {
		printlnWrapper("comparing files failed reading in out buffer: " + err.Error())
	}

	if ! bytes.Equal(fiBuf, foBuf) {
		return false
	}

	if backBufSize > 0 {
		backOffset := outSize - backBufSize
		fiBuf = make([]byte, backBufSize)
		_, err = fi.ReadAt(fiBuf, backOffset)
		if err != nil {
			printlnWrapper("comparing files failed reading in back buffer: " + err.Error())
		}
		foBuf = make([]byte, backBufSize)
		_, err = fo.ReadAt(foBuf, backOffset)
		if err != nil {
			printlnWrapper("comparing files failed reading out back buffer: " + err.Error())
		}
		if ! bytes.Equal(fiBuf, foBuf) {
			return false
		}
	}



	// buf := make([]byte, 1024)
	return true
}
