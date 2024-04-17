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

			for value := range in {
				select {
				case <-done:
					return
				default:
				}

				inChan <- value
			}
		}(in)

		in = stage(inChan)
	}

	return in
}
