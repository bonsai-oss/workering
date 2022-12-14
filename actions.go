package workering

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/go-multierror"
)

// Start starts a Worker
func (w *Worker) Start() error {
	mux.Lock()
	defer mux.Unlock()

	if w.status == Running {
		return fmt.Errorf("workerFunction already running")
	}

	ctx, cancel := context.WithCancel(context.Background())
	w.ctx = &ctx
	w.cancelFunc = &cancel

	done := make(chan any)
	go w.workerFunction(ctx, done)
	go w.watchdog(done)

	w.status = Running
	return nil
}

// watchdog watches for a done channel and manages the status of the workerFunction
func (w *Worker) watchdog(done chan any) {
	<-done
	defer func() {
		w.status = Stopped
		w.stoppedChan <- true
	}()

	if len(w.waiters) == 0 {
		return
	}

	wg := sync.WaitGroup{}
	for waiterIndex := range w.waiters {
		wg.Add(1)
		go func(wg *sync.WaitGroup, wd waiter) {
			defer wg.Done()
			wd <- true
			close(wd)
		}(&wg, w.waiters[waiterIndex])
	}
	wg.Wait()

	w.waiters = []waiter{}
}

// WaitStopped registers a waiter and returns a channel informed when the WorkerFunction is stopped
func (w *Worker) WaitStopped() <-chan any {
	var waiterInstance = make(waiter)
	w.waiters = append(w.waiters, waiterInstance)
	return waiterInstance
}

// Stop stops a Worker and waits for it to stop
func (w *Worker) Stop() error {
	mux.Lock()
	defer mux.Unlock()

	if w.status == Stopped {
		return fmt.Errorf("workerFunction already stopped")
	}
	(*w.cancelFunc)()
	<-w.stoppedChan

	// clear context and cancelFunc for reusing
	w.ctx = nil
	w.cancelFunc = nil

	return nil
}

// Status returns the Status of a Worker
func (w *Worker) Status() WorkerStatus {
	return w.status
}

// StartAll starts all Worker instances
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

// StopAll stops all Worker instances and waits for them
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
