package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var wg = sync.WaitGroup{}
var sm = sync.Mutex{}

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	// Place your code here.
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	count := 0

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for count < len(tasks) {
				if m <= 0 {
					return
				}

				sm.Lock()
				task := tasks[count]
				count++
				sm.Unlock()

				err := task()

				sm.Lock()
				if err != nil {
					m--
				}
				sm.Unlock()
			}
		}()
	}

	wg.Wait()

	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}
