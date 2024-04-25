package zlock_check

import (
	"github.com/Xuzan9396/zlog"
	"testing"
	"time"
)

func Test_lock(t *testing.T) {
	zlog.SetLog(zlog.LOG_DEBUG)
	InitLockCheck(5, 5*time.Second)
	go func() {
		defer DelLockFunc(AddLockFunc("testv1"))

		select {}

	}()

	go func() {
		defer DelLockFunc(AddLockFunc("testv2"))

		select {}

	}()

	go func() {
		for i := range GetLockChan() {
			zlog.F("lock").Error(i.Name, " ", i.Time, "s", ",id:", i.Id)
		}
	}()

	for {
		time.Sleep(1 * time.Second)
	}
}
