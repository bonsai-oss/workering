package workering

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
)

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

	delete(contexts, w)
	delete(cancelFuncs, w)
	return nil
}

func (w *Worker) Status() WorkerStatus {
	return status[w]
}

func StartAll() error {
	var err error
	for _, worker := range workers {
		if (worker).Status() == Running {
			continue
		}
		if startError := worker.Start(); startError != nil {
			err = multierror.Append(err, startError)
		}
	}

	return err
}
func StopAll() error {
	var err error
	for _, worker := range workers {
		if (worker).Status() == Running {
			continue
		}
		if startError := worker.Stop(); startError != nil {
			err = multierror.Append(err, startError)
		}
	}

	return err
}
