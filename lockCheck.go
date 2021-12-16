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

func GetLockCheck(ts ...interface{}) *LockCheck{
	onces.Do(func() {
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
		}
		g_pLockCheck = &LockCheck{
			dataFunc:make(map[uint64]*funcElem,0),
			checkTime:t,
			checkTimer: tickMin,
			chans: make(chan *LockChan,20),
		}
		go g_pLockCheck.LockRun()
	})

	return g_pLockCheck
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
	defer l.Unlock()
	for _,v:=range l.dataFunc{
		arrList = append(arrList,*v)
	}
	go l.printThread(arrList)
}

func (c *LockCheck)printThread(arrList []funcElem){
	currTime:=time.Now().Unix()
	for _,v:=range arrList{
		if currTime - v.visitTime > c.checkTime{ //超过60秒
			c.setLockChan(v.name,currTime - v.visitTime)
			//zlog.F("lock").Errorf("func spend over time:%ds,name:%s,",currTime - v.visitTime,v.name)
		}
	}
}
