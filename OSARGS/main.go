package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"time"
)

func main() {
	getCMDParam()
	showType()
	showTypeVar()
	showPporf()
}

// 一段有问题的代码
func logicCode() {
	var c chan int
	for {
		select {
		case v := <-c:
			fmt.Printf("recv from chan, value:%v\n", v)
		default:
			//time.Sleep(time.Millisecond * 500)

		}
	}
}

func showPporf() {
	var isCPUPprof bool
	var isMemPprof bool

	flag.BoolVar(&isCPUPprof, "cpu", false, "turn cpu pprof on")
	flag.BoolVar(&isMemPprof, "mem", false, "turn mem pprof on")
	flag.Parse()

	if isCPUPprof {
		file, err := os.Create("./cpu.pprof")
		if err != nil {
			fmt.Printf("create cpu pprof failed, err:%v\n", err)
			return
		}
		pprof.StartCPUProfile(file)
		defer pprof.StopCPUProfile()
		//pprof.StartCPUProfile(file)
		//defer pprof.StopCPUProfile()
	}
	for i := 0; i < 8; i++ {
		go logicCode()
	}
	time.Sleep(20 * time.Second)
	//if isMemPprof {
	//	file, err := os.Create("./mem.pprof")
	//	if err != nil {
	//		fmt.Printf("create mem pprof failed, err:%v\n", err)
	//		return
	//	}
	//	pprof.WriteHeapProfile(file)
	//	file.Close()
	//}

	//pprof.WriteHeapProfile(file)
}

func showTypeVar() {
	var name string
	//获取命令行输入 赋值到name 没有则默认值:周林 提示:请输入名字
	flag.StringVar(&name, "name", "周林", "请输入你的名字")
	//解析命令行输入
	flag.Parse()
	fmt.Println(name)
	fmt.Println(flag.Args()) //剩下的命令行参数列
}

func showType() {
	//创建标志位参数
	name := flag.String("name", "张三", "请输入名字")
	age := flag.Int("age", 9000, "请输入真实年龄")
	married := flag.Bool("married", false, "请输入是否结婚")
	mTime := flag.Duration("ct", time.Second, "请输入结婚多久了")
	//使用flag
	flag.Parse()
	fmt.Println(*name)
	fmt.Println(*age)
	fmt.Println(*married)
	fmt.Println(*mTime)
}

func getCMDParam() {
	//flag 获取命令行 附加参数
	fmt.Printf("%#v \n", os.Args)
	fmt.Println(os.Args[0], os.Args[1])
	fmt.Printf("%T \n", os.Args)
}
