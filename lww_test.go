package lww

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLWWNew(t *testing.T) {
	set := NewElementSet()
	assert.Equal(t, 0, len(set.additions.elements))
	assert.Equal(t, 0, len(set.removals.elements))
}

func TestLWWAdd(t *testing.T) {
	set := NewElementSet()
	set.Add("Hello", time.Now())
	assert.Equal(t, 1, len(set.additions.elements))
}

func TestLWWRemove(t *testing.T) {
	set := NewElementSet()
	set.Remove("Hello", time.Now())
	assert.Equal(t, 1, len(set.removals.elements))
}
