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

func TestLWWExists(t *testing.T) {
	set := NewElementSet()
	assert.Equal(t, false, set.Exists("Hello"))
	set.Add("Hello", time.Now())
	assert.Equal(t, true, set.Exists("Hello"))
	set.Remove("Hello", time.Now())
	assert.Equal(t, false, set.Exists("Hello"))
}

func TestLWWGest(t *testing.T) {
	set := NewElementSet()
	set.Add("Montreal", time.Now())
	set.Add("NYC", time.Now())
	set.Add("Toronto", time.Now())
	set.Add("Paris", time.Now())
	set.Remove("Montreal", time.Now())
	result := set.Get()
	assert.Equal(t, 3, len(result))
	assert.Contains(t, result, "NYC")
	assert.Contains(t, result, "Toronto")
	assert.Contains(t, result, "Paris")
	assert.NotContains(t, result, "Montreal")
}
