package debugger

import (
	"fmt"
	"os"
)

var debugFile *os.File

func init() {
	var err error
	debugFile, err = os.OpenFile("./debug/debug.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Could not open debug file.")
		panic(err)
	}
}

func Trace(String string) {
	debugFile.WriteString(String + "\n")
}
