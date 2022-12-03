package workering_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bonsai-oss/workering"
)

func testWorkerBuilder(input, output chan string) workering.Worker {
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

func TestWorker_Start(t *testing.T) {
	inputChannel := make(chan string)
	outputChannel := make(chan string)

	workering.Register(workering.RegisterSet{
		Name:   "test-worker",
		Worker: testWorkerBuilder(inputChannel, outputChannel),
	})

	worker := workering.Get("test-worker")
	assert.Nil(t, worker.Start())
	inputChannel <- "hello"
	assert.Equal(t, "HELLO", <-outputChannel)
	assert.Nil(t, worker.Stop())
}
