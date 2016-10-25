package lww

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimedSetNew(t *testing.T) {
	s := NewTimedSet()
	assert.Equal(t, 0, len(s.elements))
}

func TestTimedSetAdd(t *testing.T) {
	s := NewTimedSet()
	s.Add("Hi!", time.Now())
	assert.Equal(t, 1, len(s.elements))
}
