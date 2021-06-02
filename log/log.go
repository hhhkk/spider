package log

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

var logFile *os.File

var lock = sync.RWMutex{}

func init() {
	if dir, err := filepath.Abs(filepath.Dir(os.Args[0])); err == nil {
		if temp ,err := os.Create(dir+"/"+strconv.FormatInt(time.Now().Unix(),10));err ==nil{
			logFile = temp
		}else{
			// 文件系统交互异常，暂定直接停止运行,日志是目前仅有的复盘手段
			panic(err)
		}
	} else{
		// 文件系统交互异常，暂定直接停止运行,日志是目前仅有的复盘手段
		panic(err)
	}
}

func E(err ...interface{}){
	if err!=nil {
		cacheLog(err)
		I(err)
	}
}

func cacheLog(err []interface{}) {
	// 对象序列化成字符串 可能格式不太友好.
	if data, err := json.Marshal(err); err == nil {
		// 并发环境下,加锁
		lock.Lock()
		logFile.WriteString(time.Now().String())
		logFile.WriteString("\r\n")
		logFile.Write(data)
		logFile.WriteString("\r\n\r\n")
		lock.Unlock()
	}
}

func I(message ...interface{}){
	length := len(message)
	// 日至格式好看点
	for i:=0;i< length;i++ {
		fmt.Println(message[i])
	}
}


func Close()  {
	if logFile!=nil {
		logFile.Close()
		logFile = nil
	}
}

