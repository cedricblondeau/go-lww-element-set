package lww

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimedSetNew(t *testing.T) {
	s := NewTimedSet()
	assert.Equal(t, 0, len(s.elements))
}
