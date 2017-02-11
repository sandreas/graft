package main

import (
	"os"
	"github.com/urfave/cli"
	"fmt"
	"errors"
	"github.com/sandreas/graft/pattern"
	"github.com/sandreas/graft/file"
	"strconv"
	"regexp"
	"path"
	"strings"
	"math"
)

var settings *cli.Context
var dirsToRemove = make([]string, 0)

func main() {
	app := cli.NewApp()
	app.Name = "graft"
	app.Usage = "find and copy files via command line"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "move",
			Usage: "move files instead of copy",
		},
		cli.BoolFlag{
			Name:  "dry-run",
			Usage: "perform a dry run",
		},
		cli.BoolFlag{
			Name:  "regex",
			Usage: "use a real regular expression instead of glob patterns",
		},
		cli.BoolFlag{
			Name:  "case-sensitive",
			Usage: "match case sensitive",
		},
		cli.BoolFlag{
			Name:  "quiet",
			Usage: "do not show any output",
		},
	}

	app.Action = mainAction

	app.Run(os.Args)
}

func appendRemoveDir(dir string) {
	if (settings.Bool("move")) {
		dirsToRemove = append(dirsToRemove, dir)
	}
}

func handleProgress(bytesTransferred, size, chunkSize int64) (int64) {
	if size <= 0 {
		return chunkSize
	}

	percent := bytesTransferred / size

	progressChars := int(math.Floor(float64(percent * 10)) * 2)

	progressBar := fmt.Sprintf("[%-21s] %3d/100%%", strings.Repeat("=", progressChars) + ">", percent * 100)
	prnt("\x0c" + progressBar)
	// fmt.Print("\r" + progressBar)
	return chunkSize
}

func mainAction(c *cli.Context) error {
	settings = c
	sourcePattern := ""
	if c.NArg() < 1 {
		return errors.New("missing required parameter source-pattern, use --help parameter for usage instructions")
	}

	sourcePattern = c.Args().Get(0)
	destinationPattern := ""
	if c.NArg() > 1 {
		destinationPattern = c.Args().Get(1)
	}

	if destinationPattern == "" {
		prntln("search: " + sourcePattern)
	} else if(settings.Bool("move")){
		prntln("move: " + sourcePattern + " => " + destinationPattern)
	} else {
		prntln("copy: " + sourcePattern + " => " + destinationPattern)
	}
	prntln("")
	prntln("")

	patternPath, pat := pattern.ParsePathPattern(sourcePattern)

	if ! settings.Bool("regex") {
		pat = pattern.GlobToRegex(pat)
	}


	caseInsensitiveQualifier := "(?i)"
	if settings.Bool("case-sensitive") {
		caseInsensitiveQualifier = ""
	}

	compiledPattern, err := pattern.CompileNormalizedPathPattern(patternPath, caseInsensitiveQualifier + pat)
	if compiledPattern.NumSubexp() == 0 {
		compiledPattern, err = pattern.CompileNormalizedPathPattern(patternPath, caseInsensitiveQualifier + "(" + pat + ")")
	}

	if err != nil {
		return errors.New("could not compile source pattern " + patternPath + ", " + pat)
	}

	matchingPaths, err := file.WalkPathByPattern(patternPath, compiledPattern)

	if err != nil {
		return err
	}

	if destinationPattern == "" {
		for _, element := range matchingPaths {
			findElementHandler(element, compiledPattern)
		}
		return nil
	}

	for _, element := range matchingPaths {
		transferElementHandler(element, destinationPattern, compiledPattern)
	}

	if settings.Bool("move") {
		for _, dirToRemove := range dirsToRemove {
			os.Remove(dirToRemove)
		}
	}
	return nil
}

func prntln(a ...interface{}) (n int, err error) {
	if ! settings.Bool("quiet") {
		return fmt.Println(a...)
	}
	return n, err
}

func prnt(a...interface{}) (n int, err error) {
	if ! settings.Bool("quiet") {
		return fmt.Print(a...)
	}
	return n, err
}

func findElementHandler(element string, compiledPattern *regexp.Regexp) {
	elementMatches := pattern.BuildMatchList(compiledPattern, element)
	prntln(element)
	for i := 0; i < len(elementMatches); i++ {
		prntln("    $" + strconv.Itoa(i + 1) + ": " + elementMatches[i])
	}

}

func transferElementHandler(src, destinationPattern string, compiledPattern  *regexp.Regexp) {
	dst := compiledPattern.ReplaceAllString(src, destinationPattern)

	prntln(src + " => " + dst)

	if settings.Bool("dry-run") {
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
			return
		}

		if dstStat.IsDir() {
			appendRemoveDir(dst)
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
	_, dstDirErr := os.Stat(dstDir)
	if ! os.IsNotExist(dstDirErr) {
		os.MkdirAll(dstDir, srcDirStat.Mode())
	}

	if settings.Bool("move") {
		renameErr := os.Rename(src, dst)
		if renameErr == nil {
			appendRemoveDir(srcDir)
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

	dstPointer, dstPointerErr := os.OpenFile(dst, os.O_RDWR | os.O_CREATE, srcStat.Mode())
	if dstPointerErr != nil {
		prntln("Could not create destination file")
		return
	}

	file.CopyResumed(srcPointer, dstPointer, handleProgress)
}
