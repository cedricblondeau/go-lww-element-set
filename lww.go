package lww

import (
	"strings"
	"time"

	redis "gopkg.in/redis.v4"
)

// ElementSet is a Last-Writer-Wins (LWW) Element Set data structure
// implemented using two TimedSet data structures.
type ElementSet struct {
	additions TimedSet
	removals  TimedSet
}

// NewMapElementSet returns an empty and ready-to-use map-backed LWW element set
func NewMapElementSet() *ElementSet {
	return &ElementSet{
		additions: NewMapTimedSet(),
		removals:  NewMapTimedSet(),
	}
}

// NewRedisElementSet returns an empty and ready-to-use redis-backed LWW element set
func NewRedisElementSet(prefixKey string, c *redis.Client) *ElementSet {
	return &ElementSet{
		additions: NewRedisTimedSet(strings.Join([]string{prefixKey, "_lww_additions"}, ""), c),
		removals:  NewRedisTimedSet(strings.Join([]string{prefixKey, "_lww_removals"}, ""), c),
	}
}

// NewRedisElementSetWithCustomMarshalling returns an empty
// and ready-to-use redis-backed LWW element set with
// custom marshalling and unmarshalling functions
func NewRedisElementSetWithCustomMarshalling(
	prefixKey string,
	c *redis.Client,
	marshal func(value interface{}) string,
	unmarshal func(value string) interface{}) *ElementSet {

	additions := NewRedisTimedSet(strings.Join([]string{prefixKey, "_lww_additions"}, ""), c)
	additions.marshal = marshal
	additions.unmarshal = unmarshal

	removals := NewRedisTimedSet(strings.Join([]string{prefixKey, "_lww_removals"}, ""), c)
	removals.marshal = marshal
	removals.unmarshal = unmarshal

	return &ElementSet{
		additions: additions,
		removals:  removals,
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
