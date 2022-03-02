package LogExt

import (
	"errors"
	"fmt"
	"path"
	"runtime"
	"strings"
	"time"
)

type LogLevel uint16

//日志结构体
type Logger struct {
	LogLevel LogLevel
}

const (
	UNKNOWN LogLevel = iota
	DEBUG
	TRACE
	INFO
	WARNING
	ERROR
	FATAL
)

//log日志构造函数
func NewLog(levelStr string) Logger {
	level, err := parseLogLevel(levelStr)
	if err != nil {
		panic(err)
	}
	return Logger{LogLevel: level}
}

func parseLogLevel(in string) (LogLevel, error) {
	s := strings.ToLower(in)
	switch s {
	case "debug":
		return DEBUG, nil
	case "trace":
		return TRACE, nil
	case "info":
		return INFO, nil
	case "warning":
		return WARNING, nil
	case "error":
		return ERROR, nil
	case "fatal":
		return FATAL, nil
	default:
		err := errors.New("无效的日志级别")
		return UNKNOWN, err
	}
}

//支持日志文件切割
//每个文件20MB(配置),超过自动开一个新的日志文件,按写入第一条的时间作为文件名+序号

//需要文件大小收集器
//需要切换目标日志文件器
//文件写入器
//



func (l Logger) enable(lv LogLevel) bool {
	return lv <= DEBUG
}

func log(lv LogLevel, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	now := time.Now()
	funcName, fileName, lineNo := getInfo(3)
	fmt.Printf("[%s] [%s] [%v:%v:%v] %s \n", now.Format("2006-01-02 15:04:05"), lv, fileName, lineNo, msg, funcName)
}

func (l Logger) Debug(format string, a ...interface{}) {
	if l.enable(DEBUG) {
		log(DEBUG, format, a...)
	}
}

func getInfo(n int) (funcName, fileName string, lineNo int) {
	pc, file, lineNo, ok := runtime.Caller(1) //网上追溯调用者 0=>1=>2=>...入口
	if !ok {
		fmt.Printf("runtime.caller() failed, err \n")
		return
	}
	funcName = runtime.FuncForPC(pc).Name()
	funcName = strings.Split(funcName, ".")[1]
	fileName = path.Base(file)
	return
}

func (l Logger) Trace(format string, a ...interface{}) {
	if l.enable(TRACE) {
		log(DEBUG, format, a...)
	}

}
func (l Logger) Info(format string, a ...interface{}) {
	if l.enable(INFO) {
		log(DEBUG, format, a...)
	}
}

func (l Logger) Warning(format string, a ...interface{}) {
	if l.enable(WARNING) {
		log(DEBUG, format, a...)
	}
}

func (l Logger) Error(format string, a ...interface{}) {
	if l.enable(ERROR) {
		log(DEBUG, format, a...)
	}
}

func (l Logger) Fatal(format string, a ...interface{}) {
	if l.enable(FATAL) {
		log(DEBUG, format, a...)
	}
}
