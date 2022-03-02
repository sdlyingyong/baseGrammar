package main

import (
	"strings"
)

func main() {
	//需要字符串切割器
	//ret := Spilt("abc", "b")
	//fmt.Println(ret)
	//ret2 := Spilt("bbb", "b")
	//fmt.Println(ret2)
}

//切割字符串
//abc, b  => [a c]
func Spilt(str string, sep string) []string {
	//var ret []string		//没指定空间,下面append占满后就要重新分配内存,耗时*分配次数
	var ret = make([]string, 0, strings.Count(str, sep)+1) //分配内存一次到位 1.ss从两百万词到七百万次执行
	index := strings.Index(str, sep)
	for index >= 0 { //while true 循环
		ret = append(ret, str[:index]) //申请内存
		str = str[index+len(sep):]
		index = strings.Index(str, sep)
	}

	ret = append(ret, str)
	return ret
}

func Fib(n int) int {
	if n < 2 {
		return n
	}
	return Fib(n-1) + Fib(n-2)
}