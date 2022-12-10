package workering

import (
	"context"
	"sync"
)

// WorkerFunction is the main workload function that is executed by a Worker
type WorkerFunction func(ctx context.Context, done chan<- any)
type waiter chan any

// RegisterSet describes a Worker used for registration with Register
type RegisterSet struct {
	Name   string
	Worker WorkerFunction
}

var workers = make(map[string]*Worker)

func init() {
	workers = make(map[string]*Worker)
}

// Worker represents a single Worker
type Worker struct {
	name           string
	workerFunction WorkerFunction
	status         WorkerStatus
	ctx            *context.Context
	cancelFunc     *context.CancelFunc
	waiters        []waiter
}

var mux sync.RWMutex

// Register registers one or more workers to the registry represented by RegisterSet
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
			name:           set.Name,
			workerFunction: set.Worker,
			status:         Stopped,
			waiters:        []waiter{},
		}
	}
}

// Get returns a workerFunction by name
func Get(name string) *Worker {
	return workers[name]
}
