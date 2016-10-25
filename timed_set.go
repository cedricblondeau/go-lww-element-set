package lww

import "time"

// TimedSet is the set where added element are associated with timestamps.
// If an element already exists, greater timestamps will always prevail.
// This data structure does not support removals.
//
// It could be used to as an underlying data structure for
// building a Last-Writer-Wins (LWW) element set.
//
// This implementation uses a map data structure.
type TimedSet struct {
	elements map[interface{}]time.Time
}

// NewTimedSet returns an empty and ready-to-use TimedSet data structure.
func NewTimedSet() *TimedSet {
	return &TimedSet{
		elements: make(map[interface{}]time.Time),
	}
}

// Add adds an element in the set.
func (s *TimedSet) Add(value interface{}, t time.Time) {
	s.elements[value] = t
}
