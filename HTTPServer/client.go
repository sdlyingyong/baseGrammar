package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

//公用一个client 用于多次请求,减少开销
var (
	client = http.Client{Transport: &http.Transport{
		DisableKeepAlives: true,
	}}
)

func main() {
	urlExt()
	return
	resp, err := http.Get("http://127.0.0.1:9090/posts/Go/15-socket/")
	if err != nil {
		fmt.Println("get url failed, err :", err)
		return
	}
	//把服务端返回的数据读出来
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read resp.Body failed, err", err)
		return
	}
	fmt.Println(string(b))
}

func urlExt() {
	//准备请求参数
	data := url.Values{}
	data.Set("name", "周林")
	data.Set("age", "9000")
	queryStr := data.Encode()
	//fmt.Println(queryStr)

	//设置请求地址
	urlObj, err := url.ParseRequestURI("http://127.0.0.1:9090/print/")
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	urlObj.RawQuery = queryStr

	//发送请求
	req, err := http.NewRequest("GET", urlObj.String(), nil)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	defer resp.Body.Close()

	//响应结果打印
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	fmt.Println(string(b))
}
