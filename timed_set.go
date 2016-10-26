package lww

import "time"

// TimedSet is the set where added element are associated with timestamps.
// If an element already exists, greater timestamps will always prevail.
// This data structure does not support removals.
//
// It could be used to as an underlying data structure for
// building a Last-Writer-Wins (LWW) element set.
type TimedSet interface {
	Add(interface{}, time.Time)
	AddedAt(interface{}) (time.Time, bool)
	Each(func(interface{}, time.Time)) error
}
