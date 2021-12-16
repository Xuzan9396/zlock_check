
package logic

import (
	"sync"
	"time"
	"github.com/Xuzan9396/zlog"
	"runtime/debug"
)

var lockIndex uint64  = 0
var g_pLockCheck *LockCheck
type funcElem struct {
	name string
	visitTime int64

}
type LockCheck struct {
	sync.RWMutex //数据锁
	checkTime int64 // 超过多长时间检测
	dataFunc map[uint64]*funcElem
}

func GetLockCheck() *LockCheck{
	if g_pLockCheck == nil {
		g_pLockCheck = &LockCheck{
			dataFunc:make(map[uint64]*funcElem,0),
			checkTime:60,

		}
		go LockRun()
	}

	return g_pLockCheck
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

func LockRun(){
	defer func(){
		if err := recover(); err != nil {
			zlog.F().Error(string(debug.Stack()))
		}
	}()
	tick := time.Tick(1 * time.Minute) //1分钟
	for {
		select {
		case <-tick:
			//逻辑处理
			GetLockCheck().Print()

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
			zlog.F("lock").Errorf("func spend over time:%ds,name:%s,",currTime - v.visitTime,v.name)
		}
	}
}
