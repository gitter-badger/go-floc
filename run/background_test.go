package run

import (
	"testing"

	floc "gopkg.in/workanator/go-floc.v1"
)

func TestBackgroundInactive(t *testing.T) {
	// Construct the flow control object.
	flow := floc.NewFlow()
	defer flow.Release()

	flow.Complete(nil)

	// Construct the state object which as data contains the counter.
	state := floc.NewState(new(int))
	defer state.Release()

	// Counstruct the result job.
	job := Background(jobIncrement)

	// Run the job.
	floc.Run(flow, state, updateCounter, job)

	if getCounter(state) != 0 {
		t.Fatalf("%s expects counter to be zero", t.Name())
	}
}
