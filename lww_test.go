package lww

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLWWNew(t *testing.T) {
	set := NewElementSet()
	assert.Equal(t, 0, len(set.Get()))
}

func TestLWWAddRemoveAndExists(t *testing.T) {
	set := NewElementSet()
	assert.Equal(t, false, set.Exists("Hello"))
	set.Add("Hello", time.Now())
	assert.Equal(t, true, set.Exists("Hello"))
	set.Remove("Hello", time.Now())
	assert.Equal(t, false, set.Exists("Hello"))
}

func TestLWWAddRemoveAndGet(t *testing.T) {
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
