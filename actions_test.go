package workering_test

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/bonsai-oss/workering/v2"
)

func testWorkerBuilder(input, output chan string) workering.WorkerFunction {
	return func(ctx context.Context, done chan<- any) {
		defer (func() { done <- "" })()
		for {
			select {
			case <-ctx.Done():
				return
			case inputValue := <-input:
				output <- strings.ToUpper(inputValue)
			}
		}
	}
}

func TestWorker_Livecycle(t *testing.T) {
	inputChannel := make(chan string)
	outputChannel := make(chan string)

	t.Run("without waiters", func(t *testing.T) {
		workerName := uuid.New().String()
		workering.Register(workering.RegisterSet{
			Name:   workerName,
			Worker: testWorkerBuilder(inputChannel, outputChannel),
		})
		worker := workering.Get(workerName)
		assert.Nil(t, worker.Start())
		inputChannel <- "hello"
		assert.Equal(t, "HELLO", <-outputChannel)
		assert.Nil(t, worker.Stop())
	})

	t.Run("with waiters", func(t *testing.T) {
		workerName := uuid.New().String()
		workering.Register(workering.RegisterSet{
			Name:   workerName,
			Worker: testWorkerBuilder(inputChannel, outputChannel),
		})
		worker := workering.Get(workerName)

		assert.Nil(t, worker.Start())

		inputChannel <- "hello"
		assert.Equal(t, "HELLO", <-outputChannel)

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			assert.NotNil(t, <-worker.WaitStopped())
		}(&wg)

		time.Sleep(10 * time.Millisecond)

		assert.Nil(t, worker.Stop())
		wg.Wait()
	})

	t.Run("explicit workerFunction reusing", func(t *testing.T) {
		workerName := uuid.New().String()
		workering.Register(workering.RegisterSet{
			Name:   workerName,
			Worker: testWorkerBuilder(inputChannel, outputChannel),
		})
		worker := workering.Get(workerName)
		assert.Nil(t, worker.Start())
		inputChannel <- "hello"
		assert.Equal(t, "HELLO", <-outputChannel)
		assert.Nil(t, worker.Stop())

		assert.Nil(t, worker.Start())
		inputChannel <- "world"
		assert.Equal(t, "WORLD", <-outputChannel)
		assert.Nil(t, worker.Stop())
	})
}
