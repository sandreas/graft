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
)

//var wordPtr = flag.String("word", "foo", "a string")
//var numbPtr = flag.Int("numb", 42, "an int")
//var boolPtr = flag.Bool("fork", false, "a bool")
//var svar string

// var mapping = list.New()

type TransferPair struct {
	from string
	to string
}

type TransferSource struct {
	path string
	pattern string
}


func parseSourcePattern(sourcePattern string) (string, string) {
	path := sourcePattern
	pattern := ""
	pathExists := false
	for {
		if _, err := os.Stat(path); err == nil {
			pattern = strings.Replace(sourcePattern, path, "", 1)
			pathExists = true
			break
		}
		lastSlashIndex := strings.Index(path, "/")
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

	//if len(flagArgs) > 1 {
	//	destinationPattern = flagArgs[1]
	//}
	path, pattern := parseSourcePattern(sourcePattern)

	src := TransferSource{path, pattern}

	fmt.Println("src:", src)

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
