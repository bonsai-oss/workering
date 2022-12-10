package workering

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/go-multierror"
)

// Start starts a worker
func (w *Worker) Start() error {
	mux.Lock()
	defer mux.Unlock()

	if w.status == Running {
		return fmt.Errorf("worker already running")
	}

	ctx, cancel := context.WithCancel(context.Background())
	w.ctx = &ctx
	w.cancelFunc = &cancel

	done := make(chan any)
	go w.worker(ctx, done)
	go w.watchdog(done)

	w.status = Running
	return nil
}

// watchdog watches for a done channel and manages the status of the worker
func (w *Worker) watchdog(done chan any) {
	<-done
	defer func() { w.status = Stopped }()

	if len(w.waiters) == 0 {
		return
	}

	wg := sync.WaitGroup{}
	for waiterIndex := range w.waiters {
		wg.Add(1)
		go func(wg *sync.WaitGroup, wd waiter) {
			defer wg.Done()
			wd <- interface{}(nil)
		}(&wg, w.waiters[waiterIndex])
	}
	wg.Wait()

	// cleanup waiters
	for waiterIndex := range w.waiters {
		close(w.waiters[waiterIndex])
	}
	w.waiters = []waiter{}
}

// WaitStopped registers a waiter and returns a channel informed when the worker is stopped
func (w *Worker) WaitStopped() <-chan any {
	var waiterInstance = make(waiter)
	w.waiters = append(w.waiters, waiterInstance)
	return waiterInstance
}

// Stop stops a worker and waits for it to stop
func (w *Worker) Stop() error {
	mux.Lock()
	defer mux.Unlock()

	if w.status == Stopped {
		return fmt.Errorf("worker already stopped")
	}
	(*w.cancelFunc)()

	<-w.WaitStopped()

	// clear context and cancelFunc for reusing
	w.ctx = nil
	w.cancelFunc = nil

	return nil
}

// Status returns a workers status
func (w *Worker) Status() WorkerStatus {
	return w.status
}

// StartAll starts all workers
func StartAll() error {
	var err error
	for _, worker := range workers {
		if worker.Status() == Running {
			continue
		}
		if startError := worker.Start(); startError != nil {
			err = multierror.Append(err, startError)
		}
	}

	return err
}

// StopAll stops all workers and waits for them to stop
func StopAll() error {
	var err error
	for _, worker := range workers {
		if worker.Status() == Stopped {
			continue
		}
		if startError := worker.Stop(); startError != nil {
			err = multierror.Append(err, startError)
		}
	}

	return err
}
