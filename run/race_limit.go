package run

import (
	"sync"

	floc "github.com/workanator/go-floc"

	"github.com/workanator/go-floc/flow"
)

// RaceLimit runs jobs in parallel and waits until the first N jobs finish.
func RaceLimit(limit int, jobs ...floc.Job) floc.Job {
	return func(theFlow floc.Flow, state floc.State, update floc.Update) {
		// Validate the winner limit
		if limit < 1 || limit > len(jobs) {
			panic("invalid amount of possible race winners")
		}

		// Do not start the race if the execution is finished
		if theFlow.IsFinished() {
			return
		}

		// Create the channel which will have the value when the race is won
		done := make(chan int, limit)
		defer close(done)

		// Wrap the flow into disablable flow so the calls to Cancel and Complete
		// can be disabled when the race is won
		disFlow, disable := flow.WithDisable(theFlow)

		// Wrap the trigger to a function which allows to hit the trigger only
		// `limit` time(s)
		mutex := sync.Mutex{}
		won := 0

		limitedUpdate := func(flow floc.Flow, state floc.State, key string, value interface{}) {
			mutex.Lock()
			defer mutex.Unlock()

			if won < limit {
				won++
				update(flow, state, key, value)
				done <- won

				if won == limit {
					disable()
				}
			}
		}

		// Condition is used to synchronize start of jobs
		allAreReady := false
		startCond := sync.NewCond(&sync.Mutex{})

		// Run jobs in parallel and wait untill all of them ready to start
		for _, job := range jobs {
			go func(job floc.Job) {
				// Wait for the start of the race
				startCond.L.Lock()
				for !allAreReady {
					// Test if the flow is finished
					if theFlow.IsFinished() {
						return
					}

					// Wait the time the race starts on
					startCond.Wait()
				}
				startCond.L.Unlock()

				// Perform thejob
				job(disFlow, state, limitedUpdate)
			}(job)
		}

		// Notify all jobs they can start the race
		startCond.L.Lock()
		allAreReady = true
		startCond.Broadcast()
		startCond.L.Unlock()

		// Wait until `limit` job(s) done
		winners := 0
		for winners < limit {
			select {
			case <-disFlow.Done():
				// The execution has been finished or canceled so the trigger
				// should not be triggered anymore, therefore we disable it
				disable()
				return

			case <-done:
				// One of the jobs finished. We have a winner! Count it.
				winners++
			}
		}
	}
}
