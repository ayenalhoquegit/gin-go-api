package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// InfoLogger exported
var InfoLogger *log.Logger

// ErrorLogger exported
var ErrorLogger *log.Logger

func init() {
	absPath, err := filepath.Abs("./log")
	if err != nil {
		fmt.Println("Error reading given path:", err)
	}

	infoLog, err := os.OpenFile(absPath+"/info/info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	errorLog, err := os.OpenFile(absPath+"/error/error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	InfoLogger = log.New(infoLog, "", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(errorLog, "", log.Ldate|log.Ltime|log.Lshortfile)
}
