package zlock_check

import (
	"github.com/Xuzan9396/zlog"
	"testing"
	"time"
)

func Test_lock(t *testing.T)  {
	zlog.SetEnv(zlog.LOG_DEBUG)
	lockChek := GetLockCheck()
	index := lockChek.AddFunc("testv1")
	defer lockChek.DelFunc(index)

	go func() {
		for i := range lockChek.GetLockChan(){
			zlog.F("lock").Error(i.Name , " ",i.Time ,"s")
		}
	}()
	for{
		time.Sleep(1*time.Second)
	}
}
