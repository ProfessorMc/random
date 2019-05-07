package random

import "testing"

func TestNewQuantizer(t *testing.T) {
	q := NewQuantizer()
	q.SetIntervals(0)
	q.CreateIntervals()
	t.Log(len(q.IntervalMap))
	q.EvaluateMatches(4)
}

