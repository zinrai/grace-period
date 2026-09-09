package main

import "testing"

// Getting any of the four the wrong way round either hands a deadline to
// someone who is already set up, or leaves someone without one for good
func TestDeadlineFollowsWhetherSetupIsDone(t *testing.T) {
	for _, c := range []struct {
		hasMarker, hasDeadline bool
		want                   action
	}{
		{false, false, set},  // nothing done yet, so start the clock
		{false, true, leave}, // the clock is running, let it run
		{true, true, clear},  // done in time, take the deadline away
		{true, false, leave}, // done and no deadline, nothing to do
	} {
		got := decide(c.hasMarker, c.hasDeadline)
		if got != c.want {
			t.Errorf("hasMarker=%v hasDeadline=%v -> %v, want %v",
				c.hasMarker, c.hasDeadline, got, c.want)
		}
	}
}
