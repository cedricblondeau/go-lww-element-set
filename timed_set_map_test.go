package lww

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimedSetMapNew(t *testing.T) {
	s := newMapTimedSet()
	assert.Equal(t, 0, len(s.elements))
}

func TestTimedSetMapAdd(t *testing.T) {
	s := newMapTimedSet()
	s.add("Hi!", time.Now())
	assert.Equal(t, 1, len(s.elements))
}

func TestTimedSetMapAddedAt(t *testing.T) {
	s := newMapTimedSet()
	_, ok := s.addedAt("Nothing")
	assert.Equal(t, false, ok)

	now := time.Now()
	s.add("Hello world!", now)
	addedAt, ok := s.addedAt("Hello world!")
	assert.Equal(t, true, ok)
	assert.Equal(t, now, addedAt)
}

func TestTimedSetMapAddSameElementWithGreaterTimestamp(t *testing.T) {
	oct24, _ := time.Parse(shortDateForm, "2016-Oct-24")
	oct25, _ := time.Parse(shortDateForm, "2016-Oct-25")

	s := newMapTimedSet()
	s.add("Hi!", oct24)
	addedAt, _ := s.addedAt("Hi!")
	assert.Equal(t, oct24, addedAt)

	s.add("Hi!", oct25)
	addedAt, _ = s.addedAt("Hi!")
	assert.Equal(t, oct25, addedAt)
}

func TestTimedSetMapAddSameElementWithLesserTimestamp(t *testing.T) {
	oct23, _ := time.Parse(shortDateForm, "2016-Oct-23")
	oct24, _ := time.Parse(shortDateForm, "2016-Oct-24")

	s := newMapTimedSet()
	s.add("Hi!", oct24)
	addedAt, _ := s.addedAt("Hi!")
	assert.Equal(t, oct24, addedAt)

	s.add("Hi!", oct23)
	addedAt, _ = s.addedAt("Hi!")
	assert.Equal(t, oct24, addedAt)
}

func TestTimedSetConcurrentAdd(t *testing.T) {
	oct23, _ := time.Parse(shortDateForm, "2016-Oct-23")
	oct24, _ := time.Parse(shortDateForm, "2016-Oct-24")
	s := newMapTimedSet()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		s.add("Hi!", oct24)
	}()
	go func() {
		defer wg.Done()
		s.add("Hi!", oct23)
	}()
	wg.Wait()

	addedAt, ok := s.addedAt("Hi!")
	assert.Equal(t, true, ok)
	assert.Equal(t, oct24, addedAt)
}
