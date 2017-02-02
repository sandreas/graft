package main
import "fmt"
import (
	"flag"
//	"os"
//	"path/filepath"
	"os"
	//"container/list"
	//"io"
	"strings"
	"regexp"
	"path/filepath"
)

//var wordPtr = flag.String("word", "foo", "a string")
//var numbPtr = flag.Int("numb", 42, "an int")
//var boolPtr = flag.Bool("fork", false, "a bool")
//var svar string

// var mapping = list.New()



//type TransferSource struct {
//	path string
//	pattern string
//}

type TransferPair struct {
	from string
	to string
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
		normalizedDirectorySeparatorPath := strings.Replace(path, "\\", "/", -1)
		lastSlashIndex := strings.LastIndex(normalizedDirectorySeparatorPath, "/")
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


//func fileWalkCallback(path string, f os.FileInfo, err error) error {
//	fmt.Printf("Visited: %s\n", path)
//	mapping.PushBack(TransferPair{path, path})
//	return nil
//}


func main() {
	//flag.StringVar(&svar, "svar", "bar", "a string var")

	/*
if _, err := os.Stat("/path/to/whatever"); os.IsNotExist(err) {
  // path/to/whatever does not exist
}
In the above example we are not checking if err != nil because os.IsNotExist(nil) == false.

To check if a file exists, equivalent to Python's if os.path.exists(filename):

if _, err := os.Stat("/path/to/whatever"); err == nil {
  // path/to/whatever exists
}

	 */

	flag.Parse()

	//fmt.Println("word:", *wordPtr)
	//fmt.Println("numb:", *numbPtr)
	//fmt.Println("fork:", *boolPtr)
	//fmt.Println("svar:", svar)
	//fmt.Println("tail:", flag.Args())


	flagArgs := flag.Args();
	sourcePattern := ""
	// destinationPattern := ""

	if len(flagArgs) > 0 {
		sourcePattern = flagArgs[0]
	}
	fmt.Println("sourcePattern:", sourcePattern)

	//if len(flagArgs) > 1 {
	//	destinationPattern = flagArgs[1]
	//}
	path, pattern := parseSourcePattern(sourcePattern)

	//src := TransferSource{path, pattern}
	//
	//fmt.Println("src-path:", src.path)
	//fmt.Println("src-pattern:", src.pattern)

	preparedPath := strings.Replace(path, "\\", "/", -1)
	preparedPattern := pattern //strings.Replace(pattern, "*", ".*", -1)
	// todo: check if pattern contains groups => (*group1)(*group2), if not, treat whole pattern as group
	// preparedPattern = "(" + preparedPattern + ")"

	preparedPatternToCompile := regexp.QuoteMeta(preparedPath) + "/" + preparedPattern
	fmt.Println("pattern to compile:", preparedPatternToCompile)

	compiledPattern, err := regexp.Compile(preparedPatternToCompile)

	if(err != nil) {
		fmt.Println("invalid source pattern - could not compile")
		return;
	}

	// list := make([]string, 0, 10)

	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		fmt.Println("===================================")
		fmt.Println("path: " + path)

		normalizedPath := strings.Replace(path, "\\", "/", -1)

		fmt.Println("normalized: " + normalizedPath)

		if ! compiledPattern.MatchString(normalizedPath) {
			fmt.Println("match: no")
			return nil
		}
		fmt.Println("match: yes")
		//compiledPattern.ReplaceAllStringFunc(foundFile, func(m string) string {
		//	parts := compiledPattern.FindStringSubmatch(m)
		//	// return parts[1] + complexFunc(parts[2])
		//	// fmt.Println(m)
		//	fmt.Println("0: " + parts[0])
		//	fmt.Println("1: " + parts[1])
		//	//fmt.Println("2: " + parts[2])
		//	return m
		//})


		compiledPattern.ReplaceAllStringFunc(normalizedPath, func(m string) string {
			parts := compiledPattern.FindStringSubmatch(m)
			// return parts[1] + complexFunc(parts[2])
			// fmt.Println(m)
			// fmt.Println("0: " + parts[0])

			i := 0
			for range parts {
				// index is the index where we are
				// element is the element from someSlice for where we are
				fmt.Println("    match: " + parts[i])
				i++
			}

			//fmt.Println("2: " + parts[2])
			return m
		})
		//if info.IsDir() {
		//	return nil
		//}
		//if filepath.Ext(path) == ".sh" {
		//	list = append(list, path)
		//}
		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
		return;
	}

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
