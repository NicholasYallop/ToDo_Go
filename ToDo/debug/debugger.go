package debugger

import (
	"fmt"
	"os"
)

var debugFile *os.File

func init() {

}

func InitVariables(main_path string) {
	var err error
	debugFile, err = os.OpenFile(main_path+"\\debug\\debug.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Could not open debug file.")
		panic(err)
	}
}

func Trace(String string) {
	debugFile.WriteString(String + "\n")
}
