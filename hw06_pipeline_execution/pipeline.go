package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stg := range stages {
		in = stg(in)
	}
	return mediator(done, in)
}

func mediator(done In, in In) Out {
	bi := make(Bi)
	go func() {
		for {
			select {
			case v, ok := <-in:
				if !ok {
					close(bi)
					return
				}
				bi <- v
			case <-done:
				close(bi)
				return
			}
		}
	}()
	return bi
}
