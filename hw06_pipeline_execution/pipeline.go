package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stg := range stages {
		in = mediator(done, in)
		in = stg(in)
	}
	return mediator(done, in)
}

// mediator - буфер между stages, позволяющий закрыть входные каналы каждого stage и
// продолжать принимать данные с выхода каждого stage до закрытия выходного канала.
// Т.о. mediator позволяет завершить все горутины pipeline-а любой длины не более чем
// за два времени выполнения самого длителного stage.
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
				for {
					if _, ok := <-in; !ok {
						return
					}
				}
			}
		}
	}()
	return bi
}
