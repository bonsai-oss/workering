package workering_test

import (
	"context"
	"strings"
	"testing"

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

		go func() {
			ret := <-worker.WaitStopped()
			assert.Equal(t, "done", ret)
		}()
		go func() {
			ret := <-worker.WaitStopped()
			assert.Equal(t, "done", ret)
		}()

		assert.Nil(t, worker.Start())
		inputChannel <- "hello"
		assert.Equal(t, "HELLO", <-outputChannel)
		assert.Nil(t, worker.Stop())
	})
}
