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
		elements, err := set.Get()
		assert.Nil(t, err)
		assert.Equal(t, 0, len(elements))
	}
}

func TestLWWAddRemoveAndExists(t *testing.T) {
	for _, set := range testingSets(t) {
		exists1, exists1Err := set.Exists("Hello")
		assert.Nil(t, exists1Err)
		assert.Equal(t, false, exists1)

		addErr := set.Add("Hello", time.Now())
		assert.Nil(t, addErr)

		exists2, exists2Err := set.Exists("Hello")
		assert.Nil(t, exists2Err)
		assert.Equal(t, true, exists2)

		removeErr := set.Remove("Hello", time.Now())
		assert.Nil(t, removeErr)

		exists3, exists3Err := set.Exists("Hello")
		assert.Nil(t, exists3Err)
		assert.Equal(t, false, exists3)
	}
}

func TestLWWAddRemoveAndGet(t *testing.T) {
	for _, set := range testingSets(t) {
		set.Add("Montreal", time.Now())
		set.Add("NYC", time.Now())
		set.Add("Toronto", time.Now())
		set.Add("Paris", time.Now())
		set.Remove("Montreal", time.Now())
		result, err := set.Get()
		assert.Nil(t, err)
		assert.Equal(t, 3, len(result))
		assert.Contains(t, result, "NYC")
		assert.Contains(t, result, "Toronto")
		assert.Contains(t, result, "Paris")
		assert.NotContains(t, result, "Montreal")
	}
}
