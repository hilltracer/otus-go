package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.

func Run(tasks []Task, n, m int) error {
	if n <= 0 || len(tasks) == 0 {
		return nil
	}

	// If m <= 0 then ignore errors, perform all tasks
	ignoreErrors := m <= 0

	tasksCh := make(chan Task)
	var errorsCount int32
	var wg sync.WaitGroup

	// Start n workers
	for range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasksCh {
				if err := task(); err != nil && !ignoreErrors {
					atomic.AddInt32(&errorsCount, 1)
				}
			}
		}()
	}

	// Send tasks to workers
	for _, task := range tasks {
		if !ignoreErrors && int(atomic.LoadInt32(&errorsCount)) >= m {
			break
		}
		tasksCh <- task
	}
	close(tasksCh)

	wg.Wait()

	if !ignoreErrors && int(errorsCount) >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
