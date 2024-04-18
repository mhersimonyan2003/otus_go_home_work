package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		inChan := make(Bi)

		go func(in In) {
			defer close(inChan)
			for {
				select {
				case <-done:
					return
				case val, ok := <-in:
					if !ok {
						return
					}
					inChan <- val
				}
			}
		}(in)

		in = stage(inChan)
	}

	return in
}
