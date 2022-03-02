package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func f1(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile("./xx.txt")
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%v", err)))
	}
	w.Write(b)
}

func f2(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query()
	fmt.Println(data.Get("name"))
	fmt.Println(data.Get("age"))
	answer := `{"status":"ok"}`
	w.Write([]byte(answer))
}

func main() {
	http.HandleFunc("/posts/Go/15-socket/", f1)
	http.HandleFunc("/print/", f2)
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hello"))
	})
	fmt.Println("start http 127.0.0.1:9090")
	err := http.ListenAndServe("127.0.0.1:9090", nil)
	if err != nil {
		fmt.Println("http server failed,err:", err)
		return
	}
}
