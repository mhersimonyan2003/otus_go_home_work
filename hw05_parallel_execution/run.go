package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func controller(doneChan chan struct{}, tasks []Task, m int) (chan Task, chan struct{}, *int) {
	tasksChan := make(chan Task)
	errorsChan := make(chan struct{})

	count := 0
	errors := 0

	go func() {
		defer close(doneChan)

		for {
			select {
			case tasksChan <- tasks[count]:
				count++

				if count == len(tasks) {
					return
				}
			case <-errorsChan:
				errors++

				if errors == m {
					return
				}
			}
		}
	}()

	return tasksChan, errorsChan, &errors
}

func consumer(doneChan chan struct{}, wg *sync.WaitGroup, tasksChan chan Task, errorsChan chan struct{}) {
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
				case errorsChan <- struct{}{}:
				}
			}
		}
	}
}

func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	wg := sync.WaitGroup{}
	doneChan := make(chan struct{})

	tasksChan, errorsChan, errors := controller(doneChan, tasks, m)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go consumer(doneChan, &wg, tasksChan, errorsChan)
	}

	wg.Wait()

	if *errors >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
