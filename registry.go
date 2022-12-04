package workering

import (
	"context"
	"sync"
)

type WorkerFunction func(ctx context.Context, done chan<- any)

type RegisterSet struct {
	Name   string
	Worker WorkerFunction
}

var workers = make(map[string]*Worker)

func init() {
	workers = make(map[string]*Worker)
}

type Worker struct {
	name       string
	worker     WorkerFunction
	status     WorkerStatus
	ctx        *context.Context
	cancelFunc *context.CancelFunc
}

var mux sync.RWMutex

func Register(sets ...RegisterSet) {
	mux.Lock()
	defer mux.Unlock()
	for _, set := range sets {
		if set.Name == "" {
			panic("Worker name must not be empty")
		}
		if set.Worker == nil {
			panic("Worker must not be nil")
		}
		if _, ok := workers[set.Name]; ok {
			panic("Worker already registered")
		}

		workers[set.Name] = &Worker{
			name:   set.Name,
			worker: set.Worker,
			status: Stopped,
		}
	}
}

func Get(name string) *Worker {
	return workers[name]
}
