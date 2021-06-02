package cmd

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"spider/config"
	"spider/network"
	"spider/log"
	"sync"
)

var retryCache = make([]string,0)

var lock = sync.RWMutex{}

func Start()  {
	// 模拟url数据 https://tools.applemediaservices.com/app/1458211480?country=cn,
	urls:= mock()
	// 执行任务分配
	grab(urls)
	execRetry()
}

func mock() []string {
	data := make([]string,0)
	data = append(data,"https://tools.applemediaservices.com/app/1458211480?country=cn")
	return data
}

func execRetry() {
	if len(retryCache)>0 {
		retyrList:=make([]string,0)
		retyrList = append(retyrList,retryCache...)
		retryCache = make([]string,0)
		grab(retyrList)
		if len(retryCache)>0 {
			execRetry()
		}
	}
}

func grab(urls []string) {
	length := len(urls)
	// 并发量控制
	var NumberOfParallels int

	// 容错处理,防止channel 缓存大于 设定的并发值,导致程序提前退出
	if length > config.NumberOfParallels(){
		NumberOfParallels = config.NumberOfParallels()
	}else{
		NumberOfParallels = length
	}

	schedule := make(chan int, NumberOfParallels)

	for i := 0; i < length; i++ {
		go worker(Thread{
			url:      urls[i],
			schedule: schedule,
		})

		//控制当前最大并行量
		<-schedule
	}
}

// 爬取
func worker(thread Thread) {
	defer func() {
		if err := recover(); err != nil {
			// 容错兜底,可能不太科学,逻辑不严谨会出现死递归.
			getDataError(thread.url,err)
		}
	}()

	// 网络请求数据
	data := thread.getData()

	// 数据解析
	bundle_id := pars(data)

	log.I(bundle_id)
}

// 试题中的 https://tools.applemediaservices.com/app/1458211480?country=cn
// 请求 返回403 获取不到数据,这里写下业务思路吧
func pars(data []byte) string {
	// 数据解析最通用的办法还是用正则表达式.
	return bytes.NewBuffer(data).String()
}

//发起网络请求,获取html数据
func (thread *Thread) getData() []byte {
	defer thread.done()

	// 使用简单收口的网络客户端发起网络请求
	resp,err := network.DoGet(thread.url)
	if err == nil {
		defer resp.Body.Close()

		// 时间关系,不做细化处理,只处理成功或者失败的 case
		if resp.StatusCode == http.StatusOK {

			if byteData ,err := ioutil.ReadAll(resp.Body);err ==nil {

				// 获取数据成功的 case
				return byteData
			}else{

				// 可能遇到eof
				getDataError(thread.url,err)
				return nil
			}
		}else{

			// 状态码!= 200 时,可能会有服务端异常或者请求参数异常.
			// 通常情况下需要酌情处理,这里简单处理为重试以及保存错误信息,便于复盘.
			getDataError(thread.url,err)
		}
	} else {

		//这里可能是构建的请求出错,也有可能是,网络连接失败.
		getDataError(thread.url, err)
	}
	return nil
}

// 异常 case 的处理.
func getDataError(url string, err interface{}) {

	// 加入到重试列表里,本批次爬取执行结束后会重试 切片不支持并行写,加个锁
	lock.Lock()
	retryCache = append(retryCache,url)
	lock.Unlock()

	//向控制台输出日志,以及写入到日志文件中.
	log.E(url,err)
}

// 释放当前任务的channel 缓存
func (thread *Thread) done()  {
	thread.schedule<-0
}

type Thread struct {
	url      string
	schedule chan int
}
