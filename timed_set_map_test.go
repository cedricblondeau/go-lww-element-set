package lww

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const shortDateForm = "2006-Jan-02"

func TestTimedSetMapNew(t *testing.T) {
	s := NewTimedSetMap()
	assert.Equal(t, 0, len(s.elements))
}

func TestTimedSetMapAdd(t *testing.T) {
	s := NewTimedSetMap()
	s.Add("Hi!", time.Now())
	assert.Equal(t, 1, len(s.elements))
}

func TestTimedSetMapAddedAt(t *testing.T) {
	s := NewTimedSetMap()
	_, ok := s.AddedAt("Nothing")
	assert.Equal(t, false, ok)

	now := time.Now()
	s.Add("Hello world!", now)
	addedAt, ok := s.AddedAt("Hello world!")
	assert.Equal(t, true, ok)
	assert.Equal(t, now, addedAt)
}

func TestTimedSetAddSameElementWithGreaterTimestamp(t *testing.T) {
	oct24, _ := time.Parse(shortDateForm, "2016-Oct-24")
	oct25, _ := time.Parse(shortDateForm, "2016-Oct-25")

	s := NewTimedSetMap()
	s.Add("Hi!", oct24)
	addedAt, _ := s.AddedAt("Hi!")
	assert.Equal(t, oct24, addedAt)

	s.Add("Hi!", oct25)
	addedAt, _ = s.AddedAt("Hi!")
	assert.Equal(t, oct25, addedAt)
}

func TestTimedSetAddSameElementWithLesserTimestamp(t *testing.T) {
	oct23, _ := time.Parse(shortDateForm, "2016-Oct-23")
	oct24, _ := time.Parse(shortDateForm, "2016-Oct-24")

	s := NewTimedSetMap()
	s.Add("Hi!", oct24)
	addedAt, _ := s.AddedAt("Hi!")
	assert.Equal(t, oct24, addedAt)

	s.Add("Hi!", oct23)
	addedAt, _ = s.AddedAt("Hi!")
	assert.Equal(t, oct24, addedAt)
}

func TestTimedSetConcurrentAdd(t *testing.T) {
	oct23, _ := time.Parse(shortDateForm, "2016-Oct-23")
	oct24, _ := time.Parse(shortDateForm, "2016-Oct-24")
	s := NewTimedSetMap()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		s.Add("Hi!", oct24)
	}()
	go func() {
		defer wg.Done()
		s.Add("Hi!", oct23)
	}()
	wg.Wait()

	addedAt, ok := s.AddedAt("Hi!")
	assert.Equal(t, true, ok)
	assert.Equal(t, oct24, addedAt)
}