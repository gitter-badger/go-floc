package floc

//go:generate stringer -type=Result

// Result is the result of flow execution.
type Result int32

// Possible results.
const (
	None Result = iota
	Completed
	Canceled

	resultFirst = None
	resultLast  = Canceled
)

// IsNone tests if the resilt is None.
func (result Result) IsNone() bool {
	return result == None
}

// IsCompleted tests if the resilt is Completed.
func (result Result) IsCompleted() bool {
	return result == Completed
}

// IsCanceled tests if the resilt is Canceled.
func (result Result) IsCanceled() bool {
	return result == Canceled
}

// IsFinished tests if the result is either Completed or Canceled.
func (result Result) IsFinished() bool {
	return result == Completed || result == Canceled
}

// IsValid tests if the result is a valid value.
func (result Result) IsValid() bool {
	return result >= resultFirst && result <= resultLast
}

// Int32 returns the underlying value as int32. That is handy while working
// with atomic operations.
func (result Result) Int32() int32 {
	return int32(result)
}
