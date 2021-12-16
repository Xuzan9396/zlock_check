package zlock_check

import (
	"sync"
	"time"
)

var timeCheck *TimeCheck
var timeOnce sync.Once

type TimeCheck struct {
	chans chan *RespCheck
}

type RespCheck struct {
	FuncName string
	DiffTime  int64
}

func NewTimeCheck() *TimeCheck {
	timeOnce.Do(func() {
		timeCheck = &TimeCheck{
			chans: make(chan *RespCheck,1000),
		}
	})
	return timeCheck
}

func (c *TimeCheck)GetChan() chan *RespCheck{
	return c.chans
}

func (c *TimeCheck)setChan(diff int64,funcName string )  {
	select {
	case c.chans <- &RespCheck{
		FuncName: funcName,
		DiffTime: diff,
	}:
	default:


	}
}

func (c *TimeCheck)Start() int64 {
	return  time.Now().UnixNano() / 1000000
}

// 毫秒超过多少毫秒打印
func (c *TimeCheck)End(startTime int64,funcName string ,overs ...int64)  {
	lastTime := time.Now().UnixNano() / 1000000
	diff := lastTime- startTime
	var over int64
	if overs != nil {
		over = overs[0]
	}else{
		over = 100
	}
	if diff > over{
		c.setChan(diff,funcName)
	}
}