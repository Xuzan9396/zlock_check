package zlock_check

import (
	"github.com/Xuzan9396/zlog"
	"testing"
	"time"
)

func Test_lock(t *testing.T)  {
	zlog.SetEnv(zlog.LOG_DEBUG)
	index := GetLockCheck().AddFunc("testv1")
	defer GetLockCheck().DelFunc(index)
	for{
		time.Sleep(1*time.Second)
	}
}
