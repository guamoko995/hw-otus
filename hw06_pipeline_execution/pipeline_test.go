package hw06pipelineexecution

import (
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

const (
	sleepPerStage = time.Millisecond * 100
	fault         = sleepPerStage / 2
)

func TestPipeline(t *testing.T) {
	var activeTasksCount int32

	// Stage generator
	g := func(_ string, f func(v interface{}) interface{}) Stage {
		return func(in In) Out {
			out := make(Bi)
			go func() {
				atomic.AddInt32(&activeTasksCount, 1)
				defer atomic.AddInt32(&activeTasksCount, -1)
				defer close(out)
				for v := range in {
					time.Sleep(sleepPerStage)
					out <- f(v)
				}
			}()
			return out
		}
	}

	stages := []Stage{
		g("Dummy", func(v interface{}) interface{} { return v }),
		g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
		g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
		g("Stringifier", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
	}

	t.Run("simple case", func(t *testing.T) {
		in := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, nil, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)
		activeTaskCoun := atomic.LoadInt32(&activeTasksCount)

		require.Equal(t, []string{"102", "104", "106", "108", "110"}, result)
		require.Less(t,
			int64(elapsed),
			// ~0.8s for processing 5 values in 4 stages (100ms every) concurrently
			int64(sleepPerStage)*int64(len(stages)+len(data)-1)+int64(fault))
		require.Equal(t, int32(0), activeTaskCoun, "not all goroutines completed")
	})

	t.Run("done case", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		in := make(Bi)
		done := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		// Abort after 800ms
		abortDur := sleepPerStage * 2
		// Time stopet max = time stage (get Done) + time stage (1 step stage) + fault
		stopDur := sleepPerStage*2 + fault

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			<-time.After(abortDur)
			close(done)
			time.Sleep(stopDur)
			wg.Done()
		}()

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 5)
		start := time.Now()
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)
		wg.Wait()
		activeTaskCoun := atomic.LoadInt32(&activeTasksCount)

		require.Len(t, result, 0)
		require.Less(t, int64(elapsed), int64(abortDur)+int64(fault))

		require.Equal(t, int32(0), activeTaskCoun, "not all goroutines completed during the stop time")
	})
}
