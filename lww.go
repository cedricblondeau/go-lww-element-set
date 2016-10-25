package lww

import (
	"time"
)

// ElementSet is a Last-Writer-Wins (LWW) Element Set data structure
// implemented using two TimedSet data structures.
type ElementSet struct {
	additions *TimedSet
	removals  *TimedSet
}

// NewElementSet returns an empty and ready-to-use LWW element set
func NewElementSet() *ElementSet {
	return &ElementSet{
		additions: NewTimedSet(),
		removals:  NewTimedSet(),
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
