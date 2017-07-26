package main

import (
	"fmt"
	"os"
	"github.com/alexflint/go-arg"
	"os/user"
	"log"
	"io"
	"github.com/sandreas/graft/pattern"
	"github.com/sandreas/graft/newpattern"
)

//var (
//	app = kingpin.New("graft", "A command-line tool to locate and transfer files")
//	sourcePatternParameter = app.Arg("source-pattern", "source pattern - used to locate files (e.g. src/*)").Required().String()
//	destinationPatternParameter = app.Arg("destination-pattern", "destination pattern for transfer (e.g. dst/$1)").Default("").String()
//
//	exportTo = app.Flag("export-to", "export source listing to file, one line per item").Default("").String()
//	filesFrom = app.Flag("files-from", "import source listing from file, one line per item").Default("").String()
//
//	minAge = app.Flag("min-age", " minimum age (e.g. 2d, 8w, 2016-12-24, etc. - see docs for valid time formats)").Default("").String()
//	maxAge = app.Flag("max-age", "maximum age (e.g. 2d, 8w, 2016-12-24, etc. - see docs for valid time formats)").Default("").String()
//
//	caseSensitive = app.Flag("case-sensitive", "be case sensitive when matching files and folders").Bool()
//	dryRun = app.Flag("dry-run", "dry-run / simulation mode").Bool()
//	hideMatches = app.Flag("hide-matches", "hide matches in search mode ($1: ...)").Bool()
//	move = app.Flag("move", "move / rename files - do not make a copy").Bool()
//	quiet = app.Flag("quiet", "quiet mode - do not show any output").Bool()
//	regex = app.Flag("regex", "use a real regex instead of glob patterns (e.g. src/.*\\.jpg)").Bool()
//	times = app.Flag("times", "transfer source modify times to destination").Bool()
//	serve = app.Flag("serve", "start a server on this port").Default("0").String()
//
//	debug = app.Flag("debug", "enable debug logging").Bool()
//
//)

//var dirsToRemove = make([]string, 0)
//var minAgeTime time.Time
//var maxAgeTime time.Time

// cli tools:
// https://github.com/alexflint/go-arg
// https://github.com/jessevdk/go-flags
// https://github.com/spf13/pflag
// https://github.com/octago/sflags

// Input / Colors:
// https://github.com/dixonwille/wlog

//
//type DatabaseOptions struct {
//	Host     string
//	Username string
//	Password string
//}
//
//type LogOptions struct {
//	LogFile string
//	Verbose bool
//}
//
//func main() {
//	var args struct {
//		DatabaseOptions
//		LogOptions
//	}
//	arg.MustParse(&args)
//}

type PositionalArguments struct {
	Source      string `arg:"positional"`
	Destination string `arg:"positional"`
}

type BooleanFlags struct {
	Verbose bool `arg:"-v,help:be verbose"`
	Debug bool `arg:"-d,help:debug mode with logging to Stdout and into $HOME/.graft/application.log"`
}

var args struct {
	PositionalArguments
	BooleanFlags
}

func main() {
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "--help")
	}
	arg.MustParse(&args)

	initLogging()
	log.Printf("graft is starting...")

	sourcePattern := newpattern.NewSourcePattern(args.Source)
	sourceFiles := findFilesForPattern(sourcePattern)

	for key, value := range sourceFiles {
		println(key)
		println(value)
	}

	//fmt.Printf("Source: %v\n", args.Source)
	//fmt.Printf("Destination: %v\n", args.Destination)
	//fmt.Printf("Verbose: %v\n", args.Verbose)

}
func findFilesForPattern(sourcePattern *newpattern.SourcePattern) map[string]string {
	m := make(map[string]string)

	return m
}

func (PositionalArguments) Description() string {
	return "graft 0.2 - a command line application to search for and transfer files\n"
}

func initLogging() {
	if ! args.Debug {
		return
	}
	log.SetOutput(os.Stdout)

	homeDir, err := createHomeDirectoryIfNotExists()
	if err != nil {
		log.Println("could not create home directory: ", homeDir, err)
	}
	logFileName := homeDir + "/graft.log"
	os.Remove(logFileName)
	logFile, err := os.OpenFile(logFileName, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
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


//kingpin.MustParse(app.Parse(os.Args[1:]))
//
//sourcePattern := *sourcePatternParameter
//destinationPattern := *destinationPatternParameter
//
//if *debug {
//	initDebug()
//}
//
//
//serveOnPort, err := strconv.Atoi(*serve)
//if err != nil {
//	prntln("Invalid argument for serve: " + err.Error())
//	return
//}
//
//graftHomeDir := getGraftHomeDirectory()
//
//patternPath, pat := pattern.ParsePathPattern(sourcePattern)
//
//
//isFile, sourcePathStat, err := file.IsFile(patternPath)
//
//if err != nil {
//	prntln("Could not check pattern path: " + err.Error())
//	return
//}
//
//if isFile {
//	if strings.HasSuffix(destinationPattern, "/") || strings.HasSuffix(destinationPattern, "\\") {
//		destinationPattern += sourcePathStat.Name()
//	}
//	transferElementHandler(sourcePattern, destinationPattern)
//	return
//}
//
//if serveOnPort != 0 {
//	destinationPattern = ""
//}
//
//if destinationPattern == "" {
//	searchIn := patternPath
//	searchFor := ""
//	if pat != "" {
//		searchFor = pat
//	}
//	prntln("search in '" + searchIn + "': " + searchFor)
//
//} else if (*move) {
//	prntln("move: " + sourcePattern + " => " + destinationPattern)
//} else {
//	prntln("copy: " + sourcePattern + " => " + destinationPattern)
//}
//prntln("")
//
//if ! *regex {
//	pat = pattern.GlobToRegex(pat)
//}
//
//caseInsensitiveQualifier := "(?i)"
//if *caseSensitive {
//	caseInsensitiveQualifier = ""
//}
//
//// append $ for end of string
//if (!strings.HasSuffix(pat, "$")) || strings.HasSuffix(pat, "\\$") {
//	pat += "$"
//}
//
//compiledPattern, err := pattern.CompileNormalizedPathPattern(patternPath, caseInsensitiveQualifier + pat)
//if err == nil && compiledPattern.NumSubexp() == 0 && pat != "" {
//	compiledPattern, err = pattern.CompileNormalizedPathPattern(patternPath, caseInsensitiveQualifier + "(" + pat + ")")
//}
//
////prntln(compiledPattern)
//
//if err != nil {
//	prntln("could not compile source pattern, please use slashes to qualify paths (recognized path: " + patternPath + ", pattern" + pat + ")")
//	return
//}
//
//var matchingPaths []string
//
//if *filesFrom != "" {
//	if ! file.Exists(*filesFrom) {
//		prntln("Could not load files from " + *filesFrom)
//		return
//	}
//	matchingPaths, err = file.ReadAllLinesFunc(*filesFrom, file.SkipEmptyLines)
//} else {
//	if *minAge != "" {
//		minAgeTime, err = pattern.StrToAge(*minAge, time.Now())
//		if err != nil {
//			prntln("Could not parse --min-age: " + err.Error())
//			return
//		}
//	}
//	if *maxAge != "" {
//		maxAgeTime, err = pattern.StrToAge(*maxAge, time.Now())
//		if err != nil {
//			prntln("Could not parse --max-age: " + err.Error())
//			return
//		}
//	}
//
//	//matchingPaths, err = file.WalkPathByPattern(patternPath, compiledPattern, progressHandlerWalkPathByPattern)
//	matchingFiles, _ := file.WalkPathFiltered(patternPath, func(f file.File, err error) (bool) {
//		normalizedPath := pattern.NormalizeDirSep(f.Path)
//		if ! compiledPattern.MatchString(normalizedPath) {
//			return false
//		}
//
//		return minAgeFilter(f) && maxAgeFilter(f)
//	}, progressHandlerWalkPathByPattern)
//
//	for _, element := range matchingFiles {
//		matchingPaths = append(matchingPaths, element.Path)
//	}
//
//	if *exportTo != "" {
//		exportFile(*exportTo, matchingPaths)
//	}
//}
//
//if err != nil {
//	prntln("Could not load sources path " + patternPath + ":", err.Error())
//	return
//}
//
//
//
//if destinationPattern == "" {
//	for _, element := range matchingPaths {
//		findElementHandler(element, compiledPattern)
//	}
//
//	if (serveOnPort != 0) {
//		createGraftHomePathIfNotExists()
//		sftpd.NewSimpleServer(graftHomeDir, "0.0.0.0", 2022, "graft", "graft", matchingPaths)
//	}
//	return
//}
//
//dstPath, dstPatt := pattern.ParsePathPattern(destinationPattern)
//
//// replace $1_ with ${1}_ to prevent problems
//dollarUnderscore, _ := regexp.Compile("\\$([1-9][0-9]*)_")
//destinationPattern = dollarUnderscore.ReplaceAllString(destinationPattern, "${$1}_")
//
//var dst string
//for _, element := range matchingPaths {
//	if dstPatt == "" {
//		dst = pattern.NormalizeDirSep(dstPath + element[len(patternPath) + 1:])
//	} else {
//		dst = compiledPattern.ReplaceAllString(pattern.NormalizeDirSep(element), pattern.NormalizeDirSep(destinationPattern))
//	}
//	transferElementHandler(element, dst)
//}
//
//if *move {
//	for _, dirToRemove := range dirsToRemove {
//		os.Remove(dirToRemove)
//	}
//}
//return
//}

//
//func initDebug() {
//	createGraftHomePathIfNotExists()
//
//	logFile, err := os.OpenFile(getGraftHomeDirectory() + "/graft.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
//	if err != nil {
//
//	}
//	defer logFile.Close()
//
//	mw := io.MultiWriter(os.Stdout, logFile)
//	log.SetOutput(mw)
//}
//
//func getGraftHomeDirectory() string {
//	usr, err := user.Current()
//
//	if err != nil {
//		println("Could not determine current user ", err)
//		os.Exit(1)
//	}
//	return usr.HomeDir + "/.graft"
//}
//
//func createGraftHomePathIfNotExists() string {
//	graftHomePath := getGraftHomeDirectory()
//	mode := int(0755)
//	if _, err := os.Stat(graftHomePath); err != nil {
//		err := os.Mkdir(graftHomePath, os.FileMode(mode))
//		if err != nil {
//			println("Could not create home directory " + graftHomePath)
//			os.Exit(1)
//		}
//	}
//	return graftHomePath
//
//}
//
//func minAgeFilter(f file.File) (bool) {
//	if *minAge == "" {
//		return true
//	}
//	return minAgeTime.UnixNano() > f.ModTime().UnixNano()
//}
//
//func maxAgeFilter(f file.File) (bool) {
//	if *maxAge == "" {
//		return true
//	}
//	return maxAgeTime.UnixNano() < f.ModTime().UnixNano()
//}
//
//func progressHandlerWalkPathByPattern(entriesWalked, entriesMatched int64, finished bool) (int64) {
//	var progress string;
//	if entriesMatched == 0 {
//		progress = fmt.Sprintf("scanning - total: %d", entriesWalked)
//	} else {
//		progress = fmt.Sprintf("scanning - total: %d,  matches: %d", entriesWalked, entriesMatched)
//	}
//	// prnt("\x0c" + progressBar)
//	prnt("\r" + progress)
//	if finished {
//		prntln("")
//		prntln("")
//	}
//	if (entriesWalked > 1000) {
//		return 500
//	}
//	return 100
//}
//
//func exportFile(file string, lines []string) {
//	f, err := os.Create(*exportTo)
//	if err != nil {
//		prntln("could not create export file " + file + ": " + err.Error())
//		return;
//	}
//	_, err = f.WriteString(strings.Join(lines, "\n"))
//	defer f.Close()
//	if err != nil {
//		prntln("could not write export file " + file + ": " + err.Error())
//	}
//
//}
//
//func appendRemoveDir(dir string) {
//	if (*move) {
//		dirsToRemove = append(dirsToRemove, dir)
//	}
//}
//
//var timerLastUpdate time.Time
//var reportInterval int64
//var bytesLastUpdate int64
//
//func startTimer(bytesTransferred, interval int64) {
//	if bytesTransferred == 0 {
//		timerLastUpdate = time.Now()
//		reportInterval = interval
//		bytesLastUpdate = 0
//	}
//}
//
//func getReportStatus(bytesTransferred, size int64) (bool, float64, float64) {
//	timeDiffNano := time.Now().UnixNano() - timerLastUpdate.UnixNano()
//	timeDiffSeconds := float64(timeDiffNano) / float64(time.Second)
//	// if timeDiffSeconds >= float64(reportInterval) {
//	if timeDiffNano >= reportInterval {
//		bytesDiff := bytesTransferred - bytesLastUpdate
//		bytesPerSecond := float64(float64(bytesDiff) / float64(timeDiffSeconds))
//		percent := float64(bytesTransferred) / float64(size)
//
//		bytesLastUpdate = bytesTransferred
//		timerLastUpdate = time.Now()
//
//		return true, bytesPerSecond, percent
//	}
//
//	return false, 0, 0
//}
//
//func handleProgress(bytesTransferred, size, chunkSize int64) (int64) {
//
//	if size <= 0 {
//		return chunkSize
//	}
//
//	startTimer(bytesTransferred, 1 * int64(time.Second))
//	shouldReport, bytesPerSecond, percent := getReportStatus(bytesTransferred, size)
//	if shouldReport {
//		bandwidthOutput := "   " + bytefmt.FormatBytes(bytesPerSecond, 2, true) + "/s"
//		charCountWhenFullyTransmitted := 20
//		progressChars := int(math.Floor(percent * float64(charCountWhenFullyTransmitted)))
//		normalizedInt := percent * 100
//		percentOutput := strconv.FormatFloat(normalizedInt, 'f', 2, 64)
//		if bytesPerSecond == 0 {
//			bandwidthOutput = ""
//		}
//		progressBar := fmt.Sprintf("[%-" + strconv.Itoa(charCountWhenFullyTransmitted + 1) + "s] " + percentOutput + "%%" + bandwidthOutput, strings.Repeat("=", progressChars) + ">")
//
//		prnt("\r" + progressBar)
//	}
//
//	if bytesTransferred == size {
//		prntln("")
//	}
//
//	return chunkSize
//}
//
//func prntln(a ...interface{}) (n int, err error) {
//	if ! *quiet {
//		return fmt.Println(a...)
//	}
//	return n, err
//}
//
//func prnt(a...interface{}) (n int, err error) {
//	if ! *quiet {
//		return fmt.Print(a...)
//	}
//	return n, err
//}
//
//func findElementHandler(element string, compiledPattern *regexp.Regexp) {
//	prntln(element)
//	if *hideMatches || *serve != "0" {
//		return
//	}
//	elementMatches := pattern.BuildMatchList(compiledPattern, element)
//	for i := 0; i < len(elementMatches); i++ {
//		prntln("    $" + strconv.Itoa(i + 1) + ": " + elementMatches[i])
//	}
//
//}
//
//func transferElementHandler(src, dst string) {
//
//	prntln(src + " => " + dst)
//
//	if *dryRun {
//		return
//	}
//
//	srcStat, srcErr := os.Stat(src)
//
//	if srcErr != nil {
//		prntln("could not read source: ", srcErr)
//		return
//	}
//
//	dstStat, _ := os.Stat(dst)
//	dstExists := file.Exists(dst)
//	if srcStat.IsDir() {
//		if ! dstExists {
//			if os.MkdirAll(dst, srcStat.Mode()) != nil {
//				prntln("Could not create destination directory")
//			}
//			appendRemoveDir(dst)
//			fixTimes(dst, srcStat)
//			return
//		}
//
//		if dstStat.IsDir() {
//			appendRemoveDir(dst)
//			fixTimes(dst, srcStat)
//			return
//		}
//
//		prntln("destination already exists as file, source is a directory")
//		return
//	}
//
//	if dstExists && dstStat.IsDir() {
//		prntln("destination already exists as directory, source is a file")
//		return
//	}
//
//	srcDir := path.Dir(src)
//	srcDirStat, _ := os.Stat(srcDir)
//
//	dstDir := path.Dir(dst)
//	if ! file.Exists(dstDir) {
//		os.MkdirAll(dstDir, srcDirStat.Mode())
//	}
//
//	if *move {
//		renameErr := os.Rename(src, dst)
//		if renameErr == nil {
//			appendRemoveDir(srcDir)
//			fixTimes(dst, srcStat)
//			return
//		}
//		prntln("Could not rename source")
//		return
//	}
//
//	srcPointer, srcPointerErr := os.Open(src)
//	if srcPointerErr != nil {
//		prntln("Could not open source file")
//		return
//	}
//	dstPointer, dstPointerErr := os.OpenFile(dst, os.O_WRONLY | os.O_CREATE, srcStat.Mode())
//
//	if dstPointerErr != nil {
//		prntln("Could not create destination file", dstPointerErr.Error())
//		return
//	}
//
//	file.CopyResumed(srcPointer, dstPointer, handleProgress)
//	fixTimes(dst, srcStat)
//}
//
//func fixTimes(dst string, inStats os.FileInfo) {
//	if *times {
//		os.Chtimes(dst, inStats.ModTime(), inStats.ModTime())
//	}
//}
