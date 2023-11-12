package workerpool_test

import (
	"errors"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/workerpool"
	"io"
	"testing"
	"time"
)

const testWorkerPoolSize = 10
const testWorkerPoolJobCount = 100
const testWorkerPoolTimeout = 1 * time.Second

func TestWorkerPool(t *testing.T) {
	job := func(id int) int {
		return id
	}

	timeout := time.After(testWorkerPoolTimeout)

	pool := workerpool.New[int](testWorkerPoolSize)
	pool.Start()
	defer pool.Stop()

	for i := 0; i < testWorkerPoolJobCount; i++ {
		pool.Submit(job)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < testWorkerPoolJobCount; i++ {
			id, err := pool.Result()
			if errors.Is(err, io.EOF) {
				t.Errorf("channel closed unexpectedly")
				return
			}
			if id != i%testWorkerPoolSize {
				t.Errorf("expected %d, got %d", i, id)
				return
			}
		}
	}()

	select {
	case <-timeout:
		t.Fatalf("timed out after %s", testWorkerPoolTimeout.String())
	case <-done:
		break
	}
}
