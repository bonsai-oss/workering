package workering_test

import (
	"context"
	"strings"
	"testing"
	"time"

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

	workering.Register(workering.RegisterSet{
		Name:   "test-Worker",
		Worker: testWorkerBuilder(inputChannel, outputChannel),
	})

	t.Run("without waiters", func(t *testing.T) {
		worker := workering.Get("test-Worker")
		assert.Nil(t, worker.Start())
		inputChannel <- "hello"
		assert.Equal(t, "HELLO", <-outputChannel)
		assert.Nil(t, worker.Stop())
	})

	t.Run("with waiters", func(t *testing.T) {
		worker := workering.Get("test-Worker")

		var result1, result2 string

		go func(result *string) {
			<-worker.WaitStopped()
			*result = "done"
		}(&result1)
		go func(result *string) {
			<-worker.WaitStopped()
			*result = "done"
		}(&result2)

		assert.Nil(t, worker.Start())
		inputChannel <- "hello"
		assert.Equal(t, "HELLO", <-outputChannel)
		assert.Nil(t, worker.Stop())

		// must wait for the waiters to finish
		time.Sleep(50 * time.Millisecond)

		assert.Equal(t, "done", result1)
		assert.Equal(t, "done", result2)
	})

	t.Run("explicit workerFunction reusing", func(t *testing.T) {
		worker := workering.Get("test-Worker")
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
