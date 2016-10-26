package lww

import (
	"time"

	redis "gopkg.in/redis.v4"
)

// This Lua script adds or updates an element in the sorted set
// if one of the two following condition is met:
//
// - Given element (ARGV[1]) does not exists yet
// - Given element (ARGV[1]) already exists
//   but with a lesser timestamp than the given one (ARGV[2])
const redisAddScript string = `
local c = tonumber(redis.call('ZSCORE', KEYS[1], ARGV[2]))
if not c or tonumber(ARGV[1]) > c then
	redis.call('ZADD', KEYS[1], ARGV[1], ARGV[2])
	return 1
else
	return 0
end
`

// RedisTimedSet is a TimedSet implementation that uses Redis.
type RedisTimedSet struct {
	client    *redis.Client
	addScript *redis.Script
	key       string
	marshal   func(value interface{}) string
	unmarshal func(string) interface{}
}

// NewRedisTimedSet returns an empty and ready-to-use Redis-backed TimedSet data structure.
func NewRedisTimedSet(key string, client *redis.Client) *RedisTimedSet {
	s := &RedisTimedSet{
		client:    client,
		addScript: redis.NewScript(redisAddScript),
		key:       key,
	}
	s.marshal = func(value interface{}) string {
		return value.(string)
	}
	s.unmarshal = func(value string) interface{} {
		return value
	}
	return s
}

/**
 * ZSET uses IEEE 754 64-bit numbers to sort the elements (score).
 * Acceptable values are: -(2^53) >= +(2^53).
 * This will limit the timestamp precision to 1 millisecond.
 */
func timeToFloat(t time.Time) float64 {
	f := float64(t.Round(time.Microsecond).UnixNano() / 1000)
	return f
}

func floatToTime(score float64) time.Time {
	return time.Unix(0, 0).Add(time.Duration(score) * time.Microsecond)
}

// Add runs a Redis script. Redis scripts are transactional by definition
// and by extension atomic.
func (s RedisTimedSet) Add(value interface{}, t time.Time) {
	s.addScript.Run(s.client, []string{s.key}, timeToFloat(t), s.marshal(value)).Result()
}

// AddedAt returns the timestamp of a given element if it exists.
func (s RedisTimedSet) AddedAt(value interface{}) (time.Time, bool) {
	score, err := s.client.ZScore(s.key, s.marshal(value)).Result()
	if err != nil {
		return time.Time{}, false
	}
	return floatToTime(score), true
}

// SetMarshal allows to specify a specific marshalling function
func (s *RedisTimedSet) SetMarshal(f func(interface{}) string) {
	s.marshal = f
}

// SetUnmarshal allows to specify a specific unmarshalling function
func (s *RedisTimedSet) SetUnmarshal(f func(value string) interface{}) {
	s.unmarshal = f
}

// Each traverses the items in the TimedSet, calling the provided function
// for each element/timestamp association.
func (s RedisTimedSet) Each(f func(element interface{}, addedAt time.Time)) error {
	r, err := s.client.ZRangeWithScores(s.key, 0, -1).Result()
	if err != nil {
		return err
	}
	for _, v := range r {
		f(v.Member, floatToTime(v.Score))
	}
	return nil
}
