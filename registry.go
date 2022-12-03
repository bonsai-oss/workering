package workering

import (
	"context"
	"sync"
)

type Worker func(ctx context.Context, done chan<- any)

var workers map[string]Worker
var status map[*Worker]WorkerStatus
var contexts map[*Worker]context.Context
var cancelFuncs map[*Worker]context.CancelFunc

var mux sync.RWMutex

func init() {
	workers = make(map[string]Worker)
	status = make(map[*Worker]WorkerStatus)
	contexts = make(map[*Worker]context.Context)
	cancelFuncs = make(map[*Worker]context.CancelFunc)
}

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
