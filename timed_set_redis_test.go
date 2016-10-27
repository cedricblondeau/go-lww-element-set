package lww

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	redis "gopkg.in/redis.v4"
)

func testingRedisClient(t *testing.T) *redis.Client {
	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := c.Ping().Result()
	if err != nil {
		t.Error("Cannot set up redis for tests", err)
	}
	return c
}

func testingRedisTimedSet(t *testing.T) *RedisTimedSet {
	c := testingRedisClient(t)
	c.Del(redisTestingKey)
	return NewRedisTimedSet(redisTestingKey, c)
}

func TestTimeToFloat(t *testing.T) {
	time, _ := time.Parse(millisecondDateForm, "2016-Dec-24 23:59:59.998")
	f := timeToFloat(time)
	assert.Equal(t, 1.482623999998e+15, f)
}

func TestFloatToTime(t *testing.T) {
	expected, _ := time.Parse(millisecondDateForm, "2016-Dec-24 23:59:59.998")
	f := 1.482623999998e+15
	time := floatToTime(f).UTC()
	assert.Equal(t, expected, time)
}

func TestTimedSetRedisAdd(t *testing.T) {
	s := testingRedisTimedSet(t)
	s.Add("Raptors", time.Now())
	res := s.client.ZCard(redisTestingKey)
	assert.Equal(t, int64(1), res.Val())
}

func TestTimedSetRedisAddedAt(t *testing.T) {
	s := testingRedisTimedSet(t)
	now := time.Now()
	s.Add("Giants", now)
	addedAt, ok := s.AddedAt("Giants")
	assert.Equal(t, true, ok)
	expected := floatToTime(timeToFloat(now))
	assert.Equal(t, expected, addedAt)
}

func TestTimedSetRedisEach(t *testing.T) {
	s := testingRedisTimedSet(t)
	s.Add("Koala", time.Now())
	s.Add("Cat", time.Now())
	s.Add("Dog", time.Now())

	result := []string{}
	err := s.Each(func(element interface{}, addedAt time.Time) {
		result = append(result, element.(string))
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, []string{"Koala", "Cat", "Dog"}, result)
}

func TestTimedSetMapRedisSameElementWithGreaterTimestamp(t *testing.T) {
	oct24, _ := time.Parse(shortDateForm, "2016-Oct-24")
	oct25, _ := time.Parse(shortDateForm, "2016-Oct-25")

	s := testingRedisTimedSet(t)
	s.Add("Hi!", oct24)
	addedAt, _ := s.AddedAt("Hi!")
	assert.Equal(t, floatToTime(timeToFloat(oct24)), addedAt)

	s.Add("Hi!", oct25)
	addedAt, _ = s.AddedAt("Hi!")
	assert.Equal(t, floatToTime(timeToFloat(oct25)), addedAt)
}

func TestTimedSetMapRedisSameElementWithLesserTimestamp(t *testing.T) {
	oct23, _ := time.Parse(shortDateForm, "2016-Oct-23")
	oct24, _ := time.Parse(shortDateForm, "2016-Oct-24")

	s := testingRedisTimedSet(t)
	s.Add("Hi!", oct24)
	addedAt, _ := s.AddedAt("Hi!")
	assert.Equal(t, floatToTime(timeToFloat(oct24)), addedAt)

	s.Add("Hi!", oct23)
	addedAt, _ = s.AddedAt("Hi!")
	assert.Equal(t, floatToTime(timeToFloat(oct24)), addedAt)
}
