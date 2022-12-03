package workering

import (
	"context"
	"fmt"
	"sync"
)

var workers map[string]Worker
var status map[*Worker]int
var contexts map[*Worker]context.Context
var cancelFuncs map[*Worker]context.CancelFunc

var mux sync.RWMutex

const (
	_ = iota
	Running
	Stopped
)

func init() {
	workers = make(map[string]Worker)
	status = make(map[*Worker]int)
	contexts = make(map[*Worker]context.Context)
	cancelFuncs = make(map[*Worker]context.CancelFunc)
}

type Worker func(ctx context.Context, done chan<- any)

func Register(name string, worker Worker) {
	if workers == nil {
		workers = make(map[string]Worker)
	}
	if _, ok := workers[name]; ok {
		panic("worker already registered")
	}
	workers[name] = worker
	status[&worker] = Stopped

}

func Get(name string) Worker {
	return workers[name]
}

func (w *Worker) Start() error {
	mux.Lock()
	defer mux.Unlock()

	if status[w] == Running {
		return fmt.Errorf("worker already running")
	}

	ctx, cancel := context.WithCancel(context.Background())
	contexts[w] = ctx
	cancelFuncs[w] = cancel

	done := make(chan any)
	go (*w)(ctx, done)
	go w.watchdog(done)

	status[w] = Running
	return nil
}

func (w *Worker) watchdog(done chan any) {
	<-done
	mux.Lock()
	defer mux.Unlock()
	status[w] = Stopped
}

func (w *Worker) Stop() error {
	mux.Lock()
	defer mux.Unlock()

	if status[w] == Stopped {
		return fmt.Errorf("worker already stopped")
	}
	cancelFuncs[w]()
	status[w] = Stopped
	return nil
}

func (w *Worker) Status() int {
	return status[w]
}
