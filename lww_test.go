package lww

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func testingRedisElementSet(t *testing.T) *ElementSet {
	c := testingRedisClient(t)
	c.Del(strings.Join([]string{redisTestingKey, "_lww_additions"}, ""))
	c.Del(strings.Join([]string{redisTestingKey, "_lww_removals"}, ""))
	return NewRedisElementSet(redisTestingKey, c)
}

func testingSets(t *testing.T) []*ElementSet {
	sets := []*ElementSet{}
	sets = append(sets, NewMapElementSet())
	sets = append(sets, testingRedisElementSet(t))
	return sets
}

func TestLWWNew(t *testing.T) {
	for _, set := range testingSets(t) {
		assert.Equal(t, 0, len(set.Get()))
	}
}

func TestLWWAddRemoveAndExists(t *testing.T) {
	for _, set := range testingSets(t) {
		assert.Equal(t, false, set.Exists("Hello"))
		set.Add("Hello", time.Now())
		assert.Equal(t, true, set.Exists("Hello"))
		set.Remove("Hello", time.Now())
		assert.Equal(t, false, set.Exists("Hello"))
	}
}

func TestLWWAddRemoveAndGet(t *testing.T) {
	for _, set := range testingSets(t) {
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
}
