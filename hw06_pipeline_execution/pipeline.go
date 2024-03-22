package hw06pipelineexecution

import (
	"sync"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

var (
	wg = sync.WaitGroup{}
	sm = sync.Mutex{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	result := make(chan interface{})
	values := make([]interface{}, 0)
	index := 0

	for value := range in {
		values = append(values, nil)
		stageInChan := make(chan interface{})
		wg.Add(1)

		go func(index int) {
			defer wg.Done()
			var stageOutChan In
			stageOutChan = stageInChan

			for _, stage := range stages {
				stageOutChan = stage(stageOutChan)
			}

			select {
			case stageValue := <-stageOutChan:
				sm.Lock()
				values[index] = stageValue
				sm.Unlock()
			case <-done:
				sm.Lock()
				values = nil
				sm.Unlock()
				return
			}
		}(index)

		stageInChan <- value
		sm.Lock()
		index++
		sm.Unlock()
	}

	wg.Wait()

	go func() {
		defer close(result)

		for _, value := range values {
			result <- value
		}
	}()

	return result
}
