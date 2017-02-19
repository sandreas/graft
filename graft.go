package main

import (
	"os"
	"gopkg.in/alecthomas/kingpin.v2"
	"fmt"
	"github.com/sandreas/graft/pattern"
	"github.com/sandreas/graft/file"
	"strconv"
	"regexp"
	"path"
	"strings"
	"math"
	"time"
)

var (
	app = kingpin.New("graft", "A command-line tool to locate and transfer files")
	sourcePatternParameter = app.Arg("source-pattern", "source pattern - used to locate files (e.g. src/*)").Required().String()
	destinationPatternParameter = app.Arg("destination-pattern", "destination pattern for transfer (e.g. dst/$1)").Default("").String()

	exportTo = app.Flag("export-to", "export source listing to file, one line per found item").Default("").String()
	filesFrom = app.Flag("files-from", "import source listing from file, one line per item").Default("").String()

	minAge = app.Flag("min-age", " minimum age (e.g. -2 days, -8 weeks, 2015-10-10, etc.)").Default("").String()
	maxAge = app.Flag("max-age", "maximum age (e.g. 2 days, 8 weeks, 2015-10-10, etc.)").Default("").String()


	caseSensitive = app.Flag("case-sensitive", "be case sensitive when matching files and folders").Bool()
	dryRun = app.Flag("dry-run", "dry-run / simulation mode").Bool()
	hideMatches = app.Flag("hide-matches", "hide matches in search mode ($1: ...)").Bool()
	move = app.Flag("move", "move / rename files - do not make a copy").Bool()
	quiet = app.Flag("quiet", "quiet mode - do not show any output").Bool()
	regex = app.Flag("regex", "use a real regex instead of glob patterns (e.g. src/.*\\.jpg)").Bool()
	times = app.Flag("times", "transfer source modify times to destination").Bool()
)

var dirsToRemove = make([]string, 0)

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))
	sourcePattern := *sourcePatternParameter
	destinationPattern := *destinationPatternParameter

	patternPath, pat := pattern.ParsePathPattern(sourcePattern)
	if destinationPattern == "" {
		searchIn := patternPath
		if patternPath == "" {
			searchIn = "./"
		}

		searchFor := ""
		if pat != "" {
			searchFor = pat
		}
		prntln("search in '" + searchIn + "': " + searchFor)

	} else if (*move) {
		prntln("move: " + sourcePattern + " => " + destinationPattern)
	} else {
		prntln("copy: " + sourcePattern + " => " + destinationPattern)
	}
	prntln("")

	if ! *regex {
		pat = pattern.GlobToRegex(pat)
	}

	caseInsensitiveQualifier := "(?i)"
	if *caseSensitive {
		caseInsensitiveQualifier = ""
	}

	compiledPattern, err := pattern.CompileNormalizedPathPattern(patternPath, caseInsensitiveQualifier + pat)
	if err == nil && compiledPattern.NumSubexp() == 0 && pat != "" {
		compiledPattern, err = pattern.CompileNormalizedPathPattern(patternPath, caseInsensitiveQualifier + "(" + pat + ")")
	}

	if err != nil {
		prntln("could not compile source pattern, please use slashes to qualify paths (recognized path: " + patternPath + ", pattern" + pat + ")")
		return
	}

	var matchingPaths []string

	if *filesFrom != "" {
		if ! file.Exists(*filesFrom) {
			prntln("Could not load files from " + *filesFrom)
			return
		}
		matchingPaths, err = file.ReadAllLinesFunc(*filesFrom, file.SkipEmptyLines)
	} else {
		 //matchingPaths, err = file.WalkPathByPattern(patternPath, compiledPattern, progressHandlerWalkPathByPattern)
		matchingFiles, _ := file.WalkPathFiltered(patternPath, func(f file.File, err error)(bool) {
			normalizedPath := pattern.NormalizeDirSep(f.Path)
			if ! compiledPattern.MatchString(normalizedPath) {
				return false
			}

			return minAgeFilter(f) && maxAgeFilter(f)
		}, progressHandlerWalkPathByPattern)

		for _, element := range matchingFiles {
			matchingPaths = append(matchingPaths, element.Path)
		}

		if *exportTo != "" {
			exportFile(*exportTo, matchingPaths)
		}
	}

	if err != nil {
		prntln("Could not load sources path " + patternPath + ":", err.Error())
		return
	}



	if destinationPattern == "" {
		for _, element := range matchingPaths {
			findElementHandler(element, compiledPattern)
		}
		return
	}

	dstPath, dstPatt := pattern.ParsePathPattern(destinationPattern)
	var dst string
	for _, element := range matchingPaths {
		if dstPatt == "" {
			dst = pattern.NormalizeDirSep(dstPath + element[len(patternPath)+1:])
		} else {
			dst = compiledPattern.ReplaceAllString(pattern.NormalizeDirSep(element), pattern.NormalizeDirSep(destinationPattern))
		}
		transferElementHandler(element, dst)
	}

	if *move {
		for _, dirToRemove := range dirsToRemove {
			os.Remove(dirToRemove)
		}
	}
	return
}

func minAgeFilter(f file.File)(bool) {
	if *minAge == "" {
		return true
	}

	minAgeTime, err := pattern.StrToAge(*minAge, time.Now())
	if err != nil {
		return false
	}
	return minAgeTime.UnixNano() > f.ModTime().UnixNano()
}

func maxAgeFilter(f file.File)(bool) {
	if *maxAge == "" {
		return true
	}

	maxAgeTime, err := pattern.StrToAge(*maxAge, time.Now())
	if err != nil {
		return false
	}
	return maxAgeTime.UnixNano() < f.ModTime().UnixNano()
}


func progressHandlerWalkPathByPattern(entriesWalked, entriesMatched int64, finished bool) (int64) {
	var progress string;
	if entriesMatched == 0 {
		progress = fmt.Sprintf("scanning - total: %d", entriesWalked)
	} else {
		progress = fmt.Sprintf("scanning - total: %d,  matches: %d", entriesWalked, entriesMatched)
	}
	// prnt("\x0c" + progressBar)
	prnt("\r" + progress)
	if finished {
		prntln("")
		prntln("")
	}
	if(entriesWalked > 1000) {
		return 500
	}
	return 100
}


func exportFile(file string, lines []string) {
	f, err := os.Create(*exportTo)
	if err != nil {
		prntln("could not create export file " + file + ": " + err.Error())
		return;
	}
	_, err = f.WriteString(strings.Join(lines, "\n"))
	defer f.Close()
	if err != nil {
		prntln("could not write export file " + file + ": " + err.Error())
	}

}

func appendRemoveDir(dir string) {
	if (*move) {
		dirsToRemove = append(dirsToRemove, dir)
	}
}

func handleProgress(bytesTransferred, size, chunkSize int64) (int64) {

	if size <= 0 {
		return chunkSize
	}

	percent := float64(bytesTransferred) / float64(size)
	charCountWhenFullyTransmitted := 20
	progressChars := int(math.Floor(percent * float64(charCountWhenFullyTransmitted)))
	normalizedInt :=percent * 100
	progressBar := fmt.Sprintf("[%-" + strconv.Itoa(charCountWhenFullyTransmitted + 1)+ "s] %3d%%", strings.Repeat("=", progressChars) + ">", normalizedInt)

	prnt("\r" + progressBar)
	if bytesTransferred == size {
		prntln("")
	}
	// fmt.Print("\r" + progressBar)
	return chunkSize
}

func prntln(a ...interface{}) (n int, err error) {
	if ! *quiet {
		return fmt.Println(a...)
	}
	return n, err
}

func prnt(a...interface{}) (n int, err error) {
	if ! *quiet {
		return fmt.Print(a...)
	}
	return n, err
}

func findElementHandler(element string, compiledPattern *regexp.Regexp) {
	prntln(element)
	if *hideMatches {
		return
	}
	elementMatches := pattern.BuildMatchList(compiledPattern, element)
	for i := 0; i < len(elementMatches); i++ {
		prntln("    $" + strconv.Itoa(i + 1) + ": " + elementMatches[i])
	}

}

func transferElementHandler(src, dst string) {

	prntln(src + " => " + dst)

	if *dryRun {
		return
	}

	srcStat, srcErr := os.Stat(src)

	if srcErr != nil {
		prntln("could not read source: ", srcErr)
		return
	}

	dstStat, _ := os.Stat(dst)
	dstExists := file.Exists(dst)
	if srcStat.IsDir() {
		if ! dstExists {
			if os.MkdirAll(dst, srcStat.Mode()) != nil {
				prntln("Could not create destination directory")
			}
			appendRemoveDir(dst)
			fixTimes(dst, srcStat)
			return
		}

		if dstStat.IsDir() {
			appendRemoveDir(dst)
			fixTimes(dst, srcStat)
			return
		}

		prntln("destination already exists as file, source is a directory")
		return
	}

	if dstExists && dstStat.IsDir() {
		prntln("destination already exists as directory, source is a file")
		return
	}

	srcDir := path.Dir(src)
	srcDirStat, _ := os.Stat(srcDir)

	dstDir := path.Dir(dst)
	if ! file.Exists(dstDir) {
		os.MkdirAll(dstDir, srcDirStat.Mode())
	}

	if *move {
		renameErr := os.Rename(src, dst)
		if renameErr == nil {
			appendRemoveDir(srcDir)
			fixTimes(dst, srcStat)
			return
		}
		prntln("Could not rename source")
		return
	}

	srcPointer, srcPointerErr := os.Open(src)
	if srcPointerErr != nil {
		prntln("Could not open source file")
		return
	}
	dstPointer, dstPointerErr := os.OpenFile(dst, os.O_WRONLY | os.O_CREATE, srcStat.Mode())

	if dstPointerErr != nil {
		prntln("Could not create destination file", dstPointerErr.Error())
		return
	}

	file.CopyResumed(srcPointer, dstPointer, handleProgress)
	fixTimes(dst, srcStat)
}

func fixTimes(dst string, inStats os.FileInfo) {
	if *times {
		os.Chtimes(dst, inStats.ModTime(), inStats.ModTime())
	}
}