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
	result := time.Now().Format("[2006-01-02 15:04:05] ")
	result += fmt.Sprintf(format, arg...)
	fmt.Println(result)
	file, err := os.OpenFile(path, os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Log error.")
	}
	file.WriteString(result)
	file.Close()
}