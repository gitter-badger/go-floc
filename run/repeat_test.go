package run

import (
	"testing"

	floc "gopkg.in/workanator/go-floc.v1"
	"gopkg.in/workanator/go-floc.v1/guard"
)

func TestRepeat(t *testing.T) {
	// Construct the flow control object.
	flow := floc.NewFlow()
	defer flow.Release()

	// Construct the state object which as data contains the counter.
	state := floc.NewState(new(int))
	defer state.Release()

	// Counstruct the result job.
	const times = 10

	job := Repeat(
		times,
		jobIncrement,
	)

	// Run the job.
	floc.Run(flow, state, updateCounter, job)

	expect := times
	v := getCounter(state)
	if v != expect {
		t.Fatalf("%s expects counter to be %d but has %d", t.Name(), expect, v)
	}
}

func TestRepeatInterrupt(t *testing.T) {
	// Construct the flow control object.
	flow := floc.NewFlow()
	defer flow.Release()

	// Construct the state object which as data contains the counter.
	state := floc.NewState(new(int))
	defer state.Release()

	// Counstruct the result job.
	const times = 10

	job := Repeat(
		times,
		Sequence(
			jobIncrement,
			guard.Complete(nil),
		),
	)

	// Run the job.
	floc.Run(flow, state, updateCounter, job)

	expect := 1
	v := getCounter(state)
	if v != expect {
		t.Fatalf("%s expects counter to be %d but has %d", t.Name(), expect, v)
	}
}
