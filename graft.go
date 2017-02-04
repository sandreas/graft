package main

import "fmt"
import (
	"flag"
	"os"
	"strings"
	"regexp"
	"path/filepath"
)

type TransferSource struct {
	path string
}

const ERR_MISSING_PARAMS = 1
const ERR_COULD_NOT_COMPILE_SOURCE_PATTERN = 2

//var wordPtr = flag.String("word", "foo", "a string")
//var numbPtr = flag.Int("numb", 42, "an int")
//var svar string
var debug = flag.Bool("debug", true, "enable debug messages")
var help = flag.Bool("help", false, "show help")

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
	// todo: check if pattern contains groups => (*group1)(*group2), if not, treat whole pattern as group
	// preparedPattern = "(" + preparedPattern + ")"

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

func main() {
	flag.Parse()

	//fmt.Println("word:", *wordPtr)
	//fmt.Println("numb:", *numbPtr)
	//fmt.Println("fork:", *boolPtr)
	//fmt.Println("svar:", svar)
	//fmt.Println("tail:", flag.Args())


	flagArgs := flag.Args();


	sourcePattern := ""
	if *help || len(flagArgs) < 1 {
		exitWithHelp("Please specify at least valid source pattern")

	}
	sourcePattern = flagArgs[0]
	path, pattern := parseSourcePattern(sourcePattern)

	dbg("src - parameter:", sourcePattern)
	dbg("src - parsedPath: ", path)
	dbg("src - pattern: ", pattern)

	destinationPattern := ""
	if len(flagArgs) > 1 {
		destinationPattern = flagArgs[1]
	}
	dbg("dst - parameter:", destinationPattern)

	compiledPattern, err := compilePattern(path, pattern)
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
		printlnWrapper(paths[i] + " => " + sourcePattern.ReplaceAllString(paths[i], replacement))

		//normalizedPath := normalizeDirSep(paths[i])
		//sourcePattern.ReplaceAllStringFunc(normalizedPath, func(m string) string {
		//	parts := sourcePattern.FindStringSubmatch(m)
		//	i := 1
		//	for range parts[1:] {
		//		println("  $1: " + parts[i])
		//		i++
		//
		//	}
		//	return m
		//})

	}
}
