package lww

import "time"

// timedSet is the set where added element are associated with timestamps.
// If an element already exists, greater timestamps will always prevail.
// This data structure does not support removals.
//
// It's useful as an underlying data structure for
// building a Last-Writer-Wins (LWW) element set.
type timedSet interface {
	add(interface{}, time.Time) error
	addedAt(interface{}) (time.Time, bool, error)
	each(func(interface{}, time.Time) error) error
}
