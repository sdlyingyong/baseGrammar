package main

import (
	"basicGrar/LogExt"
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"unicode"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/nsqio/go-nsq"
)

const (
	n1 = iota //0
	n2 = 100  //100
	n3 = iota //2
	n4        //3
)
const n5 = iota //0

func main() {
	if false {

		showIota()
		showFmtNum()
		showFloat()
		showBool()
		showStr()

		showIf()
		showFor()
		multiTable()
		showBreak()
		showGoto()
		showSwitch()

		showArray()
		showSlice()
		showAppend()
		showCopy()
		showPoint()
		showMap()

		showFunc()
		showStrNum()
		showFunc2()
		showFunc3()
		showFuncAsPam()
		showFuncAsRet()
		showNoNameFunc()
		showSendFunc()
		showCliseFunc()

		showDeferExec()
		showDefer()
		showDeferRegister()

		showRecover()
		showFmt()
		showDispatchCoin()
		showRecursive()
		showContruct()
		showReceiveVal()

		showStutAdmnByFunc()
		showNoNameConst()
		showJson()
		showStutAdmnByObj()

		showInterfaceDog()
		showInteChin()
		showInterfaceType()

		showFileReadByIO()
		showFileReadByBufio()
		showFileReadByIOutil()
		showWriteFileByOs()
		showWriteFileByBufio()
		showWriteFileByIoutil()
		showCopyFile()
		showWriteOffset()
		showInsertFile()

		showTime()
		showGetInfo(1)
		showSliceDemo()
		showLogExt()

		showReflect()
		showReflectSet()
		showReflectBool()
		showReflectMethod()
		showReflectJson()
		showWorkIniLoad()
		showGoIni()

		showStrToInt()
		showStrconv()

		showGoroutine()
		showRandInt()
		showWaitGroup()
		showGPM()

		showChan()
		showWorkChan()
		showChanClose()
		showOnlyChan()
		showWorkGRPool()
		showWorkGRNum()
		showSelect()

		showChangeParamFunc()
		showSyncMutex()
		showRWLock()
		showSyncOnce()
		showSyncMap()
		showAtomicAdd()
		showAtomic()
		showSyncPool()

		showTCPServer()
		showTCPClient()
		showTCPStickyBuns()

		showLinkListRepeat()
		showUpStep(10)

		showMysqlConn()
		showMysqlSelect()
		showMysqlInsert()
		showMysqlDelete()
		showMysqlUpdate()
		showMysqlInsertPrepare()
		showMysqlTrans()

		showSqlx()
		showRedis()
		showRedisZset()

		showNsqProducer()
		showNsqConsumer()

		showGoRoutineExit()
		showContext()
	}

}

func sonWorkerCont(ctx context.Context) {
LOOP:
	for {
		fmt.Println("worker son")
		time.Sleep(time.Second)
		select {
		case <-ctx.Done():
			break LOOP
		default:
		}
	}
}

func workerCont(ctx context.Context) {
	go sonWorkerCont(ctx)
LOOP:
	for {
		fmt.Println("worker")
		time.Sleep(time.Second)
		select {
		case <-ctx.Done():
			break LOOP
		default:
		}
	}
	wg.Done()
}

//按照时间关闭产生的所有子goroutine
func showContext() {
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go workerCont(ctx)
	//go workerCont(context.TODO())
	time.Sleep(time.Second * 3)
	cancel() //通知结束子goroutine
	wg.Wait()
	fmt.Println("stop all goroutine ... over")
}

var exitChanCont = make(chan bool, 1)

func fCont() {
	defer wg.Done()
FORLOOP:
	for {
		select {
		case <-exitChanCont:
			break FORLOOP
		default:
			fmt.Println("周林")
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func showGoRoutineExit() {
	wg.Add(1)
	go fCont()
	time.Sleep(time.Second * 5)
	exitChanCont <- true
	close(exitChanCont)
	wg.Wait()
}

type MyHandle struct {
	Title string
}

//HandleMessage(message *Message) error
func (m *MyHandle) HandleMessage(msg *nsq.Message) (err error) {
	fmt.Println("msg: ", m.Title, msg.NSQDAddress, string(msg.Body))
	return
} //顺序: 1.返回值赋值 2.defer做完 3.跳转返回

func initConsumer(topic, channel, address string) (err error) {
	config := nsq.NewConfig()
	config.LookupdPollInterval = 15 * time.Second
	c, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		fmt.Println("new consumer fail, err :", err)
		return
	}
	consumer := &MyHandle{Title: "沙河一号"}
	c.AddHandler(consumer)
	if err := c.ConnectToNSQLookupd(address); err != nil {
		return err
	}
	return nil
}

func showNsqConsumer() {
	err := initConsumer("topic_demo", "first", "127.0.0.1:4161")
	if err != nil {
		fmt.Println("init consumer fail, err :", err)
		return
	}
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT)
	<-c
}

var nsqProd *nsq.Producer

func initNsqProd(str string) (err error) {
	config := nsq.NewConfig()
	nsqProd, err = nsq.NewProducer(str, config)
	if err != nil {
		fmt.Println("init nsq fail, err :", err)
		return err
	}
	return nil
}

func showNsqProducer() {
	nsqAddr := "127.0.0.1:4150"
	err := initNsqProd(nsqAddr)
	if err != nil {
		fmt.Println("initNsqProd fail, err :", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		data, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("read string from stdin fail, err :", err)
			continue
		}
		data = strings.TrimSpace(data)
		if strings.ToUpper(data) == "Q" {
			break
		}
		err = nsqProd.Publish("topic_demo", []byte(data))
		if err != nil {
			fmt.Println("prod fail,err :", err)
			continue
		}
	}
}

//用来按照分数排名的结构 score =>分数 member=>对象
func showRedisZset() {
	initRedis()
	key := "rank"
	items := []redis.Z{
		redis.Z{Score: 90, Member: "PHP"},
		redis.Z{Score: 80, Member: "JAVA"},
		redis.Z{Score: 70, Member: "Go"},
	}
	num, err := redisC.ZAdd(key, items...).Result()
	if err != nil {
		fmt.Printf("zset err : %v \n", err)
		return
	}
	fmt.Println("zset num :", num)

	//获取从高到低的排名
	//ZREVRANGEBYSCORE redis命令
	retS, err := redisC.ZRevRangeByScore(key, redis.ZRangeBy{"60", "100", 0, 10}).Result()
	if err != nil {
		fmt.Println("get by score err: ", err)
		return
	}
	fmt.Println("按照热度排序的列表是: ", retS)
}

var redisC *redis.Client

func initRedis() (err error) {
	redisC = redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
	_, err = redisC.Ping().Result()
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	fmt.Println("redis 连接成功")
	return nil
}

func showRedis() {
	err := initRedis()
	if err != nil {
		fmt.Println("err: ", err)
	}
	//redis exe

}

type LogMS struct {
	Id         int
	Mark       string
	Created_at []uint8
	Updated_at []uint8
}

func showSqlx() {
	dsn := "root:root@tcp(127.0.0.1:3306)/metal"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Println("conn mysql err: ", err)
		return
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	//查询
	sqlStr := "select id,mark,created_at,updated_at from log where id = ?"
	var log1 LogMS
	err = db.Get(&log1, sqlStr, 2)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	fmt.Println("查询到记录", log1.Id, log1.Mark, log1.Created_at)

}

func showMysqlTrans() {
	db, err := initMysql()
	if err != nil {
		fmt.Println("err :", err)
		return
	}
	defer db.Close()

	//开启事务
	tx, err := db.Begin()
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		fmt.Println("start trans err: ", err)
		return
	}
	//遇到异常回滚操作
	sqlStr := "update log set mark = ? where id = ?"
	randMark := rand.Intn(100)
	ret, err := db.Exec(sqlStr, randMark, 2)
	if err != nil {
		tx.Rollback()
		fmt.Println("err: ", err)
		return
	}
	affRow1, err := ret.RowsAffected()
	if err != nil {
		fmt.Println("err:", err)
		tx.Rollback()
		return
	}
	if affRow1 != 1 {
		tx.Rollback()
		fmt.Println("未修改")
		return
	}
	//提交事务
	tx.Commit()
	fmt.Println("事务提交成功")
}

//预处理形式执行语句
//1.发送语句给mysql,给他预编译 2.发送需要的参数 3.mysql匹配参数和语句,进行执行.
//作用: 1.防止sql注入攻击 2.减少每次编译sql语句的消耗
func showMysqlInsertPrepare() {
	db, err := initMysql()
	if err != nil {
		fmt.Println("err :", err)
		return
	}
	defer db.Close()

	sqlStr := "INSERT INTO `log` ( `mark`, `created_at`, `updated_at`) VALUES ( ?, ?, ?)"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	now := time.Now()
	m := map[string]time.Time{
		"1": now,
		"2": now,
		"3": now,
	}
	defer stmt.Close()

	for k, v := range m {
		//stmt.Exec("1",now,now)
		_, err := stmt.Exec(k, v, v)
		if err != nil {
			fmt.Println("err: ", err)
			return
		}
	}
	fmt.Println("批量预处理形式插入成功!")

}

func showMysqlUpdate() {
	db, err := initMysql()
	if err != nil {
		fmt.Println("err :", err)
		return
	}
	defer db.Close()

	sqlStr := "update log set mark = ? where id = ?"
	ret, err := db.Exec(sqlStr, rand.Intn(100), 2)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	n, err := ret.RowsAffected()
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	fmt.Printf("更新了%v行", n)
}

func showMysqlDelete() {
	db, err := initMysql()
	if err != nil {
		fmt.Println("err :", err)
		return
	}
	defer db.Close()

	sqlStr := "delete from log where id = ?"
	ret, err := db.Exec(sqlStr, 1)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	n, err := ret.RowsAffected()
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	fmt.Printf("删除了%v行记录", n)
}
func showMysqlInsert() {
	db, err := initMysql()
	if err != nil {
		fmt.Println("err: ", err)
	}
	defer db.Close()
	sqlStr := "INSERT INTO `log` ( `mark`, `created_at`, `updated_at`) VALUES ( ?, ?, ?)"
	ret, err := db.Exec(sqlStr, 1, time.Now(), time.Now())
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	newId, _ := ret.LastInsertId()
	fmt.Println("插入成功,id为 ", newId)
}

func showMysqlSelect() {
	db, err := initMysql()
	if err != nil {
		fmt.Println("err: ", err)
	}
	//在err之后注册关闭函数,因为资源创建失败时就关闭会报错
	defer db.Close()
	db.SetMaxOpenConns(10) //最大建立连接数
	db.SetMaxIdleConns(5)  //最大闲置连接数

	//查询单条
	var art1 art
	//1.查询语句
	sqlStr := "select id,title,content from article"
	//2.执行
	rowObj := db.QueryRow(sqlStr)
	//rowObj := db.QueryRow(sqlStr).Scan(&art1.id, &art1.title, art1.content)
	//3.拿到结果
	err = rowObj.Scan(&art1.id, &art1.title, &art1.content) //必须对obj对象调用scan,该方法会释放数据库连接
	if err != nil {
		fmt.Println("err: ", err)
	}
	fmt.Println(art1.title)

	//查询多条
	sqlStr = "select id,title,content from article where id > ?"
	rows, err := db.Query(sqlStr, 10)
	if err != nil {
		fmt.Println("err :", err)
		return
	}
	//关闭rows持有的数据库连接
	defer rows.Close()

	list := make(map[string]art)
	for rows.Next() {
		err := rows.Scan(&art1.id, &art1.title, &art1.content)
		if err != nil {
			fmt.Println("err: ", err)
			return
		}
		list[art1.id] = art1
		//fmt.Printf(" 查询结果 id:%v title:%v content: %v \n",art1.id,art1.title,art1.content)
	}
	fmt.Printf("得到记录%v条", len(list))
}

type art struct {
	id      string
	title   string
	content string
}

func initMysql() (db *sql.DB, err error) {
	dsn := "root:root@tcp(127.0.0.1:3306)/metal"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		errors.New("mysql 无连接")
	}
	fmt.Println("数据库连接成功")
	return db, nil
}

var dbG *sql.DB

func showMysqlConn() {
	initMysql()
}

type a struct {
	val  int
	next *a
}

func showUpStep(n int) int {
	if n == 1 {
		return 1
	}
	if n == 2 {
		return 2
	}
	return showUpStep(n-1) + showUpStep(n-2)
}

func showLinkListRepeat() {
	//x 走一步
	//y走两步
	//如果某一时刻 他们在同一个节点相遇
	//说明链表有闭环
}

func DecodeSB(reader *bufio.Reader) (string, error) {
	//读取消息长度
	lengthByte, _ := reader.Peek(4)
	lengthBuff := bytes.NewBuffer(lengthByte)
	var length int32
	err := binary.Read(lengthBuff, binary.LittleEndian, &length)
	if err != nil {
		return "", err
	}
	//buffer返回缓冲中现有的可读取字节数
	if int32(reader.Buffered()) < length+4 {
		return "", err
	}
	//读取真正的消息数据
	pack := make([]byte, int(4+length))
	_, err = reader.Read(pack)
	if err != nil {
		return "", err
	}
	return string(pack[4:]), nil
}

func EncodeSB(msg string) ([]byte, error) {
	//读取消息,转为int32
	length := int32(len(msg))
	pkg := new(bytes.Buffer)
	//写入消息头
	err := binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}
	//写入消息实体
	err = binary.Write(pkg, binary.LittleEndian, []byte(msg))
	if err != nil {
		return nil, err
	}
	//返回
	return pkg.Bytes(), nil
}

func processSB(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	//var buf [1024]byte
	for {
		//n, err := reader.Read(buf[:])
		recvStr, err := DecodeSB(reader)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("read from client failed, err:", err)
			break
		}
		//recvStr := string(buf[:n])
		fmt.Println("收到client发来的数据：", recvStr)
	}
}

func showTCPStickyBuns() {
	listen, err := net.Listen("tcp", "127.0.0.1:30000")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	defer listen.Close()
	fmt.Println("开启服务", "tcp", "127.0.0.1:30000")
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		go processSB(conn)
	}
}

func showTCPClient() {
	//1.与server建立连接
	conn, err := net.Dial("tcp", "127.0.0.1:20000")
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	//2,通信
	conn.Write([]byte("hello wang ye!"))
	conn.Close()
}

func processConn(conn net.Conn) {
	//与客户端通信
	var tmp [128]byte
	n, err := conn.Read(tmp[:])
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	fmt.Println(string(tmp[:n]))
}

func showTCPServer() {
	//1.本地服务端口启动
	listen, err := net.Listen("tcp", "127.0.0.1:20000")
	if err != nil {
		fmt.Println("127.0.0.1:20000 tcp start error:", err)
		return
	}
	//2.等待别人来和我建立连接
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("err: ", err)
			return
		}
		go processConn(conn)
	}
}

//存放对象的内存池,减少内存占用,无引用自动回收减轻cg压力
func showSyncPool() {
	//放在内存池,无引用自动回收
	pool := &sync.Pool{New: func() interface{} {
		return new(person)
	}}
	//从池中拿到对象
	p1 := pool.Get().(*person)
	fmt.Println("获取了内存池中的对象 p1 :", p1)
	p1.age = 20
	//用完还回内存池
	pool.Put(p1)
	fmt.Println("把对象放回内存池中")
	//拿到空
	fmt.Println("pool拿到对象 person :", pool.Get().(*person))
	fmt.Println("pool里的对象空了,调用get : ", pool.Get().(*person))
}

func showAtomic() {
	x64 = int64(200)
	//如果当前是a值,就改成b值
	ok := atomic.CompareAndSwapInt64(&x64, 200, 300)
	fmt.Println(ok, x64)
}

var x64 int64

func addAc() {
	defer wg.Done()
	//lock.Lock()
	//x = x + 1
	//lock.Unlock()
	//原子操作 加
	atomic.AddInt64(&x64, 1)
}

func showAtomicAdd() {
	wg.Add(100000)
	for i := 0; i < 100000; i++ {
		go addAc()
	}
	wg.Wait()
	fmt.Println(x64)
}

var m = make(map[string]int)
var mS sync.Map

func getSM(key string) (value interface{}) {
	//return m[key]
	ret, _ := mS.Load(key)
	return ret
}
func setSM(key string, value int) {
	//m[key] = value
	mS.Store(key, value)
}
func showSyncMap() {
	wg := sync.WaitGroup{}
	//普通map 20个并发就会报错  fatal error: concurrent map writes
	//sync.Map 可以支持并发读写
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func(n int) {
			key := strconv.Itoa(n)
			setSM(key, n)
			fmt.Printf("k=%v,v:=%v \n", key, getSM(key))
			wg.Done()
		}(i)
	}
	wg.Wait()
}

var once sync.Once

func showSyncOnce() {
	f := func() {
		add(1, 2)
	}
	once.Do(f) //传入闭包函数
}

var rwlock sync.RWMutex

func read() {
	defer wg.Done()
	//lock.Lock()
	rwlock.RLock()
	fmt.Println(x)
	time.Sleep(time.Millisecond)
	//lock.Unlock()
	rwlock.RUnlock()
}
func write() {
	defer wg.Done()
	//lock.Lock()
	rwlock.Lock()
	x = x + 1
	time.Sleep(time.Millisecond * 5)
	//lock.Unlock()
	rwlock.Unlock()
}

func showRWLock() {
	start := time.Now()
	for i := 0; i < 100; i++ {
		go write()
		wg.Add(1)
	}
	for i := 0; i < 100; i++ {
		go read()
		wg.Add(1)
	}
	wg.Wait()
	fmt.Println(time.Now().Sub(start))
	//sync.mutex 3.13s 1.57s 缩短了50%时间
}

var x = 0
var lock sync.Mutex

func addSc() {
	lock.Lock()
	for i := 0; i < 50000; i++ {
		x = x + 1
	}
	lock.Unlock()
	wg.Done()
}

func showSyncMutex() {
	wg.Add(2)
	go addSc()
	go addSc()
	wg.Wait()
	fmt.Println(x)
}

//可变参数
func f1CP(a ...interface{}) { //允许interface{}类型 n种数量的参数传入
	fmt.Printf("type:%T value:%#v \n", a, a)
}

func showChangeParamFunc() {
	var s = []interface{}{1, 2, 3, 4}
	f1CP(s)
	f1CP(s...) //拆开,把每个元素传入
	f1CP()
	f1CP(1)
	f1CP(1, 2, 3, 4, 5, 6)
	f1CP(1, 2, 3, 4, 5, 6)
	f1CP(1, 2, 3, 4, 5, false, "a", []int{1, 2}, [...]int{12, 3}, map[string]int{"ms": 1})
}

func showSelect() {
	ch := make(chan int, 10)
	for i := 0; i < 10; i++ {
		select {
		case x := <-ch:
			fmt.Println(x)
		case ch <- i:
		}
	}
}

type jobGr struct {
	value int64
}
type resultGr struct {
	jobGr *jobGr
	sum   int64
}

var jobChan = make(chan *jobGr, 100)
var resultChan = make(chan *resultGr, 100)

func zhoulin(zl chan<- *jobGr) {
	for {
		x := rand.Int63()
		newJob := &jobGr{x}
		zl <- newJob
		time.Sleep(time.Millisecond * 500)
	}
}

func baodelu(zl <-chan *jobGr, resultChan chan<- *resultGr) {
	defer wg.Done()
	for {
		job := <-zl
		sum := int64(0)
		n := job.value
		for n > 0 {
			sum += n % 10
			n = n / 10
		}
		newRet := &resultGr{
			jobGr: job,
			sum:   sum,
		}
		resultChan <- newRet
	}

}

func showWorkGRNum() {
	wg.Add(1)
	go zhoulin(jobChan)
	//开启24个goroutine执行
	for i := 0; i < 24; i++ {
		go baodelu(jobChan, resultChan)
	}
	for result := range resultChan {
		fmt.Printf("value:%d sum:%d \n", result.jobGr.value, result.sum)
	}
	wg.Wait()
}

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Sprintf("worker:%d start job:%d \n", id, j)
		time.Sleep(time.Second)
		fmt.Printf("worder:%d end job:%d \n", id, j)
		results <- j * 2
	}
}

func showWorkGRPool() {
	//需要并发池
	jobs := make(chan int, 100)
	results := make(chan int, 100)
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)
	for a := 1; a <= 5; a++ {
		<-results
	}
}

func onlySendChan(ch1 <-chan int) {
}
func onlyReceiveChan(ch1 chan<- int) {
}
func showOnlyChan() {
	ch1 := make(chan int, 10)
	onlySendChan(ch1)
	onlyReceiveChan(ch1)
}

func showChanClose() {
	ch1 := make(chan int, 2)
	ch1 <- 10
	ch1 <- 20
	close(ch1)
	//chan用for range没有值会退出
	//for x := range ch1 {
	//	fmt.Println(x)
	//}

	<-ch1
	<-ch1
	x, ok := <-ch1
	//取值的chan已关闭,返回零值和ok=false
	fmt.Println(x, ok)
}

var sOnce sync.Once

func f1ch(ch1 chan int) {
	for i := 0; i < 100; i++ {
		ch1 <- i
	}
	close(ch1)
	wg.Done()
}

func f2ch(ch1, ch2 chan int) {
	for {
		x, ok := <-ch1
		if !ok {
			break
		}
		ch2 <- x * x
	}
	sOnce.Do(func() {
		close(ch2)
	})
	wg.Done()
}

func showWorkChan() {
	a := make(chan int, 100)
	b := make(chan int, 100)
	wg.Add(3)
	go f1ch(a)
	go f2ch(a, b)
	go f2ch(a, b)
	wg.Wait()
	for ret := range b {
		fmt.Println(ret)
	}
}

func showChan() {
	//chan 需要make分配空间才能使用
	msgCn := make(chan string)
	//对无缓冲无接收的chan发送会报错
	//all goroutines are asleep - deadlock!
	//msgCn <- "你好"
	fmt.Println(msgCn)

	//通过chan通信
	b := make(chan int, 1)
	//defer close(b)	//会自动关闭
	b <- 10
	fmt.Println(b)
	wg.Add(1)
	go func(b <-chan int) {
		fmt.Println(<-b)
		wg.Done()
	}(b)
	wg.Wait()
}

func aGMP() {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		fmt.Printf("A: %d \n", i)
	}
}

func bGMP() {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		fmt.Printf("B: %d \n", i)
	}
}

func showGPM() {
	runtime.GOMAXPROCS(runtime.NumCPU()) //默认是最大 等于本机cpu核心数
	wg.Add(2)
	go aGMP()
	go bGMP()
	wg.Wait()
}

var wg = sync.WaitGroup{}

func f1GR(in int) {
	defer wg.Done()
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
	fmt.Println(in)
}

func showWaitGroup() {
	wg.Add(3) //设置计数器 3
	for i := 0; i < 10; i++ {
		go f1GR(i)
		//wg.Done()		//执行完计数器 -1
	}
	wg.Wait() //等待wg的计数器减为 0
}

func showRandInt() {
	rand.Seed(time.Now().UnixNano()) //每次使用新的随机数种子 int64
	for i := 0; i < 5; i++ {
		r1 := rand.Int()
		r2 := rand.Intn(10)
		fmt.Println(0-r1, 0-r2)
	}
}

func hello(in int) {
	fmt.Println("hello", in)
}

func showGoroutine() {
	for i := 0; i < 1; i++ {
		//go hello(i)
		//匿名函数
		go func(i int) {
			fmt.Println(i)
		}(i)
	}
	fmt.Println("main")
	time.Sleep(1 * time.Second)
}

func showStrconv() {
	//从字符串解析出数字
	str := "10000"
	//ret1 := int64(str)	//会报错的操作
	ret1, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%T %#v \n", ret1, ret1)
	retInt, _ := strconv.Atoi(str)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%T %#v \n", retInt, retInt)

	//把数字转换成字符串类型
	i := int32(97)
	//ret2 := string(i)
	ret2 := fmt.Sprintf("%d", i)
	fmt.Printf("%T %#v \n", ret2, ret2)

	//解析字符串bool值
	boolStr := "true"
	boolValue, _ := strconv.ParseBool(boolStr)
	fmt.Printf("%T %v \n", boolValue, boolValue)

	//解析字符串float
	floatStr := "12.13"
	floatValue, _ := strconv.ParseFloat(floatStr, 64)
	fmt.Printf("%T %v \n", floatValue, floatValue)

	intStr := strconv.Itoa(100)
	fmt.Printf("%T %v \n", intStr, intStr)

}

func showGoIni() {
}

func showStrToInt() {
	str := "10"
	i64, _ := strconv.ParseInt(str, 10, 64)
	fmt.Printf("%T %v", i64, i64)
}

type MysqlConfig struct {
	Address  string `ini:"address"`
	Port     string `ini:"port"`
	UserName string `ini:"username"`
	Password string `ini:"password"`
}

type RedisConfig struct {
	Host     string `ini:"host"`
	Port     string `ini:"port"`
	Username string `ini:"username"`
	Password string `ini:"password"`
}

type Config struct {
	MysqlConfig `ini:"mysql"`
	RedisConfig `ini:"redis"`
}

func showWorkIniLoad() {
	var cfg Config
	var x = new(int)
	err := iniLoad("./config.ini", &cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cfg, x)
}

func iniLoad(fileName string, data interface{}) (err error) {
	//todo 实现加载配置文件
	t := reflect.TypeOf(data)
	if t.Kind() != reflect.Ptr {
		err = fmt.Errorf("data should be a pointer")
		return
	}
	if t.Elem().Kind() != reflect.Struct {
		err = fmt.Errorf("data should be a pointer")
		return
	}
	//需要文件读取器
	fmt.Println(fileName)
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	lineSlice := strings.Split(string(b), "\r\n")
	fmt.Printf("%#v \n", lineSlice)
	//需要一行行读取器
	for idx, line := range lineSlice {
		//如果是注释就跳过
		line = strings.TrimSpace(string(line))
		if strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}
		//如果[在开头就是节
		if strings.HasPrefix(line, "[") {
			if line[0] != '[' || line[len(line)-1] != ']' {
				err = fmt.Errorf("line: %d syntax error", idx+1)
				return
			}
			//取消收尾的空行
			sectionName := strings.TrimSpace(line[1 : len(line)-1])
			if len(strings.TrimSpace(line[1:len(line)-1])) == 0 {
				err = fmt.Errorf("line:%d syntax error", idx+1)
				return
			}
			//根据字符串section去data里面找对应结构体
			v := reflect.ValueOf(data)
			fmt.Println(v, sectionName)
			//for i := 0; i < v.NumField(); i++ {
			//	field := t.Field(i)
			//	if sectionName == field.Tag.Get("ini"){
			//		structName := field.Name
			//		fmt.Println("找到%s对应的嵌套结构体:%s \n",sectionName, structName)
			//	}
			//}

		} else {
			//如果不是 [ 开头  就是=分隔的键值对
			//=分隔一行,左边是key 右边是value
			//if strings.Contains(line,"=") {
			if strings.Index(line, "=") == -1 || strings.HasPrefix(line, "=") {
				err = fmt.Errorf("line:%d syntax error", idx+1)
				return
			}
			index := strings.Index(line, "=")
			key := strings.TrimSpace(line[:index])
			value := strings.TrimSpace(line[index+1:])
			//根据structName去结构体取出字段值
			v := reflect.ValueOf(data)
			fmt.Println(key, value, v)
			//structObj := v.Elem().FieldByName(structName)
			//if structObj.Kind() != reflect.Struct{
			//	fmt.Println("data中的%s字段应该是一个结构体",structName)
			//}
			//// 如果key = tag,给这个字段赋值
			//for i := 0;i < structObj.NumField();i++ {
			//	//field := sType.Field(i)
			//
			//}
		}

	}

	//需要配置映射器
	//检查data,必须是指针类型,因为要对其赋值
	return nil

}

type studentJson struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

func reflectField(in interface{}) {
	t := reflect.TypeOf(in)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fmt.Printf("name:%s index:%d type:%v json tag:%v \n", field.Name, field.Index, field.Type, field.Tag.Get("json"))
	}

	if scoreField, ok := t.FieldByName("Score"); ok {
		fmt.Printf("name:%s index:%d type:%v json tag:%v \n", scoreField.Name, scoreField.Index, scoreField.Type, scoreField.Tag.Get("json"))
	}
}

func showReflectJson() {
	stu1 := studentJson{Name: "小王子", Score: 90}
	reflectField(stu1)
}

func showReflectMethod() {
	//通过反射调用值
	s := student{}
	reflectMethod(s)
}

func (s student) Study() string {
	msg := "好好学习 天天向上"
	fmt.Println(msg)
	return msg
}
func (s student) Sleep() string {
	msg := "好好睡觉,按时长大"
	fmt.Println(msg)
	return msg
}

func reflectMethod(in interface{}) {
	t := reflect.TypeOf(in)
	v := reflect.ValueOf(in)

	fmt.Println(t.NumMethod())
	for i := 0; i < v.NumMethod(); i++ {
		methodType := v.Method(i).Type()
		fmt.Printf("method name:%s \n", t.Method(i).Name)
		fmt.Printf("method: %s \n", methodType)
		var args = []reflect.Value{}
		v.Method(i).Call(args)
	}
}

func showReflectBool() {
	var a *int
	//判断指针是否为空
	fmt.Println("var a *int isNil ", reflect.ValueOf(a).IsNil())
	//判断返回值是否有效
	fmt.Println("var nil isNil ", reflect.ValueOf(nil).IsValid())
	b := struct{}{}
	fmt.Println("不存在的结构体成员", reflect.ValueOf(b).FieldByName("abc").IsValid())
	fmt.Println("不存在的结构体方法", reflect.ValueOf(b).MethodByName("abc").IsValid())
	//map
	c := map[string]int{}
	fmt.Println("map中不存在的键: ", reflect.ValueOf(c).MapIndex(reflect.ValueOf("娜扎")).IsValid())
}

func showReflectSet() {
	var a int64 = 100
	reflectSetValue(&a)
	fmt.Println(a)
}

func reflectSetValue(in interface{}) {
	v := reflect.ValueOf(in)
	if v.Elem().Kind() == reflect.Int64 {
		v.Elem().SetInt(200) //设置新值
	}
}

func showReflect() {
	str := `{"name":"周林","age":9000}`
	var p person
	json.Unmarshal([]byte(str), &p)
	fmt.Println(p.name, p.age)

	var a *float32 //指针类型 reflect name 是空的
	var b myInt
	var c rune
	reflectType(a) //type:  kind:ptr
	reflectType(b) //type: myInt kind:int
	reflectType(c) //type: int32 kind:int32

	//array & slice	diff
	arr := [...]int{1, 2, 3}
	s1 := []int{1, 2, 3}
	reflectType(arr) //type:  kind:array
	reflectType(s1)  //type:  kind:slice

	p2 := person{name: "zhangsan"}
	reflectType(p2) //type: person kind:struct

	type book struct {
		name string
	}
	bk2 := book{name: "钢铁是怎样炼成的"}
	reflectType(bk2) //type: book kind:struct

	//数组、切片、Map、指针	reflect 的name是空的
}

func reflectType(x interface{}) {
	v := reflect.TypeOf(x)
	fmt.Printf("type: %v kind:%v \n", v.Name(), v.Kind())
}

func showSliceDemo() {
	s1 := make([]int, 0, 10)
	_ = append(s1, 1, 2, 3, 4)
	fmt.Println(s1)
}

//日志库
//支持往不同地方输出日志
//按照级别来分类提示 debug trace info warning error fatal
//日志支持开关控制 上线只有info往下级别才能输出
//完整的日志记录 需要有时间 行号 文件名 日志级别 日志信息
//
func showLogExt() {
	//需要文件输出器
	//需要日志分级器
	log := LogExt.NewLog("debug")
	for i := 0; i < 3; i++ {
		log.Debug("这是一条debug日志", 123)
		//log.Info("这是一条 Info 日志")
		//log.Error("这是一条 Error 日志")
	}
}

//调用者信息
//打印调用者 文件 行数
func showGetInfo(n int) {
	pc, file, line, ok := runtime.Caller(1) //网上追溯调用者 0=>1=>2=>...入口
	if !ok {
		fmt.Printf("runtime.caller() failed, err \n")
	}
	//return pc,file,line
	fmt.Println(pc)
	fmt.Println(file)
	fmt.Println(line)
}

func showTime() {
	//时间类型
	now := time.Now()
	fmt.Println(now.String())
	fmt.Println(now.Year())
	fmt.Println(now.Month())
	fmt.Println(now.Day())
	fmt.Println(now.Hour())
	//%02d 按照2位补齐 不足的前面填充0
	fmt.Printf("现在时间是北京时间%04v年%02d月%02v日%02v小时%02v分%02v秒 \n", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())

	//获取时间戳
	timeStap := now.Unix()
	fmt.Printf("unix : %v\n ", timeStap)
	//纳秒时间戳 nano
	//nanoSec < microSec <  milliSec <  sec < min < hour
	timeStap2 := now.Nanosecond()
	fmt.Println(timeStap2)
	fmt.Printf("nanoSecond : %v\n", timeStap2)
	timeStap3 := now.UnixNano()
	fmt.Printf("unixNano :%v \n", timeStap3)

	//时间 => 字符串
	fmt.Println(now.Format("2006-01-02 15:04:05")) //06 12345
	//0.000 毫秒值
	fmt.Println(now.Format("2006-01-02 15:04:05.000")) //06 12345
	//xxx AM
	fmt.Println(now.Format("2006-01-02 15:04:05 PM")) //06 12345

	//字符串 => 时间
	timeObj, err := time.Parse("2006-01-02", "2022-02-11") //layout value
	if err != nil {
		fmt.Printf("parse time failed, err %v", err)
		return
	}
	fmt.Println(timeObj)

	//按照时区解析字符串
	//需要时区加载器
	loc, err := time.LoadLocation("US/Hawaii")
	if err != nil {
		fmt.Println(err)
		return
	}
	//需要转换器
	timeObj, err = time.ParseInLocation("2006/01/02 15:04:05", "2022/02/11 15:33:00", loc)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("按照时区解析字符串", timeObj)
	fmt.Println(timeObj.Sub(now))

	//一小时后
	later := now.Add(time.Hour)
	fmt.Println(later)

	//时间差
	subT := now.Sub(later)
	fmt.Println(subT)

	//等于 时区和时间都检查
	isSameTime := now.Equal(later)
	fmt.Println(isSameTime)
	tmpT, err := time.Parse("20060102", "20220211")
	if err != nil {
		fmt.Println(err)
		return
	}
	isSameTime = now.Equal(tmpT)
	fmt.Println(isSameTime)

	//之前
	isBefe := now.Before(later)
	fmt.Println(isBefe)
	//之后
	isAftr := now.After(later)
	fmt.Println(isAftr)

	//间隔一秒执行
	ticker := time.Tick(time.Second)
	for i := range ticker {
		fmt.Println(i)
	}

}

func showInsertFile() {
	//需要配置器
	fName := "./main.go"
	offset := int64(1)
	whence := 0
	insertStr := " hello ,this is insert string. "
	//需要文件读取器 接受光标
	firstCont, err := readFileWithOffset(fName, offset, whence)
	if err != nil {
		fmt.Println(err)
		return
	}
	//需要新建文件器
	tmpFName := "insert.tmp"
	err = addFile(tmpFName)
	if err != nil {
		fmt.Println(err)
		return
	}
	//把offset位置前的内容存到新的临时文件
	_, err = writeFile(firstCont, tmpFName)
	if err != nil {
		fmt.Println(err)
		return
	}
	//需要光标移动器
	//需要文件写入器
	//把插入内容和写入临时文件
	_, err = appendFile(insertStr, tmpFName)
	if err != nil {
		fmt.Println(err)
		return
	}
	//读取到后面的内容全部追加到临时文件
	fObj, err := os.OpenFile(fName, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	fObj.Seek(offset, whence)
	tmpCont := make([]byte, 0, 100)
	for {
		rdNum, err := fObj.Read(tmpCont)
		if err != nil {
			fmt.Println(err)
			return
		}
		if err == io.EOF {
			if rdNum == 0 {
				break
			}
		}
		_, err = appendFile(string(tmpCont), tmpFName)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	//文件改名器
	_, err = copyFileExe(tmpFName, fName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("成功插入")
}

func appendFile(content string, to string) (size uint64, err error) {
	//os.OpenFile(to, os.O_APPEND, 066)
	return 0, nil
}

func readFileWithOffset(fName string, offset int64, whence int) (string, error) {
	return "", nil
}

func showWriteOffset() {

	fObj, err := os.OpenFile("./main.log", os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fObj.Close()
	//移动光标
	fObj.Seek(1, 0)
	var ret [1]byte

	_, err = fObj.Read(ret[:])
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = fObj.WriteString(string(ret[:]) + " insert ") //覆盖
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(ret[:]))

}

func showCopyFile() {
	fSize, err := copyFileExe("./main.log", "./main_c.log")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("拷贝成功 新文件大小:%v字节 \n", fSize)
}
func copyFileExe(from, to string) (size uint64, err error) {
	//需要目标读取器
	cont, err := readFile(from)
	if err != nil {
		fmt.Println(err)
		return
	}
	//需要新增文件器
	err = addFile(to)
	if err != nil {
		fmt.Println(err)
		return
	}

	//需要文件写入器
	fSize, err := writeFile(cont, to)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("拷贝完成")
	return fSize, nil
}

func writeFile(cont, to string) (uint64, error) {
	//需要文件写入器
	err := ioutil.WriteFile(to, []byte(cont), 0666)
	if err != nil {
		return 0, err
	}
	fInfo, err := os.Stat(to)
	if err != nil {
		fmt.Println(err)
	}
	//fObj, err := os.OpenFile(to, os.O_RDONLY, 0666)
	//if err != nil {
	//	return 0, err
	//}
	//fInfo, err := fObj.Stat()
	//if err != nil {
	//	return 0, err
	//}
	return uint64(fInfo.Size()), nil
}

func addFile(to string) error {
	//需要文件创建器
	fObj, err := os.OpenFile(to, os.O_CREATE, 0666)
	defer fObj.Close()
	if err != nil {
		return err
	}
	return nil
}

func readFile(from string) (content string, err error) {
	//需要文件读取器
	tmp, err := ioutil.ReadFile(from)
	if err != nil {
		return
	}
	return string(tmp), nil
}

func showWriteFileByIoutil() {
	str := "nginx 访问者data:[{}] \n"
	err := ioutil.WriteFile("./main.log", []byte(str), 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("写入成功")
}

func showWriteFileByBufio() {
	//需要文件操作者
	fObj, err := os.OpenFile("./main.log", os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fObj.Close()
	//需要写入字符串到文件
	fieWrr := bufio.NewWriter(fObj)
	str := "nginx 访问者data:[{}] \n"
	for i := 1; i <= 10; i++ {
		fieWrr.WriteString(fmt.Sprint(i, " ", str))
	}
	fieWrr.Flush() //保存到文件中
	fmt.Println("写入成功")
}

func showOsOpenFileType() {
	//创建没有的	os.O_CREATE
	//只写			os.O_WRONLY
	//只读			os.O_RDONLY
	//同时读写		os.O_RDWR
	//追加			os.O_APPEND
	//清空已有的	os.O_TRUNC
}

func showWriteFileByOs() {
	showOsOpenFileType()
	//需要文件资源打开器
	//选择 文件名 打开方式 执行权限
	fObj, err := os.OpenFile("./main.log", os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fObj.Close()
	//需要实际写入日志
	str := "nginx 访问者data:[{}] \n"
	fObj.Write([]byte(str))
	fObj.WriteString(str)
	fmt.Println("保存完成")
}

func showFileReadByBufio() {
	//需要文件选择器
	fileObj, err := os.Open("./main.go")
	if err != nil {
	}
	defer fileObj.Close()
	//需要文件读取器
	bufoRear := bufio.NewReader(fileObj)
	for {
		cont, err := bufoRear.ReadString('\n')
		if err == io.EOF {
			//剩余一部分的文字
			if len(cont) != 0 {
				fmt.Println(cont)
			}
			fmt.Println("文件读取完成")
			return
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(cont)
	}
}

func showFileReadByIOutil() {
	//需要文件选择
	//需要文件读取器
	cont, err := ioutil.ReadFile("./main.go")
	if err != nil {
	}
	fmt.Println(string(cont))
}

func showFileReadByIO() {
	//需要文件选择器
	fileObj, err := os.Open("./main.go")
	if err != nil {
		fmt.Println(err)
		return
	}
	//需要关闭文件资源
	defer fileObj.Close()
	//读取文件全部内容
	var tmp = make([]byte, 128)
	for {
		num, err := fileObj.Read(tmp[:])
		if err == io.EOF {
			fmt.Println("文件读取完成")
			return
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("读取了%v个字节 \n", num)
		fmt.Println(string(tmp))
	}
}

func showInterfaceType() {
	//不确定要保存什么类型,可以用空接口
	m1 := make(map[string]interface{}, 16)
	m1["name"] = "周林"
	m1["age"] = 22
	m1["merried"] = true
	m1["hobby"] = [...]string{"唱", "跳", "rap"}

	printAny(m1)
	printAny(m1["hobby"])

	assign(m1)
	assign(m1["merried"])
	assign(m1["age"])
	assign(m1["name"])
}

func printAny(in interface{}) {
	fmt.Printf("type: %T value: %v \n", in, in)
}

func assign(in interface{}) {
	fmt.Printf("%T \n", in)
	switch in.(type) {
	case string:
		fmt.Println(in.(string))
	case int:
		fmt.Println(in.(int))
	case int64:
		fmt.Println(in.(int64))
	case bool:
		fmt.Println(in.(bool))
	}
}

func assignStr(in interface{}) {
	str, ok := in.(string)
	if !ok {
		fmt.Println("猜错了")
	} else {
		fmt.Println("传进来是个字符串")
	}
	fmt.Println(str)
}

func showInteChin() {
	var a1 animalIC
	ca1 := catIC{
		name: "咪咪",
		feet: 4,
	}

	a1 = ca1
	a1.eat("猫粮")
	fmt.Printf("%T \n", a1)
	fmt.Println(a1)

	kfc := chickenIC{feet: 4}
	a1 = kfc
	a1.eat("鸡饲料")
	fmt.Printf("%T \n", a1)
	fmt.Println(a1)

	//指针和值都能接收
	var aP1 animalIC
	aP1 = &kfc
	fmt.Printf("%T \n", aP1)
	aP1 = kfc
	fmt.Printf("%T \n", aP1)
	fmt.Println(aP1)
}

type animalIC interface {
	move()
	eat(string)
}

type catIC struct {
	name string
	feet int8
}

func (c catIC) move() {
	fmt.Println("猫走路")
}
func (c catIC) eat(in string) {
	fmt.Println("猫咪吃东西", in)
}

type chickenIC struct {
	feet int8
}

func (c chickenIC) move() {
	fmt.Println("鸡走路")
}
func (c chickenIC) eat(in string) {
	fmt.Println("鸡吃东西", in)
}

type catI struct{}
type dogI struct{}

//定义接口
type animalI interface {
	speak(string) string //函数名(参数),返回值
}

func (this dogI) speak(string) (ret string) {
	fmt.Println("汪汪汪")
	return
}

func (this catI) speak(string) (ret string) {
	fmt.Println("喵喵喵")
	return ""
}

func do(this animalI) {
	fmt.Println(this.speak(""))
}

func showInterfaceDog() {
	var c1 = catI{}
	var d1 = dogI{}

	do(c1)
	do(d1)
}

type studentO struct {
	id   uint64
	name string
}

type studentOMgr struct {
	studentOSave map[uint64]studentO
}

func showStutAdmnByObj() {
	smgr := studentOMgr{studentOSave: make(map[uint64]studentO, 10)}
	for {
		showStutMenu()
		fmt.Println("请输入序号:")
		var in uint64
		fmt.Scanln(&in)
		fmt.Printf("你选择的序号是:%v \n", in)

		//需要分配器
		switch in {
		//需要执行器
		case 1:
			smgr.getStut()
		case 2:
			smgr.addStut()
		case 3:
			smgr.editStut()
		case 4:
			smgr.delStut()
		case 5:
			fmt.Println("谢谢使用!")
			os.Exit(0)
		}
	}
}

func showStutMenu() {
	//需要
	fmt.Println("欢迎光临学生管理系统!")
	fmt.Println(`
		1.查看所有学生
		2.新增学生
		3.修改学生
		4.删除学生
		5.退出
	`)
}

func (this *studentOMgr) addStut() {
	//需要输入者
	var (
		id   uint64
		name string
	)
	fmt.Println("请输入学号:")
	fmt.Scanln(&id)
	fmt.Println("请输入姓名:")
	fmt.Scanln(&name)
	//需要初始化学生
	newStut := new(studentO)
	newStut.id = id
	newStut.name = name
	//需要存储学生
	this.studentOSave[newStut.id] = *newStut
}

func (this studentOMgr) editStut() {
	//需要输入器
	var (
		id    uint64
		name  string
		field string
	)
	fmt.Println("请输入你需要修改的学生学号:")
	fmt.Scanln(&id)
	fmt.Printf("你输入的学号是:%v\n", id)
	//需要读取器
	editStu, ok := this.studentOSave[id]
	if !ok {
		fmt.Println("输入的学号不存在,请重新输入")
		return
	}
	fmt.Println("请输入需要修改的字段 name/id")
	fmt.Scanln(&field)
	fmt.Printf("请输入字段%v 的新内容:  \n", field)

	switch field {
	//需要存储器
	case "name":
		fmt.Scanln(&name)
		editStu.name = name
		this.studentOSave[id] = editStu
	case "id":
		delete(this.studentOSave, editStu.id)
		editStu.id = id
		this.studentOSave[id] = editStu
	default:
		fmt.Println("请检查输入内容是否正确")
	}

}
func (this *studentOMgr) delStut() {
	//需要输入处理者
	var (
		id uint64
	)
	fmt.Println("请输入需要删除的序号:")
	fmt.Scanln(&id)
	//需要删除存储器
	delete(this.studentOSave, id)
}
func (this *studentOMgr) getStut() {
	//访问存储器
	for _, v := range this.studentOSave {
		fmt.Printf("学生学号:%v 姓名:%v \n", v.id, v.name)
	}
}
func (this *studentOMgr) newStut(id uint64, name string) studentO {
	//实例化学生并返回
	return studentO{
		id:   id,
		name: name,
	}
}

type personJ struct {
	Name string `json:"name" db :"name" ini:"name"`
	Age  int    `json:"age"`
}

type personJn struct {
	name string `name:"name" db:"name" ini:"name"`
	Age  int    `json:"age"`
}

func showJson() {
	//对象转为json输出
	//序列化
	p1 := personJ{Name: "周林", Age: 9000}
	b, _ := json.Marshal(p1)
	fmt.Printf("%v \n", string(b))

	//反序列化
	//json转为对象
	str := "{\"name\":\"周林\",\"age\":9000}"
	var p2 personJ
	json.Unmarshal([]byte(str), &p2)
	fmt.Println(p2)

	//小写字段是外部看不到的
	p3 := personJn{
		Age:  22,
		name: "周林",
	}
	b, _ = json.Marshal(p3)
	fmt.Printf("%v \n", string(b))
}

type address struct {
	city string
}

type company struct {
	name     string
	province string
	//嵌套的结构体
	address
}

type animal struct {
}

type cat struct {
	name string
	animal
}

func (this *cat) miao() {
	fmt.Println("喵喵喵~")
}
func (this *animal) run() {
	fmt.Println("在动")
}

func showNoNameConst() {
	comy := company{name: "apple", province: "", address: address{city: "加利福尼亚"}}
	fmt.Println(comy.address.city)
	fmt.Println(comy.city)
}

type student struct {
	id   int64
	name string
}

func showStutAdmnByFunc() {
	for {
		showStutAdmn()
	}

}

func showStutAdmn() {
	//打印菜单
	fmt.Println("欢迎光临学生管理系统!")
	fmt.Println(`
		1.查看所有学生
		2.新增学生
		3.删除学生
		4.退出
	`)
	//需要接收器
	fmt.Println("请输入你要干啥:")
	//等待用户要做什么
	inpt := 0
	fmt.Scanln(&inpt)
	fmt.Printf("你选择了第%v个选项 \n", inpt)
	//需要转发者
	//需要处理者
	switch inpt {
	case 1:
		showAllStut()
	case 2:
		addStut()
	case 3:
		delStut()
	case 4:
		fmt.Println("谢谢使用!")
		os.Exit(0)
	default:
		fmt.Println("请正确输入")
	}
}

var stuSave = make(map[int64]*student, 10)

func showAllStut() {
	//需要读取存储
	for k, v := range stuSave {
		fmt.Printf("学号:%v 姓名:%v \n", k, v.name)
	}
}

func newStut(id int64, name string) *student {
	return &student{id: id, name: name}
}

func addStut() {
	//需要输入者
	var (
		id   int64
		name string
	)
	fmt.Print("请输入学生学号:")
	fmt.Scanln(&id)
	fmt.Print("请输入学生姓名:")
	fmt.Scanln(&name)
	//构造函数
	newStu := newStut(id, name)
	//需要存储器
	stuSave[id] = newStu
}

func delStut() {
	//需要输入者
	fmt.Println("请输入想删除的学生学号:")
	var id int64
	fmt.Scanln(&id)
	//操作存储者
	delete(stuSave, id)
	fmt.Printf("你删除了学号为%v的学生信息", id)
}

func showReceiveVal() {
	d1 := new(dog)
	d1.name = "周林"
	d1.wang()

	yuanshuai := new(person)
	yuanshuai.name = "ys"
	yuanshuai.age = 79
	makeUp2(yuanshuai)
	fmt.Printf("%v \n", yuanshuai)

	//自定义类型
	mi1 := new(myInt)
	mi1.hello()

	p2 := person{
		name: "周冠花",
		age:  15,
	}
	fmt.Println(p2)
}

type myInt int

func (mi myInt) hello() {
	fmt.Println("我是一个int")
}

//可以在外面访问
type Dog struct {
	name string
}

func (d dog) wang() {
	fmt.Println("汪汪汪~", d.name)
}

type person struct {
	name   string
	age    int
	gender string
	hobby  []string
}

func showContruct() {
	type myInt int
	type yourInt = int

	var n myInt
	n = 10
	fmt.Println(n)
	fmt.Printf("%T %v \n", n, n)

	var n2 yourInt
	n2 = 100
	fmt.Println(n2)
	fmt.Printf("%T %v \n", n2, n2)

	var c rune
	c = 20013 + 1
	fmt.Printf("%T %v %c \n", c, c, c)

	for i := uint32(0); i < 5; i++ {
		num := 20013 + i
		c = rune(num)
		fmt.Printf("%T %v %c \n", c, c, c)
	}

	var p person
	p.name = "周林"
	p.age = 100
	p.gender = "男"
	p.hobby = []string{"篮球", "足球", "双色球"}

	var p2 person
	p2.name = "张三"
	p2.age = 10
	fmt.Printf("name %v gender %v\n", p2.name, p2.gender)

	var s struct {
		name string
		age  int
	}

	s.name = "匿名结构体类型"
	s.age = 10
	fmt.Printf("%v %T \n", s, s)

	p3 := person{}
	p3.name = "李四"
	p3.age = 17
	makeUp(p3)
	fmt.Println(p3)
	makeUp2(&p3)
	fmt.Println(p3)

	pointPern := new(person)
	fmt.Printf("%T %v \n", pointPern, pointPern)

	var a int
	a = 100
	b := &a
	fmt.Printf("%T %v \n", a, a)
	fmt.Printf("%T %v \n", b, b)
	fmt.Printf("%p \n", &a)
	fmt.Printf("%p %v\n", b, b)

	type x struct {
		a int8
		b int8
		c int8
	}

	m := x{
		a: int8(10),
		b: int8(20),
		c: int8(30),
	}

	fmt.Printf("%p \n", &(m.a))
	fmt.Printf("%p \n", &(m.b))
	fmt.Printf("%p \n", &(m.c))

}

//对象+构造函数
type dog struct {
	name string
}

func newDog(name string) *dog {
	return &dog{name: name}
}

//构造函数 newObj 根据内容构造一个对象
//拷贝的值占用空间比较大的,不如直接引用对象的地址
func newPern(name string, age int) *person {
	return &person{
		name: name,
		age:  age,
	}
}

func makeUp(in person) {
	in.age += 1
}

func makeUp2(in *person) {
	//(*in).age += 1
	in.age += 1
}

func showRecursive() {
	rece1()
}

func rece1() {
	n := 5
	fmt.Printf("%v的阶乘结果是%v", n, factorial(n))
}

func factorial(in int) int {
	//终结的条件 达到底层的return
	if in > 0 {
		result := in * factorial(in-1)
		return result
	}
	return 1
}

func showDispatchCoin() {
	//准备条件
	users := []string{
		"Matthew", "Sarah", "Augustus", "Heidi", "Emilie", "Peter", "Giana", "Adriano", "Aaron", "Elizabeth",
	}
	distribution := make(map[string]int, len(users))
	left := dispatchCoin(50, users, distribution)
	//返回结果
	fmt.Println("剩余金币:", left)
	fmt.Println("每人分得金币: ", distribution)
}

//分配金币
func dispatchCoin(coinsNum int, users []string, ret map[string]int) (left int) {
	//需要执行分配的操作者
	for _, name := range users {
		for _, v := range name {
			//需要分配器
			dispNum := dispSet(v)
			ret[name] += dispNum
			coinsNum -= dispNum
		}
	}
	return coinsNum
}

func dispSet(in rune) (addNum int) {
	switch in {
	case 'e', 'E':
		return 1
	case 'i', 'I':
		return 2
	case 'o', 'O':
		return 3
	case 'u', 'U':
		return 4
	}
	return 0
}

func showFmt() {
	fmt.Print("沙河")
	fmt.Print("娜扎")
	fmt.Println("...")
	fmt.Println("沙河")
	fmt.Println("娜扎")

	//%v 值
	//%T 类型

	//%c 字符
	//%s 字符串
	//%p 指针
	//%f 浮点数
	//%t 布尔值

	//数值类型
	//%b binary二进制
	//%o 八进制
	//%x %X 十六进制

	//浮点数

	//获取屏幕输入
	var s, d, f string
	//fmt.Scan(&s,&d,&f)
	//fmt.Println(s,d,f)

	fmt.Scanln(&s, &d, &f)
	fmt.Println(s, d, f)
}

//panic recover
func showRecover() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
			fmt.Println("释放数据库连接...")
		}
	}()
	fmt.Println("a")
	panic("权限检查失败") //程序崩溃退出
	fmt.Println("b")
}

//传一个指针到匿名函数中
func df6() (x int) {
	//1,返回值赋值	x=5
	//2.defer	x = 6
	//3.RET返回 return x
	defer func(x *int) {
		(*x)++
	}(&x)
	return 5
}

func showDefer2() {

}

func showDefer() {
	//1.返回值赋值
	//2.defer
	//3.RET跳转到调用处
	fmt.Println(df1())
	fmt.Println(df2())
	fmt.Println(df3())
	fmt.Println(df4())
	fmt.Println(df6())
}

func showDeferRegister() {
	a := 1
	b := 2
	defer calc("1", a, calc("10", a, b))
	a = 0
	defer calc("2", a, calc("20", a, b))
	b = 1
}

func calc(index string, a, b int) int {
	ret := a + b
	fmt.Println(index, a, b, ret)
	return ret
}

//执行顺序
//a = 1
//b = 2
//遇到defer 进行登记  calc(1,1 calc(10,1,2))
//发现参数是函数的返回值 先执行函数求出固定值 再登记defer
//执行 calc(10,1,2)		打印 10 1 2 3
//defer 登记 calc(1,1,3)
//a = 0
//遇到defer 进行登记 calc(2,a,calc(20,0,2))
//发现参数是函数的返回值,先执行函数来得到固定参数,再登记defer
//执行 calc(20,0,2) 打印 20 0 2 2
//defer 登记 calc(2,0,2)
//b = 1
//函数结尾开始return
//从defer登记取出最近一条并执行 calc(2,0,2)
//执行 calc(2,0,2) 打印 2 0 2 2
//从defer登记取出最近第二条并执行 calc(1,1,3)
//执行 calc(1,1,3) 打印 1 1 3 4
//defer取完结束
//函数ret跳转

//10 1 2 3
//20 0 2 2
//2 0 2 2
//1 1 3 4

//从上到下阅读执行代码
//遇到defer 把函数参数固定,登记下来(如果参数是函数的返回值,先把函数执行拿到返回值)
//下个defer 压栈,遵循先进后出
//函数的return开始 1.赋值返回值 2.defer 3.ret跳转

func df4() (x int) {
	defer func(x int) {
		x++
	}(x)
	return 5

	//1.返回值赋值
	//x = 5
	//2.defer
	//fun x 传进来
	//3.RET跳转到调用处
	//x = 5
}

func df3() (y int) {
	//1.返回值赋值
	//y = x = 5
	//2.defer
	//x = 6
	//3.RET跳转到调用处
	//x=5

	x := 5
	defer func() {
		x++
	}()
	return x
}

func df2() (x int) {
	//1.返回值赋值
	//x = 5
	//2.defer
	//x ++
	//3.RET跳转到调用处
	//return x (x=6)
	defer func() {
		x++
	}()
	return 5
}

func df1() int {
	x := 5
	defer func() {
		x++ //返回的是x不是返回值
	}()
	return x
}

func showDeferExec() {
	fmt.Println("start")
	defer fmt.Println("嘿嘿嘿")
	fmt.Println("end")
}

//闭包函数
//是一种函数,包含了外部作用域的一个变量
func showCliseFunc() {
	ret := adder(100)
	fmt.Println(ret)
	ret2 := ret(200)
	fmt.Println(ret2)

	//目标能够 cf1(cf2)
	ret3 := cf3(cf2, 100, 200)
	cf1(ret3)
}

func cf1(f func()) {
	fmt.Println("this is f1")
}

func cf2(x, y int) {
	fmt.Println("this is f2")
	fmt.Println(x + y)
}

//cf1(cf3)
func cf3(f func(int, int), m, n int) func() {
	tmp := func() {
		f(m, n)
	}
	return tmp
}

func adder(x int) func(int) int {
	return func(y int) int {
		x += y
		return x
	}
}

//底层原理:
//函数可以作为返回值
//函数内部找不到变量,可以层层往外找,直到全局
func showSendFunc() {
	lixiang(add, 100, 200)
}

func lixiang(func1 func(uint64, uint64) uint64, m, n uint64) {
	tmp := func() uint64 {
		return func1(m, n)
	}
	ret := tmp()
	fmt.Println(ret)
}

//匿名函数
var f1 = func(x, y int) {
	fmt.Println(x + y)
}

func showFuncAsRet() {
	funcName := chooseFunc("add")
	fmt.Printf("%T %v", funcName, funcName)
}

func chooseFunc(funcName string) (returnFunc func(uint64, uint64) uint64) {
	if funcName == "add" {
		return add
	}
	return nil
}

func add(x, y uint64) uint64 {
	return x + y
}

func loadFunc(x, y uint64, op func(uint64, uint64) uint64) uint64 {
	return op(x, y)
}

func showFuncAsPam() {
	//函数作为参数传入
	ret := loadFunc(1, 2, add)
	fmt.Println(ret)
}

func showFunc3() {
	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}
}

var envX = 100

func showFunc2() {
	//函数内查找变量的顺序
	//1.内部找
	//2,往外找,找到全局变量里
	//都没有就报错
	fmt.Println(envX)
}

func showNoNameFunc() {
	func(msg string) {
		fmt.Println("no name func get param :", msg)
	}("ok")
}

func showStrNum() {
	//判断字符串中汉字的数量
	s1 := "my name is 张三"
	fmt.Printf(" \" %v \"中汉字数量是: %v \n", s1, getHanNum(s1))

	//字符串中出现次数
	s2 := "how do you do "
	fmt.Printf(" %v中各个单词出现的次数为: %v \n", s2, getStrInNum(s2))

	//检查是否是回文
	s3 := "上海海上"
	//堆栈 压进去 先进后出
	//除以一半 从中间开始 s[0] s[len-1]
	fmt.Printf(" %v是否是回文:  %v \n", s3, checkPalindrome(s3))

}

func checkPalindrome(in string) bool {
	//需要文字格式转换器
	r := str2Rune(in)
	//需要检查器
	for i := 0; i < len(r)/2; i++ {
		if r[i] != r[len(r)-1-i] {
			return false
		}
	}
	return true
}

func str2Rune(in string) []rune {
	//转成rune slice
	r := make([]rune, 0, len(in))
	for _, v := range in {
		r = append(r, v)
	}
	return r
}

func getHanNum(in string) int {
	//遍历每个字符
	strNum := 0
	for _, v := range in {
		//检查是语言,计数器+1
		if isLang(v, "chinese") {
			strNum += 1
		}
	}
	return strNum
}

func getStrInNum(in string) map[string]int {
	s3 := strings.Split(in, " ")
	m1 := make(map[string]int, 10)
	for _, v := range s3 {
		if v == "" {
			continue
		}
		if _, ok := m1[v]; ok {
			m1[v] += 1
		} else {
			m1[v] = 1
		}
	}
	return m1
}

func isLang(s rune, lang string) bool {
	if lang == "chinese" || lang == "china" {
		if unicode.Is(unicode.Han, s) {
			return true
		}
	}
	return false
}

func showFunc() {

	f2()
	fmt.Println(f3(1, 1))
	fmt.Println(f4(1, 1))
	fmt.Println(f5())
	fmt.Println(f6(1, 1, false))
}

func f2() {
	fmt.Println("f2")
}

//无命名的返回值
func f3(x, y int) int {
	return x + y
}

//命名返回值,相当于在函数中声明了一个变量
func f4(x, y int) (ret int) {
	ret = x + y
	//命名返回值,可以不写return的变量名
	return
}

func f5() (int, string) {
	return 200, "ok"
}

//连续两个相同类型参数可以简写
func f6(x, y int, z bool) int {
	return x + y
}

//go 不允许使用默认参数
func f7(x, y int) int {
	return x + y
}

func showMap() {
	//var m1 map[string]int	//需要初始化(没有在内存中分配空间)
	//初始化map 估算好map容量,避免运行期间重新开辟内存的开销
	m1 := make(map[string]int, 10)
	fmt.Println(m1)
	m1["理想"] = 9000
	m1["jiwuming"] = 35
	fmt.Println(m1)

	if v, ok := m1["jiwuming"]; ok {
		fmt.Println(v)
	} else {
		fmt.Println("查无此key")
	}

	//遍历
	for k := range m1 {
		fmt.Printf("%v => %v \n", k, m1[k])
	}

	//删除键值对
	delete(m1, "jiwuming")
	delete(m1, "沙河")

	m1["理想"] = 2

	fmt.Println(m1)
	rand.Intn(100)

	//map slice 组合

	s1 := make([]map[int]string, 10, 10)
	//内部的map也需要初始化才能赋值
	s1[0] = make(map[int]string, 1)
	s1[0][10] = "沙河"
	fmt.Println(s1)

	//也需要初始化
	m2 := make(map[string][]int, 10)
	m2["北京"] = []int{10, 11, 12}
	fmt.Println(m2)

}

func showPoint() {
	//取地址 &
	//根据地址取值 *

	n := 18
	//取变量使用的内存地址
	fmt.Println(&n) //0xc0000ac058
	p := &n
	fmt.Printf("%T \n", p) //*int

	//取内存地址中保存的值
	m := *p
	fmt.Println(m)

	//面试题 什么结果?
	//找不到值 赋值时候报错,因为找不到内存地址
	var a *int
	fmt.Println(a)
	//*a = 100
	//fmt.Println(a)

	//分配内存地址,并存储值
	var a2 = new(int)
	*a2 = 100
	fmt.Println(*a2)

}

func showCopy() {
	a1 := []int{1, 3, 5}
	a2 := a1
	//var a3 []int	//空数组 无长度 无容量
	a3 := make([]int, 3, 5) //长度3 容量5
	copy(a3, a1)
	fmt.Println(a1, a2, a3)
	a1[0] = 100
	fmt.Println(a1, a2, a3)

	a1 = append(a1[:1], a1[2:]...)
	fmt.Println(a1, a2, a3)

	x1 := [...]int{1, 3, 5} //数组
	s1 := x1[:]             //切片
	fmt.Println(s1, len(s1))

	s1 = append(s1[:1], s1[:2]...) //修改切片就是修改底层数组
	fmt.Println(s1, len(s1))
	fmt.Println(x1, len(x1))

	//面试题
	s2 := make([]int, 5, 10)
	for i := 0; i < 10; i++ {
		s2 = append(s2, i)
	}
	fmt.Printf("s2 value: %v cap: %v \n", s2, cap(s2))

	//
	a4 := [...]int{1, 3, 5, 7, 9, 11, 13, 15, 17}
	s4 := a4[:]

	s4 = append(s4[0:1], s4[2:]...)
	fmt.Println(s4)
	fmt.Println(a4)

}

func showAppend() {
	//append 为切片增加元素
	s6 := []string{"北京", "上海", "广州"}
	//s6[3] = "无锡"	//会索引越界,因为用了不存在的索引
	fmt.Printf("slice :%v len:%v cap:%v \n", s6, len(s6), cap(s6))

	s6 = append(s6, "深圳") //容量从3变成6,
	// 因为切片使用的底层数组放不下的时候,go底层换了一个更大容量的数组
	fmt.Printf("slice :%v len:%v cap:%v \n", s6, len(s6), cap(s6))

	s7 := []string{"武汉", "西安", "苏州"}
	s6 = append(s6, s7...)
	fmt.Printf("slice :%v len:%v cap:%v \n", s6, len(s6), cap(s6))
}

func showSlice() {
	//切片定义
	var a = []int{1, 2, 3}
	fmt.Println(a)

	//数组得到切片,开始包含结束不包含
	arr1 := [5]int{0, 1, 2, 3, 4}
	slice1 := arr1[1:3]
	fmt.Printf("你的数组是%v 切片是%v 长度是%v 容量是%v \n",
		arr1, slice1, len(slice1), cap(slice1))

	//切片的切片,上限是切片的容量cap,
	slice2 := slice1[:2]
	fmt.Printf("你的数组是%v 切片是%v 长度是%v 容量是%v \n",
		arr1, slice2, len(slice2), cap(slice2))

	//切片要用len==0判断是否为空,不能用nil去判断
	arr2 := [2]int{1, 2}
	var emySle = arr2[0:0]
	fmt.Printf("emySle==nil  %v \n", emySle == nil)
	fmt.Printf("len(emySle) == 0 %v \n", len(emySle) == 0)

	//拷贝的切片底层是同一个引用数组,修改会互相影响
	var slice3 = []int{1, 2, 3}
	slice4 := slice3
	slice4[0] = 10
	fmt.Printf("slice3 = %v slice4 = %v \n", slice3, slice4)

	//make函数创造切片
	s5 := make([]int, 0, 5) //已有长度0 分配容量5
	fmt.Printf("slice :%v len:%v cap:%v \n", s5, len(s5), cap(s5))

}

func showArray() {
	var ar1 [3]int //长度3 容量3 值是 [0 0 0]
	var ar2 [10]int
	fmt.Println(len(ar1), cap(ar1), ar1)
	fmt.Println(len(ar2), cap(ar2), ar2)

	//定义数组
	var type1 = [3]int{1, 2, 3}
	var type2 = [...]int{1, 2, 3}
	var type3 = [...]int{0: 1, 1: 2, 2: 3}
	fmt.Println(type1, type2, type3)

	//多维数组
	numArr := [3][2]string{
		{"腾讯", "阿里"},
		{"华为", "小米"},
	}
	fmt.Println(numArr)
	//...推导长度,外层能用
	numArr2 := [...][2]float64{
		{116.420699, 39.97562},
	}
	fmt.Println(numArr2)
	//不支持内层使用...
	//numArr3 := [2][...]float64{
	//	{116.420699,39.97562},
	//}
	//fmt.Println(numArr3)

	//数组是值类型

}

func multiTable() {
	for i := 1; i <= 9; i++ {
	forloop1:
		for j := 1; j <= 9; j++ {
			if j > i {
				continue forloop1
			}
			ret := fmt.Sprintf("%v x %v = %v ", i, j, i*j)
			fmt.Print(ret)
		}
		fmt.Println()
	}
}

func showBreak() {
	i := 0
BREAKDEMO1:
	for {
		if i == 10 {
			break BREAKDEMO1
		}
		fmt.Println("value is ", i)
		i++
	}
}

func showGoto() {
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if j == 2 {
				goto breakTag
			}
		}
	}
breakTag:
	fmt.Println("结束for循环")
}

func showSwitch() {
	list := []int{0, 1, 3, 4, 5}
	for _, val := range list {
		ret := ""
		switch val {
		case 1, 3, 5, 7:
			ret = "奇数"
		case 2, 4, 6, 8:
			ret = "偶数"
			fallthrough
		default:
			ret = fmt.Sprintf("当前值为(%v)", val)
		}
		fmt.Println("结果是: ", ret)
	}
}
func showFor() {
	i := 0
	for ; i < 10; i++ {
		fmt.Println(i)
	}
	for i <= 10 {
		fmt.Println(i)
		i++
	}
}

func showIf() {
	if time := "2022"; time == "" {
		fmt.Println("time null")
	} else if time != "" {
		fmt.Println("time not null")
	} else {
		fmt.Println("no time")
	}
}

func showStr() {
	c1 := 'h'
	c2 := 'l'
	//单独一个字符用''
	c3 := '沙'
	fmt.Println("c1 :", c1, " c2 :", c2, "c3 :", c3)
	fmt.Printf("type: %T value:%v", c3, c3)
	fmt.Println()

	//拼接字符串
	s1 := "你的余额是:"
	s2 := "100元"
	str := fmt.Sprintf("%s%s", s1, s2)
	fmt.Println("str is :", str)

	//分隔
	res := strings.Split(str, ":")
	fmt.Println("res is :", res)

	//包含
	fmt.Println(strings.Contains(str, "余额"))

	//前后缀
	fmt.Println(strings.HasPrefix(str, "你的"))
	fmt.Println(strings.HasSuffix(str, "元"))

	//出现的位置
	fmt.Println(strings.Index(str, "余额"))
	fmt.Println(strings.LastIndex(str, "元"))

	//插入间隔字符
	fmt.Println(strings.Join([]string{"收到消息:", str, "over"}, "--"))

	//修改字符串
	tmp := []rune(str)
	tmp[len(tmp)-1] = '¥'
	fmt.Println(string(tmp))

	//转为utf-8字节
	str = "u body weight : 100kg"
	tmp2 := []byte(str)
	tmp2[0] = byte('i')
	fmt.Println(string(tmp2))

	//统计中文难
	str = "hello沙河小王子"
	tmp = []rune(str)
	fmt.Println("长度为:", len(tmp))

}

func showBool() {
	var bl bool //false
	fmt.Printf("type: %T, value : %v", bl, bl)
}

func showFloat() {
	fmt.Println(math.MaxFloat32)
	fmt.Println(math.MaxFloat32)

	f1 := 1.2456
	fmt.Printf("%T \n", f1) //变量在go中的值类型	float64
	f2 := float32(f1)
	fmt.Printf("%T", f2) //float32
}

func showFmtNum() {
	var a int = 10
	fmt.Printf("%b", a) //base 2 转成二进制
	fmt.Println("")

	fmt.Printf("%o", a) //base 8 八进制
	fmt.Println("")

	fmt.Printf("%x", a) //base 16 十六机制
	fmt.Println("")

	fmt.Printf("%d", a) //base 10 十进制
	fmt.Println("")

	//定义一个八进制的数字
	var b int = 077
	fmt.Printf("%d", b)
	fmt.Println("")
	fmt.Printf("%b", b)
	fmt.Println("")
	fmt.Printf("%x", b)
	fmt.Println("")
	fmt.Printf("%o", b)
	fmt.Println("")

	//定义一个十六进制的数字
	var c int = 0xff
	fmt.Printf("%d", c)
	fmt.Println("")
	fmt.Printf("%b", c)
	fmt.Println("")
	fmt.Printf("%x", c)
	fmt.Println("")
}

func showIota() {
	fmt.Println("n1", n1)
	fmt.Println("n2", n2)
	fmt.Println("n3", n3)
	fmt.Println("n3", n4)
	fmt.Println("n3", n5)
}
