# Spider 文档

## 项目结构
    spider
        --main.go
        --/cmd/cmd.go 
        --/config/config.go
        --/log/log.go
        --/model/model.go
        --/network/network.go
        --go.mod
        --Readme.md 

## 结构说明
    1.main.go 是 app入口 .负责基础模块以及启动业务.
    2.cmd.go 负责业务逻辑.
    3.log.go 负责日志打印以及日志持久化
    4.config.go 负责读取终端传入参数
    5.model.go 用来处理业务持久化数据(暂未使用)
    6.network.go 负责基础网络层,目前用于网络请求收口,简单代理池创建.
    7.go.mod 项目依赖管理(暂未使用)

## 使用文档

### mac


默认并发数

终端执行

`./spider `

设定并发数

终端执行

`./spider 10`

### windows

默认并发数

cmd 执行

`spider.exe `

设定并发数

cmd执行

`spider.exe 100`

### linux


默认并发数

终端执行

`./spider_linux `

设定并发数

终端执行

`./spider_linux 10`


