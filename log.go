package squid

import (
	"fmt"
	"os"
	"time"
)

var path = "./out.log"

func SetPath(p string) {
	path = p
}

func Log(format string, arg ...interface{}) {
	result := time.Now().Format("\n [2006-01-02 15:04:05] ")
	result += fmt.Sprintf(format, arg...)
	fmt.Print(result)
	if !Exists(path) {
		file, err := os.Create(path)
		defer file.Close()
		if err != nil {
			fmt.Println("Create log file error.")
		}
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0666)
	defer file.Close()
	if err != nil {
		fmt.Println("Log error.")
	}
	file.WriteString(result)
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}