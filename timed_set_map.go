package lww

import (
	"sync"
	"time"
)

// MapTimedSet is a TimedSet implementation that uses a map data structure.
// Maps in Go are not thread safe by default and that's why we use mutual exclusion.
type MapTimedSet struct {
	elements map[interface{}]time.Time
	l        sync.RWMutex // we name it because we don't want to expose it
}

// NewMapTimedSet returns an empty and ready-to-use map-backed TimedSet data structure.
func NewMapTimedSet() *MapTimedSet {
	return &MapTimedSet{
		elements: make(map[interface{}]time.Time),
	}
}

// Add adds an element in the set if one of the following condition is met:
// - Given element does not exists yet
// - Given element already exists but with a lesser timestamp than the given one
//
// This function is thread-safe.
func (s *MapTimedSet) Add(value interface{}, t time.Time) {
	s.l.Lock()
	defer s.l.Unlock()
	addedAt, ok := s.elements[value]
	if !ok || (ok && t.After(addedAt)) {
		s.elements[value] = t
	}
}

// AddedAt returns the timestamp of a given element if it exists.
//
// The second return value (bool) indicates whether the element exists or not.
// If the given element does not exists, the second return (bool) is false.
//
// This function is thread-safe.
func (s *MapTimedSet) AddedAt(value interface{}) (time.Time, bool) {
	s.l.RLock()
	defer s.l.RUnlock()
	t, ok := s.elements[value]
	return t, ok
}

// Each traverses the items in the TimedSet, calling the provided function
// for each element/timestamp association.
func (s *MapTimedSet) Each(f func(element interface{}, addedAt time.Time)) error {
	s.l.RLock()
	defer s.l.RUnlock()
	for element, addedAt := range s.elements {
		f(element, addedAt)
	}
	return nil
}
