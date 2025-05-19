package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func bridge(done In, in In) Out {
	if done == nil {
		return in
	}

	out := make(Bi)
	go func() {
		for {
			select {
			case v, ok := <-in:
				if !ok {
					close(out)
					return
				}
				select {
				case out <- v:
				case <-done:
					close(out)
					go func() {
						for range in {
							_ = struct{}{} // drain
						}
					}()
					return
				}
			case <-done:
				close(out)
				go func() {
					for range in {
						_ = struct{}{} // drain
					}
				}()
				return
			}
		}
	}()
	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return in
	}

	// If we used bridge here, it might stop reading from `in` early,
	// causing the generator goroutine to block on send and leak.
	curr := stages[0](in)

	for _, stage := range stages[1:] {
		curr = stage(bridge(done, curr))
	}

	// wrap the final output with a bridge to close it immediately on `done`
	return bridge(done, curr)
}
