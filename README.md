# zlock_check
## 1.死锁检测
```go
package main

import (
	"github.com/Xuzan9396/zlock_check"
	"log"
	"runtime"
	"time"
)

func main()  {
	//lockChek := zlock_check.GetLockCheck(5,5*time.Second)
	lockChek := zlock_check.GetLockCheck()
	index := lockChek.AddFunc(runFuncName())
	defer lockChek.DelFunc(index)
	go func() {
		for i := range lockChek.GetLockChan(){
			log.Println("锁住的函数",i.Name,i.Time)
		}
	}()
	for{
		time.Sleep(1*time.Second)
	}
}

func runFuncName()string{
	pc := make([]uintptr,1)
	runtime.Callers(2,pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
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
	model := zlock_check.NewTimeCheck()
	var sg sync.WaitGroup
	go func() {
		for i := range model.GetChan() {
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

}

func check(res ...int64)  {
	start := zlock_check.NewTimeCheck().Start()
	defer zlock_check.NewTimeCheck().End(start,"check",100)
	time.Sleep(2*time.Second)
	log.Println(res == nil )
}



func checkv2(res ...int64)  {
	defer zlock_check.NewTimeCheck().End(zlock_check.NewTimeCheck().Start(),"checkv2",100)
	time.Sleep(3*time.Second)
	log.Println(res == nil )
}
```