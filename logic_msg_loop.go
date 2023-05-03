package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

var g_mq *MessageQueue = NewMessageQueue()
var g_poolSize = 8

type CallBackData struct {
	msg      int
	callback func(*Msg)
}

type MessageQueue struct {
	messages  []*Msg
	callbacks []*CallBackData
	mutex     sync.Mutex
	notEmpty  *sync.Cond

	task chan *Msg
	wg   sync.WaitGroup
}

func NewMessageQueue() *MessageQueue {
	mq := &MessageQueue{
		messages: make([]*Msg, 0),
		task:     make(chan *Msg, g_poolSize),
	}
	mq.notEmpty = sync.NewCond(&mq.mutex)
	return mq
}

func (mq *MessageQueue) Push(message *Msg) {
	mq.mutex.Lock()
	mq.messages = append(mq.messages, message)
	mq.mutex.Unlock()
	mq.notEmpty.Signal()
}

func (mq *MessageQueue) Pop() *Msg {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	for len(mq.messages) == 0 {
		mq.notEmpty.Wait()
	}

	message := mq.messages[0]
	mq.messages = mq.messages[1:]
	return message
}

func (mq *MessageQueue) Add(id int, param string, data interface{}) {
	msg := &Msg{
		Id:    id,
		Param: param,
		Data:  data,
	}

	mq.Push(msg)
}

func (mq *MessageQueue) Register(msg int, c func(*Msg)) {
	data := new(CallBackData)
	data.msg = msg
	data.callback = c
	mq.callbacks = append(mq.callbacks, data)
}

func (mq *MessageQueue) GetRegister(msg int) func(*Msg) {
	for _, item := range mq.callbacks {
		if item.msg == msg {
			return item.callback
		}
	}
	return nil
}

func msg_loop_init() {
	for i := 0; i < g_poolSize; i++ {
		go worker(i, g_mq.task, &g_mq.wg)
	}
}

func msg_loop_process(msg *Msg) {
	c := g_mq.GetRegister(msg.Id & 0xF0)
	if c == nil {
		return
	}

	c(msg)
}

func msg_loop_run() {
	for {
		if len(g_mq.messages) == 0 {
			msg := new(Msg)
			msg.Id = MSG_PAGE_MAIN_IDLE
			g_mq.task <- msg
			time.Sleep(1 * time.Second)
			continue
		}
		msg := g_mq.Pop()
		// msg_loop_process(msg)
		g_mq.task <- msg
	}
}

// worker 是实际执行任务的工作者 Goroutine
func worker(id int, msgs <-chan *Msg, wg *sync.WaitGroup) {
	for msg := range msgs {
		time.Sleep(time.Microsecond * time.Duration(rand.Intn(10)))
		wg.Add(1)
		defer wg.Done()
		msg_loop_process(msg)
	}
	log.Printf("[%d/%d] exit\n", id, g_poolSize)
}
