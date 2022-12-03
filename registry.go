package workering

import (
	"context"
	"sync"
)

type Worker func(ctx context.Context, done chan<- any)

type RegisterSet struct {
	Name   string
	Worker Worker
}

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

func Register(sets ...RegisterSet) {
	mux.Lock()
	defer mux.Unlock()
	for _, set := range sets {
		if set.Name == "" {
			panic("worker name must not be empty")
		}
		if set.Worker == nil {
			panic("worker must not be nil")
		}
		if _, ok := workers[set.Name]; ok {
			panic("worker already registered")
		}
		workers[set.Name] = set.Worker
		status[&set.Worker] = Stopped
	}
}

func Get(name string) Worker {
	return workers[name]
}
