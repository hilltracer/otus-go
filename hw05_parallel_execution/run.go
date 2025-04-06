package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
//
//nolint:gocognit // The function is divided logically, cognitive complexity is acceptable
func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return nil
	}

	if len(tasks) == 0 {
		return nil
	}

	// If m <= 0 then ignore errors, perform all tasks
	ignoreErrors := m <= 0

	tasksCh := make(chan Task)
	doneCh := make(chan struct{})

	var wg sync.WaitGroup
	var once sync.Once

	var errorsCount int32

	startWorker := func() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-doneCh:
					// If signal to stop is received, finish the processing of tasks
					return
				case task, ok := <-tasksCh:
					if !ok {
						// If tasksCh is closed, no more tasks to process
						return
					}
					if err := task(); err != nil && !ignoreErrors {
						newCount := atomic.AddInt32(&errorsCount, 1)
						if int(newCount) >= m {
							once.Do(func() {
								// Add once to prevent closing the closed channel
								close(doneCh)
							})
						}
					}
				}
			}
		}()
	}

	// Start n workers
	for range n {
		startWorker()
	}

	// Send tasks to workers
	go func() {
		defer close(tasksCh)
		for _, task := range tasks {
			select {
			case <-doneCh:
				// If signal to stop is received, finish sending tasks
				return
			default:
			}
			tasksCh <- task
		}
	}()

	wg.Wait()

	if !ignoreErrors && int(errorsCount) >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
