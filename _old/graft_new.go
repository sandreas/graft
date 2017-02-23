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
	"path/filepath"
)

var (
	app = kingpin.New("graft", "A command-line tool to locate and transfer files")
	sourceParameter = app.Arg("source-pattern", "source pattern - used to locate files (e.g. src/*)").Required().String()
	destinationParameter = app.Arg("destination-pattern", "destination pattern for transfer (e.g. dst/$1)").Default("").String()

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

var (
	source string
	destination string

	sourcePath string
	sourcePattern string
	sourcePathStat os.FileInfo


	err error
	dirsToRemove = make([]string, 0)
	exportFile os.File
)


func main() {
	/*
	if source is a file, no walk is needed
		if destination is a path, filename is appended to path
		if destination is a file, destination is kept


	if source is a path, these possibilities exist

	Pattern based:
	/data/* /data2/		- if data2 is a file, transfer is not performed
	 data/* data2/		- if data2 is a file, transfer is not performed

	Without pattern:
	 data/ data2      	- source is a path, destination should be path to, if data2 is a file, transfer is not performed
	 data/ data2/	  	- source is a path, destination is a path, if data2 is a file, transfer is not performed


	 */


	kingpin.MustParse(app.Parse(os.Args[1:]))


	source = *sourceParameter
	sourcePath, sourcePattern = pattern.ParsePathPattern(source)
	sourcePathStat, err = os.Stat(sourcePath)

	destination = *destinationParameter


	if err != nil {
		println("could not read source " + sourcePath + ": " + err.Error())
		return
	}


	printAction()
	initExport()

	if isSourceAFile() {
		handleFileSource()
		return;
	}
	handleWalkableSource()
}

func printAction() {
	if destination == "" {
		searchIn := sourcePath
		if sourcePath == "" {
			searchIn = "./"
		}

		searchFor := ""
		if sourcePattern != "" {
			searchFor = sourcePattern
		}
		prntln("search in '" + searchIn + "': " + searchFor)

	} else if (*move) {
		prntln("move: " + source + " => " + destination)
	} else {
		prntln("copy: " + source + " => " + destination)
	}
	prntln("")
}


func prntln(a ...interface{}) (n int, err error) {
	if ! *quiet {
		return fmt.Println(a...)
	}
	return n, err
}


func initExport() {
	if *exportTo != "" {
		file, err := os.Create(*exportTo)
		if err != nil {
			prntln("could not create export file " + *exportTo + ": " + err.Error())
		} else {
			exportFile = *file
			defer exportFile.Close()
		}
	}
}

func isSourceAFile()(bool) {
	return sourcePathStat.Mode().IsRegular()
}

func handleFileSource() {
	if strings.HasSuffix(destination, "/") || strings.HasSuffix(destination, "\\") {
		destination += sourcePathStat.Name()
	}

	transferElementHandler(source, destination)
}

func handleWalkableSource() {
	compiledPattern, _ := regexp.Compile("");
	hasCompiledPattern := false
	if sourcePattern != "" {
		if ! *regex {
			sourcePattern = pattern.GlobToRegex(sourcePattern)
		}

		caseInsensitiveQualifier := "(?i)"
		if *caseSensitive {
			caseInsensitiveQualifier = ""
		}

		compiledPattern, err := pattern.CompileNormalizedPathPattern(sourcePath, caseInsensitiveQualifier + sourcePattern)
		if err == nil && compiledPattern.NumSubexp() == 0 && sourcePattern != "" {
			compiledPattern, err = pattern.CompileNormalizedPathPattern(sourcePath, caseInsensitiveQualifier + "(" + sourcePattern + ")")
		}

		if err != nil {
			prntln("could not compile source pattern, please use slashes to qualify paths (recognized path: " + sourcePath + ", pattern" + sourcePattern + ")")
			return
		}
		hasCompiledPattern = true
	}



	destinationPath, destinationPattern := pattern.ParsePathPattern(destination)

	var dst string;
	var err error;

	if *filesFrom != "" {
		if ! file.Exists(*filesFrom) {
			prntln("Could not load files from " + *filesFrom)
			return
		}
		matchingPaths, err := file.ReadAllLinesFunc(*filesFrom, file.SkipEmptyLines)
		if err != nil {
			prntln("could read lines from file " + *filesFrom + ": " + err.Error())
			return
		}
		for _, p := range matchingPaths {
			if destinationPattern == "" {
				dst = pattern.NormalizeDirSep(destinationPath + p[len(sourcePath)+1:])
			} else {
				dst = compiledPattern.ReplaceAllString(pattern.NormalizeDirSep(p), pattern.NormalizeDirSep(destination))
			}


			if destination == "" {
				findElementHandler(p, compiledPattern)
			} else {
				transferElementHandler(p, dst)
			}
		}
	} else {
		matchingFiles, _ := file.WalkPathFiltered(sourcePath, func(f file.File, err error)(bool) {
			normalizedPath := filepath.ToSlash(f.Path)
			if hasCompiledPattern && !compiledPattern.MatchString(normalizedPath) {
				return false
			}

			return minAgeFilter(f) && maxAgeFilter(f)
		}, progressHandlerWalkPathByPattern)


		//for _, element := range matchingPaths {
		//	if dstPatt == "" {
		//		dst = pattern.NormalizeDirSep(dstPath + element[len(patternPath)+1:])
		//	} else {
		//		dst = compiledPattern.ReplaceAllString(pattern.NormalizeDirSep(element), pattern.NormalizeDirSep(destinationPattern))
		//	}
		//	transferElementHandler(element, dst)
		//}

		for _, f := range matchingFiles {
			p := f.Path

			if destinationPattern == "" {
				dst = pattern.NormalizeDirSep(destinationPath + p[len(sourcePath)+1:])
			} else {
				// dst = compiledPattern.ReplaceAllString(pattern.NormalizeDirSep(p), pattern.NormalizeDirSep(destination))
				dst = compiledPattern.ReplaceAllString(pattern.NormalizeDirSep(f.Path), pattern.NormalizeDirSep(destination))

			}

			if destination == "" {
				findElementHandler(f.Path, compiledPattern)
			} else {
				transferElementHandler(f.Path, dst)
			}
			exportLineIfRequested(f.Path)
		}

		if err != nil {
			prntln("could not write all lines to file " + *exportTo + ": " + err.Error())
		}
	}

	if *move {
		for _, dirToRemove := range dirsToRemove {
			os.Remove(dirToRemove)
		}
	}
}


func exportLineIfRequested(line string) {
	if exportFile != (os.File{}) {
		exportFile.WriteString(line + "\n")
	}
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
		progress = fmt.Sprintf("scanned: %d", entriesWalked)
	} else {
		progress = fmt.Sprintf("scanned: %d,  matched: %d", entriesWalked, entriesMatched)
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

func prnt(a...interface{}) (n int, err error) {
	if ! *quiet {
		return fmt.Print(a...)
	}
	return n, err
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


func findElementHandler(element string, compiledPattern *regexp.Regexp) {
	prntln(element)
	if *hideMatches || compiledPattern.String() == "" {
		return
	}
	elementMatches := pattern.BuildMatchList(compiledPattern, element)
	for i := 0; i < len(elementMatches); i++ {
		prntln("    $" + strconv.Itoa(i + 1) + ": " + elementMatches[i])
	}

}

func transferElementHandler(src, dst string) {

	srcStat, srcErr := os.Stat(src)
	if srcErr != nil {
		prntln("could not read source: ", srcErr)
		return
	}
	dstStat, _ := os.Stat(dst)
	dstExists := file.Exists(dst)

	prntln(src + " => " + dst)

	if *dryRun {
		return
	}


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