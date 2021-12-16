package zlock_check

import (
	"github.com/Xuzan9396/zlog"
	"runtime/debug"
	"sync"
	"time"
)

var lockIndex uint64  = 0
var g_pLockCheck *LockCheck
var onces sync.Once
type funcElem struct {
	name string
	visitTime int64

}
type LockChan struct {
	Name string
	Time int64
}
type LockCheck struct {
	sync.RWMutex //数据锁
	checkTime int64 // 大于这个时间判断锁失败
	checkTimer time.Duration
	dataFunc map[uint64]*funcElem
	chans chan *LockChan
}

/**
第一个参数：执行参数大于当前多少s打印 默认60s
第二个参数: 多长时间检测一次 默认 1分钟
 */
func GetLockCheck(ts ...interface{}) *LockCheck{
	onces.Do(func() {

		g_pLockCheck = &LockCheck{
			dataFunc:make(map[uint64]*funcElem,0),
			chans: make(chan *LockChan,20),
		}
		g_pLockCheck.setTickMin(ts) // 设置检测时间
		go g_pLockCheck.LockRun()
	})

	return g_pLockCheck
}

func (c *LockCheck)setTickMin(ts ...interface{})  {
	var t int64
	var tickMin time.Duration
	if ts != nil && len(ts) > 0  {
		switch s := ts[0].(type) {
		case int:
			t = int64(s)
		case int64:
			t = s
		}
		if len(ts) == 1{
			tickMin = time.Minute * 1
		}else{
			tickMin = ts[1].(time.Duration)
		}
	}else{
		t = 60
		tickMin = time.Minute * 1
	}
	c.checkTime = t
	c.checkTimer = tickMin
}

func (c *LockCheck)GetLockChan() chan *LockChan  {
	return c.chans
}
func (c *LockCheck)setLockChan(name string ,t int64 )  {
	select {
	case c.chans <- &LockChan{Name: name,Time: t}:
	default:

	}
}

func (l *LockCheck)AddFunc(funcName string)uint64{
	elem:=&funcElem{
		name:funcName,
		visitTime:time.Now().Unix(),
	}
	l.Lock()
	defer l.Unlock()
	lockIndex++
	l.dataFunc[lockIndex]=elem
	return lockIndex
}

func (l *LockCheck)DelFunc(index uint64){
	l.Lock()
	defer l.Unlock()
	delete(l.dataFunc,index)
}

func (c *LockCheck)LockRun(){
	defer func(){
		if err := recover(); err != nil {
			zlog.F("lock").Error(string(debug.Stack()))
		}
	}()
	tick := time.Tick(c.checkTimer) //
	for {
		select {
		case <-tick:
			//逻辑处理
			c.Print()

		}
	}
}

func (l* LockCheck)Print(){
	defer func(){
		if err := recover(); err != nil {
			zlog.F("lock").Error(string(debug.Stack()))
		}
	}()
	arrList:=make([]funcElem,0)
	l.Lock()
	for _,v:=range l.dataFunc{
		arrList = append(arrList,*v)
	}
	l.Unlock()
	l.printThread(arrList)
}

func (c *LockCheck)printThread(arrList []funcElem){
	currTime:=time.Now().Unix()
	for _,v:=range arrList{
		if currTime - v.visitTime > c.checkTime{ //超过60秒
			c.setLockChan(v.name,currTime - v.visitTime)
		}
	}
}
