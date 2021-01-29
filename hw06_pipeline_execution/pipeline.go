package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		in = stage(makeClosableBroker(in, done))
	}
	return in
}

// makeClosableBroker creates a broker goroutine
// which can be terminated via closing `interrupt` channel.
func makeClosableBroker(in, interrupt In) Out {
	out := make(Bi)

	go func() {
		defer close(out)

		for {
			select {
			case data, ok := <-in:
				if !ok {
					return
				}
				out <- data
			case <-interrupt:
				return
			}
		}
	}()

	return out
}
