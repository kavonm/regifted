package driver

import (
	"regifted/ts"
	"regifted/util/mylog"

	"regifted/giftcollection"

	"flag"
	"fmt"
	"log"
	"os"
)

const LOGGER_NAME = "driver"
const LOGGER_SEVERITY_LEVEL = mylog.SEV_ERROR

func Main() int {

	logger := mylog.CreateLogger(LOGGER_NAME)
	logger.SetSeverityThresh(LOGGER_SEVERITY_LEVEL)

	filename, rv := getFilepath()
	if rv != 0 {
		os.Exit(rv)
	}
	fmt.Printf("Attempting to read file, Run 7 " + filename + "\n")

	fh, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		// log.Fatal(err)
		panic(err)
	}
	defer func() {
		if err := fh.Close(); err != nil {
			panic(err)
		}
	}()

	pesBuffer := make([]*ts.TSState, 1)
fmt.Println("TESTING1")
	for i := 0; i < 30; i++ {
		fmt.Println( "TESTING i = ", i )
		ts := ts.Load(fh)
		fmt.Println( "TESTING after i" )
		pesBuffer = append(pesBuffer, ts)

	}
fmt.Println("TESTING2")
	giftcollection.Regift(pesBuffer)

	return 0

}



// todo( mathew guest ) I think golang wants to use error as return codes but
// it's a little slow so I'm cheating
func getFilepath() (string, int) {
	flag.Parse()
	argc := flag.NArg()
	if argc < 1 {
		log.Printf("Usage: " + os.Args[0] + " [input ts file]\n")
		return "", 66
	}
	if argc > 1 {
		log.Printf("Ignoring all but first argument.\n")
		os.Exit(1)
	}
	filename := os.Args[1]
	return filename, 0
}
