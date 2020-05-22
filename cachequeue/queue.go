package cachequeue

import (
	"sync"
	"time"
)

const DEFAULT_SIZE = 15
type Queue struct{
	Data []interface{}
	Mutex *sync.Mutex

	timeSpy bool  //是否根据时间弃掉数据，不可修改
	ExpireAfter time.Duration
	flag chan int //time queue 停止
	cap int
	timeStep time.Duration
}

type TimeWrapper struct{
	Data interface{}
	CreatAt int64
}

func NewEmpty() *Queue{
	return &Queue{Data: make([]interface{}, 0, DEFAULT_SIZE), Mutex: &sync.Mutex{},}
}
func New(size, cap int) *Queue{
	return &Queue{Data: make([]interface{}, size, cap), cap: cap, Mutex: &sync.Mutex{},}
}
func NewSize(size int) *Queue{
	return &Queue{Data: make([]interface{}, size, 2*size), Mutex: &sync.Mutex{},}
}
func NewCap(cap int) *Queue{
	return &Queue{Data: make([]interface{}, 0, cap), cap: cap, Mutex: &sync.Mutex{},}
}
func TimeQueue(expireAfter time.Duration,cap int,tsp time.Duration) *Queue{
	if tsp <= 0{
		tsp = 10*time.Second
	}
	return &Queue{
		Data:        make([]interface{}, 0, cap),
		timeSpy:     true,
		ExpireAfter: expireAfter,
		cap:         cap,
		timeStep:    tsp,
		Mutex:       &sync.Mutex{},
	}
}

func (q *Queue) Head() (interface{},int){
	if len(q.Data) != 0{
		return nil, -1
	}
	if !q.timeSpy{
		return q.Data[0],0
	}else{
		wrapper, index := q.Data[0],0
		return wrapper.(TimeWrapper).Data,index
	}
}

func (q *Queue) SafeHead() (interface{}, int) {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()
	return q.Head()
}