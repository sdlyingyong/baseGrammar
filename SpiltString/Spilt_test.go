package main

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestSpilt(t *testing.T) {
	got := Spilt("abc", "b")
	want := []string{"a", "c"}
	//got == want	//切片应用类型 不能直接判断相等
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v ,but got %v \n", want, got)
		return
	}
}

func TestSpilt2(t *testing.T) {
	got := Spilt("a:b:c", ":")
	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v , but got %v \n", want, got)
	}
}

func TestSpilt3(t *testing.T) {
	//测试用例类型
	type test struct {
		input string
		sep   string
		want  []string
	}
	//各种情况
	tests := []test{
		{input: "abc", sep: "b", want: []string{"a", "c"}},
		{input: "a:b:c", sep: ":", want: []string{"a", "b", "c"}},
	}
	//遍历切片,依次执行测试用例
	for _, v := range tests {
		got := Spilt(v.input, v.sep)
		if !reflect.DeepEqual(got, v.want) {
			t.Errorf(" want %#v, but got %#v \n", v.want, got)
		}
	}
}

func TestSpilt4(t *testing.T) {
	//测试用例类型
	type test struct {
		input string
		sep   string
		want  []string
	}
	tests := map[string]test{
		"simple": {input: "abc", sep: "b", want: []string{"a", "c"}},
		//"wrong":       {input: "a:b:c", sep: ",", want: []string{"a:b:c"}},
		//"more":        {input: "abcd", sep: "bc", want: []string{"a", "d"}},
		//"leading sep": {input: "沙河有沙又有河", sep: "沙", want: []string{"河有", "又有河"}},
	}
	//t.Run() 子测试
	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			got := Spilt(v.input, v.sep)
			if !reflect.DeepEqual(got, v.want) {
				t.Errorf("name:%s failed. want %#v, but got %#v \n", k, v.want, got)
			}
		})
	}
}

//基准测试 最少会运行1秒钟
func BenchmarkSpilt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Spilt("a:b:c:d:e", ":")
	}
}

//性能对比测试
//跑1遍/10遍/100遍 对比每次的耗时
func benchmarkFib(b *testing.B, n int) {
	for i := 0; i < b.N; i++ {
		Fib(n)
	}
}

func BenchmarkFib1(b *testing.B) {
	benchmarkFib(b, 1)
}

func BenchmarkFib2(b *testing.B) {
	benchmarkFib(b, 2)
}

func BenchmarkFib10(b *testing.B) {
	benchmarkFib(b, 10)
}

func BenchmarkFib20(b *testing.B) {
	benchmarkFib(b, 20)
}

//不计时
func BenchmarkFibRetTime(b *testing.B) {
	time.Sleep(time.Second)
	b.ResetTimer() //开始性能测试计时
	Fib(10)
}

func ExampleSpilt() {
	fmt.Println(Spilt("a:b:c", ":"))
}
