package config

import (
	"os"
	"spider/log"
	"strconv"
)

var numberOfParallels int

// 读取命令行参数。
func init() {
	length := len(os.Args)
	if length > 1 {
		var temp = os.Args[1]
		if value, err := strconv.Atoi(temp); err == nil {
			numberOfParallels = value
			log.I("Concurrency parameter setting is successful,\r\n"+
				" the current maximum number of concurrent", strconv.Itoa(numberOfParallels))
		}
	} else {
		numberOfParallels = 100
		log.I("Concurrent parameter limit is not set initially，\r\n "+
			"default value", strconv.Itoa(numberOfParallels))
	}
}

func NumberOfParallels() int {
	return numberOfParallels
}
