package zlock_check

import (
	"github.com/Xuzan9396/zlog"
	"testing"
	"time"
)

func Test_lock(t *testing.T)  {
	zlog.SetEnv(zlog.LOG_DEBUG)
	InitLockCheck(5,5*time.Second)
	index := AddLockFunc("testv1")
	defer DelLockFunc(index)

	go func() {
		for i := range GetLockChan(){
			zlog.F("lock").Error(i.Name , " ",i.Time ,"s")
		}
	}()
	for{
		time.Sleep(1*time.Second)
	}
}
