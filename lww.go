package lww

import (
	"strings"
	"time"

	redis "gopkg.in/redis.v4"
)

// ElementSet is a Last-Writer-Wins (LWW) Element Set data structure
// implemented using two TimedSet data structures.
type ElementSet struct {
	additions timedSet
	removals  timedSet
}

// NewMapElementSet returns an empty and ready-to-use map-backed LWW element set
func NewMapElementSet() *ElementSet {
	return &ElementSet{
		additions: newMapTimedSet(),
		removals:  newMapTimedSet(),
	}
}

// NewRedisElementSet returns an empty and ready-to-use redis-backed LWW element set
func NewRedisElementSet(prefixKey string, c *redis.Client) *ElementSet {
	return &ElementSet{
		additions: newRedisTimedSet(strings.Join([]string{prefixKey, "_lww_additions"}, ""), c),
		removals:  newRedisTimedSet(strings.Join([]string{prefixKey, "_lww_removals"}, ""), c),
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

	additions := newRedisTimedSet(strings.Join([]string{prefixKey, "_lww_additions"}, ""), c)
	additions.marshal = marshal
	additions.unmarshal = unmarshal

	removals := newRedisTimedSet(strings.Join([]string{prefixKey, "_lww_removals"}, ""), c)
	removals.marshal = marshal
	removals.unmarshal = unmarshal

	return &ElementSet{
		additions: additions,
		removals:  removals,
	}
}

// Add marks an element to be added at a given timestamp.
//
// Because of eventual consistencty, we do not guarantee
// that the operation took effect, we just return an error
// if a network error occured.
func (s *ElementSet) Add(value interface{}, t time.Time) error {
	return s.additions.add(value, t)
}

// Remove marks an element to be removed at a given timestamp.
//
// Because of eventual consistencty, we do not guarantee
// that the operation took effect, we just return an error
// if a network error occured.
func (s *ElementSet) Remove(value interface{}, t time.Time) error {
	return s.removals.add(value, t)
}

// Exists checks if an element is marked as present in the set
// Please note that this function provides liveness guarantees only.
func (s ElementSet) Exists(value interface{}) (bool, error) {
	addedAt, added, addedErr := s.additions.addedAt(value)
	if addedErr != nil {
		return false, addedErr
	}
	if !added {
		return false, nil
	}

	removed, removedErr := s.isRemoved(value, addedAt)
	if removedErr != nil {
		return false, removedErr
	}
	if !removed {
		return true, nil
	}
	return false, nil
}

func (s ElementSet) isRemoved(value interface{}, since time.Time) (bool, error) {
	removedAt, removed, err := s.removals.addedAt(value)
	if err != nil {
		return false, err
	}
	if !removed {
		return false, nil
	}
	if since.Before(removedAt) {
		return true, nil
	}
	return false, nil
}

// Get returns set content
// Please note that this function pvovides liveness guarantees only
func (s ElementSet) Get() ([]interface{}, error) {
	var result []interface{}

	err := s.additions.each(func(element interface{}, addedAt time.Time) error {
		removed, removedErr := s.isRemoved(element, addedAt)
		if removedErr != nil {
			return removedErr
		}
		if !removed {
			result = append(result, element)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
