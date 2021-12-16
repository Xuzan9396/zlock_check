package zlock_check

import (
	"log"
	"sync"
	"testing"
	"time"
)

func Test_timeCheck(t *testing.T)  {
	model := NewTimeCheck()
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
	start := NewTimeCheck().Start()
	defer NewTimeCheck().End(start,"check",100)
	time.Sleep(2*time.Second)
	log.Println(res == nil )
}



func checkv2(res ...int64)  {
	defer NewTimeCheck().End(NewTimeCheck().Start(),"check",100)
	time.Sleep(3*time.Second)
	log.Println(res == nil )
}