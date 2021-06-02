package network

import (
	"crypto/tls"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

var clients = make([]http.Client,0)

var uas = make([]string,0)

func init() {
	//预留代理逻辑
	clients = append(clients, createNetworkClient(nil))
	//预留ua逻辑
	uas = append(uas, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Safari/605.1.15")
}

// 统一设定网络请求客户端
func createNetworkClient(proxy func(*http.Request)(*url.URL, error)) http.Client {
	return http.Client{
		Transport: &http.Transport{
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			IdleConnTimeout:       time.Second * time.Duration(30),
			TLSHandshakeTimeout:   time.Second * time.Duration(30),
			ResponseHeaderTimeout: time.Second * time.Duration(30),
			ExpectContinueTimeout: time.Second * time.Duration(30),
			Proxy: proxy,
		},
	}
}

// 收口网络请求
func DoGet(url string) (*http.Response, error) {
	index := rand.Intn(len(clients))
	if request, err := buildRequest(url, http.MethodGet, nil);err == nil{
		return clients[index].Do(request)
	}else{
		return nil,err
	}
}

// 收口请求头
func buildRequest(url, method string, body io.Reader) (*http.Request, error) {
	if request, err := http.NewRequest(method, url, body); err == nil {
		index := rand.Intn(len(uas))
		request.Header.Add("User-Agent",uas[index])
		return request, nil
	} else {
		return request, err
	}
}