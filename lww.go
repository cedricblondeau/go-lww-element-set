package lww

import (
	"time"
)

// ElementSet is a Last-Writer-Wins (LWW) Element Set data structure
// implemented using two TimedSet data structures.
type ElementSet struct {
	additions TimedSet
	removals  TimedSet
}

// NewElementSet returns an empty and ready-to-use LWW element set
func NewElementSet() *ElementSet {
	return &ElementSet{
		additions: NewTimedSetMap(),
		removals:  NewTimedSetMap(),
	}
}

// Add marks an element to be added at a given timestamp
func (s *ElementSet) Add(value interface{}, t time.Time) {
	s.additions.Add(value, t)
}

// Remove marks an element to be removed at a given timestamp
func (s *ElementSet) Remove(value interface{}, t time.Time) {
	s.removals.Add(value, t)
}

// Exists checks if an element is marked as present in the set
func (s ElementSet) Exists(value interface{}) bool {
	addedAt, added := s.additions.AddedAt(value)
	if !added {
		return false
	}
	if !s.isRemoved(value, addedAt) {
		return true
	}
	return false
}

func (s ElementSet) isRemoved(value interface{}, since time.Time) bool {
	removedAt, removed := s.removals.AddedAt(value)
	if !removed {
		return false
	}
	if since.Before(removedAt) {
		return true
	}
	return false
}

// Get returns set content
func (s ElementSet) Get() []interface{} {
	var result []interface{}
	s.additions.Each(func(element interface{}, addedAt time.Time) {
		if !s.isRemoved(element, addedAt) {
			result = append(result, element)
		}
	})
	return result
}
