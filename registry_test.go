package workering_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/bonsai-oss/workering/v2"
)

func TestRegister(t *testing.T) {
	// empty Worker
	assert.Panics(t, func() {
		workerName := uuid.New().String()
		workering.Register(workering.RegisterSet{
			Name:   workerName,
			Worker: nil,
		})
	})

	// empty name
	assert.Panics(t, func() {
		workering.Register(workering.RegisterSet{
			Name:   "",
			Worker: testWorkerBuilder(nil, nil),
		})
	})

	// normal register call
	assert.NotPanics(t, func() {
		workerName := uuid.New().String()
		workering.Register(workering.RegisterSet{
			Name:   workerName,
			Worker: testWorkerBuilder(nil, nil),
		})

		assert.NotNil(t, workering.Get(workerName))
	})

	// panic on duplicate register
	assert.Panics(t, func() {
		workerName := uuid.New().String()
		workering.Register(workering.RegisterSet{
			Name:   workerName,
			Worker: testWorkerBuilder(nil, nil),
		})
		workering.Register(workering.RegisterSet{
			Name:   workerName,
			Worker: testWorkerBuilder(nil, nil),
		})
	})
}
