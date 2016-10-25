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

// Add adds an element in the set if one of the following condition is met:
// - Given element does not exists yet
// - Given element already exists but with a lesser timestamp than the given one
func (s *TimedSet) Add(value interface{}, t time.Time) {
	addedAt, ok := s.AddedAt(value)
	if !ok || (ok && t.After(addedAt)) {
		s.elements[value] = t
	}
}

// AddedAt returns the timestamp of a given element if it exists.
//
// The second return value (bool) indicates whether the element exists or not.
// If the given element does not exists, the second return (bool) is false.
func (s *TimedSet) AddedAt(value interface{}) (time.Time, bool) {
	t, ok := s.elements[value]
	return t, ok
}
