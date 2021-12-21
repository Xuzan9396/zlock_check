# zlock_check
## 1.死锁检测
```go
package main

import (
	"github.com/Xuzan9396/zlock_check"
	"log"
	"time"
)

func main()  {
	// 超过5s检测, 5s检测一次
	zlock_check.InitLockCheck(5,5*time.Second) // 不执行则默认 取  超过60s,1分钟执行一次
	defer zlock_check.DelLockFunc(zlock_check.AddLockFunc("test"))
	go func() {
		for i := range zlock_check.GetLockChan(){
			log.Println("锁住的函数",i.Name,i.Time)
		}
	}()
	select {

	}
}

}
```

## 2.检测运行时间
```go
package main

import (
	"log"
	"sync"
	"time"
	"github.com/Xuzan9396/zlock_check"

)


func main()  {
	var sg sync.WaitGroup
	go func() {
		for i := range zlock_check.GetTimeChan() {
			log.Println("超时检测函数:",i.FuncName,i.DiffTime,"ms")
		}
	}()
	sg.Add(2)
	go func() {
		defer  sg.Done()
		check()

	}()
	go func() {
		defer sg.Done()
		checkv2()

	}()
	sg.Wait()
	time.Sleep(2*time.Second)

}

func check(res ...int64)  {
	// 最后一位参数100为毫秒
	defer zlock_check.TimeEnd(zlock_check.TimeStart(),"check",100)
	time.Sleep(2*time.Second)
	log.Println(res == nil )
}


func checkv2(res ...int64)  {
	defer zlock_check.TimeEnd(zlock_check.TimeStart(),"checkv2",100)
	time.Sleep(3*time.Second)
	log.Println(res == nil )
}
```