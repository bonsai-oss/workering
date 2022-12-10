package workering

import (
	"context"
	"sync"
)

// WorkerFunction is the function that is executed by a worker
type WorkerFunction func(ctx context.Context, done chan<- any)
type waiter chan any

// RegisterSet represents a single worker used for registration with Register
type RegisterSet struct {
	Name   string
	Worker WorkerFunction
}

var workers = make(map[string]*Worker)

func init() {
	workers = make(map[string]*Worker)
}

// Worker represents a single worker
type Worker struct {
	name       string
	worker     WorkerFunction
	status     WorkerStatus
	ctx        *context.Context
	cancelFunc *context.CancelFunc
	waiters    []waiter
}

var mux sync.RWMutex

// Register registers one or more workers to the registry
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
			name:    set.Name,
			worker:  set.Worker,
			status:  Stopped,
			waiters: []waiter{},
		}
	}
}

// Get returns a worker by name
func Get(name string) *Worker {
	return workers[name]
}
