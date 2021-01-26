package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var ErrNoWorkers = errors.New("worker count is zero. misconfig?")

type Task func() error

// Run starts tasks in `workerCount` goroutines and stops its work when receiving `errorCount` errors from tasks.
func Run(tasks []Task, workerCount int, errorCount int) (err error) {
	if workerCount <= 0 {
		return ErrNoWorkers
	}

	errCh := make(chan error)
	defer close(errCh)

	taskCh := make(chan Task)

	once := sync.Once{}
	termCh := make(chan struct{})
	terminate := func() {
		close(termCh)
	}

	wg := sync.WaitGroup{}

	// Start workers
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func() {
			defer wg.Done()
			worker(taskCh, errCh, termCh)
		}()
	}

	// Start error counting
	wg.Add(1)
	go func() {
		defer wg.Done()

		// err here is from Run scope
		err = errorCounter(errCh, termCh, errorCount)
		if err != nil {
			once.Do(terminate)
		}
	}()

	// Start task dealer
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(taskCh)
		defer once.Do(terminate)

		taskDealer(tasks, taskCh, termCh)
	}()

	wg.Wait()

	return err
}

func errorCounter(errCh <-chan error, termCh <-chan struct{}, threshold int) error {
	errorCounter := 0

	for {
		select {
		case <-termCh:
			return nil
		case _, ok := <-errCh:
			if !ok {
				return nil
			}

			errorCounter++
			if threshold > 0 && errorCounter >= threshold {
				return ErrErrorsLimitExceeded
			}
		}
	}
}

func taskDealer(tasks []Task, taskCh chan<- Task, termCh <-chan struct{}) {
	for _, task := range tasks {
		select {
		case taskCh <- task:
		case <-termCh:
			return
		}
	}
}

func worker(taskCh <-chan Task, errCh chan<- error, termCh <-chan struct{}) {
	for {
		select {
		case task, ok := <-taskCh:
			if !ok {
				return
			}

			err := task()
			if err != nil {
				select {
				case <-termCh:
					return
				case errCh <- err:
				}
			}
		case <-termCh:
			return
		}
	}
}
