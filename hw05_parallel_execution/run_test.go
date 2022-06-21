package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)
	t.Run("the value m <= 0 is interpreted as a sign to ignore errors in principle;", func(t *testing.T) {
		tasksCount := rand.Intn(100)
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := rand.Intn(200)
		maxErrorsCount := 0
		err := Run(tasks, workersCount, maxErrorsCount)

		require.NoError(t, err, "actual err - %v", err)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks completed")

		maxErrorsCount = -rand.Intn(200)
		runTasksCount = 0

		err = Run(tasks, workersCount, maxErrorsCount)

		require.NoError(t, err, "actual err - %v", err)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks completed")
	})
	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var activeTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				atomic.AddInt32(&activeTasksCount, 1)
				time.Sleep(taskSleep)
				atomic.AddInt32(&activeTasksCount, -1)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		var wg sync.WaitGroup
		var err error
		wg.Add(1)
		go func() {
			err = Run(tasks, workersCount, maxErrorsCount)
			wg.Done()
		}()
		condition := func() bool {
			return atomic.LoadInt32(&activeTasksCount) == int32(workersCount)
		}
		require.Eventually(t, condition, sumTime, time.Millisecond/2, "tasks were run sequentially?")
		wg.Wait()
		require.NoError(t, err)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	})
}
