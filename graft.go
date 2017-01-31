package main
import "fmt"
import (
	"flag"
//	"os"
)

//var wordPtr = flag.String("word", "foo", "a string")
//var numbPtr = flag.Int("numb", 42, "an int")
//var boolPtr = flag.Bool("fork", false, "a bool")
//var svar string

func main() {
	//flag.StringVar(&svar, "svar", "bar", "a string var")

	flag.Parse()

	//fmt.Println("word:", *wordPtr)
	//fmt.Println("numb:", *numbPtr)
	//fmt.Println("fork:", *boolPtr)
	//fmt.Println("svar:", svar)
	//fmt.Println("tail:", flag.Args())

	var sourcePattern = flag.Args()[0]
	// var sourcePattern = flag.Args()[0]

	fmt.Println("sourcePattern:", sourcePattern)
	//if len(flag.Args()) == 0 {
	//	flag.Usage()
	//	os.Exit(1)
	//}
	//fmt.Println(*sourcePattern);
}
