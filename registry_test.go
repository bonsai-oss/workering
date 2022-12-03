package workering_test

import (
	"context"
	"github.com/bonsai-oss/workering"
	"log"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	workerA := func(ctx context.Context, done chan<- any) {
		defer (func() { done <- "" })()
		ticker := time.NewTicker(1000 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				log.Println("tick")
			}
		}
	}

	workering.Register("workerA", workerA)

	worker := workering.Get("workerA")

	err := worker.Start()
	if err != nil {
		log.Println(err)
	}
	time.Sleep(2 * time.Second)
	err = worker.Stop()
	if err != nil {
		log.Println(err)
	}

}
