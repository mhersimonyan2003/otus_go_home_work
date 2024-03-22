package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrInvalidWorkersLimit = errors.New("invalid workers limit")
	wg                     = sync.WaitGroup{}
)

type Task func() error

func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return ErrInvalidWorkersLimit
	}
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	var ops atomic.Int32

	controller := func(doneChan chan struct{}) chan Task {
		tasksChan := make(chan Task)

		go func() {
			defer close(doneChan)

			for _, task := range tasks {
				tasksChan <- task

				if int(ops.Load()) >= m {
					return
				}
			}
		}()

		return tasksChan
	}

	consumer := func(doneChan chan struct{}, tasksChan chan Task) {
		defer wg.Done()

		for {
			select {
			case <-doneChan:
				return
			case task := <-tasksChan:
				err := task()
				if err != nil {
					select {
					case <-doneChan:
						return
					default:
						ops.Add(1)
					}
				}
			}
		}
	}

	doneChan := make(chan struct{})

	tasksChan := controller(doneChan)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go consumer(doneChan, tasksChan)
	}

	wg.Wait()

	if int(ops.Load()) >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
