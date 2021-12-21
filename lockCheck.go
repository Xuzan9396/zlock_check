package zlock_check

import (
	"github.com/Xuzan9396/zlog"
	"log"
	"runtime/debug"
	"sync"
	"time"
)

var lockIndex uint64  = 0
var g_pLockCheck *lockCheck
var onces sync.Once
var checkTime int64 = 60
var checkTimer = 1*time.Minute

type funcElem struct {
	name string
	visitTime int64

}
type LockChan struct {
	Name string
	Time int64
}
type lockCheck struct {
	sync.RWMutex //数据锁
	checkTime int64 // 大于这个时间判断锁失败
	checkTimer time.Duration
	dataFunc map[uint64]*funcElem
	chans chan *LockChan
}
/**
第一个参数：执行参数大于当前多少s打印 默认60s
第二个参数: 多长时间检测一次 默认 1分钟 时间
*/
func InitLockCheck(t int64 ,s time.Duration)  {
	if checkTime <= 0 || checkTimer <= 0*time.Second {
		log.Panic("参数错误")
	}
	checkTime,checkTimer = t,s
}


func getLockCheck() *lockCheck{
	onces.Do(func() {

		g_pLockCheck = &lockCheck{
			dataFunc:make(map[uint64]*funcElem,0),
			chans: make(chan *LockChan,20),
			checkTime: checkTime,
			checkTimer: checkTimer,
		}
		//g_pLockCheck.setTickMin(ts...) // 设置检测时间
		go g_pLockCheck.LockRun()
	})

	return g_pLockCheck
}



func GetLockChan() chan *LockChan  {
	return getLockCheck().chans
}
func (c *lockCheck)setLockChan(name string ,t int64 )  {
	select {
	case c.chans <- &LockChan{Name: name,Time: t}:
	default:

	}
}

func AddLockFunc(funcName string)uint64{
	elem:=&funcElem{
		name:funcName,
		visitTime:time.Now().Unix(),
	}
	getLockCheck().Lock()
	defer getLockCheck().Unlock()
	lockIndex++
	getLockCheck().dataFunc[lockIndex]=elem
	return lockIndex
}

func DelLockFunc(index uint64){
	getLockCheck().Lock()
	defer getLockCheck().Unlock()
	delete(getLockCheck().dataFunc,index)
}

func (c *lockCheck)LockRun(){
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

func (l* lockCheck)Print(){
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

func (c *lockCheck)printThread(arrList []funcElem){
	currTime:=time.Now().Unix()
	for _,v:=range arrList{
		if currTime - v.visitTime > c.checkTime{ //超过60秒
			c.setLockChan(v.name,currTime - v.visitTime)
		}
	}
}
