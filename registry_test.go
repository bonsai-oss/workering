package workering_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bonsai-oss/workering/v2"
)

func TestRegister(t *testing.T) {
	// empty Worker
	assert.Panics(t, func() {
		workering.Register(workering.RegisterSet{
			Name:   "test-Worker",
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
		workering.Register(workering.RegisterSet{
			Name:   "test-worker2",
			Worker: testWorkerBuilder(nil, nil),
		})

		assert.NotNil(t, workering.Get("test-Worker"))
	})

	// panic on duplicate register
	assert.Panics(t, func() {
		workering.Register(workering.RegisterSet{
			Name:   "test-Worker-duplicate",
			Worker: testWorkerBuilder(nil, nil),
		})
		workering.Register(workering.RegisterSet{
			Name:   "test-Worker-duplicate",
			Worker: testWorkerBuilder(nil, nil),
		})
	})
}
