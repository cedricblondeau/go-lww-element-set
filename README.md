# Last-Writer-Wins (LWW) Element Set

> A LWW CRDT implementation in Go.

## CR-What?!

A conflict-free replicated data type (CRDT) is a type of data structure that is 
used to achieve strong eventual consistency and monotonicity (ie, there are no rollbacks) 
across a set of nodes in a distributed system.

This package focuses on Last-Writer-Wins (LWW) Element Set 
data structure that uses timestamped adds and removes.

## Implementations

[ElementSet](lww.go) defines the LWW Element Set logic.
An [ElementSet](lww.go) is backed by two underlying 
timestamped sets ([TimedSet](timed_set.go)).

[TimedSet](timed_set.go) is an interface that defines 
a set where added element are associated with timestamps.
This package provides two implementations of this interface:
- MapTimedSet (using Go maps)
- RedisTimedSet (using Redis)

#### Map-backed implementation

`NewMapElementSet()` returns a LWW backed with two [MapTimedSet](timed_set_map.go).

A [MapTimedSet](timed_set_map.go) is a [TimedSet](timed_set.go) backed 
with a [Go map](https://blog.golang.org/go-maps-in-action).
Because Go maps are [not safe for concurrent use](https://golang.org/doc/faq#atomic_maps), 
mutual exclusion is used.

```go
import (
	"time"
	"github.com/cedricblondeau/go-lww-element-set"
)

ms := lww.NewMapElementSet()
ms.Add("Hello", time.Now())
ms.Add("Hi!", time.Now())
ms.Remove("Hello", time.Now())
ms.Get() // ["Hi!"]
```

#### Redis-backed implementation

`NewRedisElementSet()` returns a LWW backed with two [RedisTimedSet](timed_set_redis.go).

A [RedisTimedSet](timed_set_redis.go) is a [TimedSet](timed_set.go) backed 
with a [Redis sorted set](http://redis.io/topics/data-types#sorted-sets).
This implementation uses a Redis script ([which is transactional by definition and by extension atomic](http://redis.io/topics/transactions#redis-scripting-and-transactions)) 
to add elements.

```go
import (
	"time"

	"github.com/cedricblondeau/go-lww-element-set"
	redis "gopkg.in/redis.v4"
)

rc := redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})
rs := lww.NewRedisElementSet("shopping_cart", rc)
rs.Add("Product #1", time.Now())
rs.Add("Product #2", time.Now())
rs.Remove("Product #1", time.Now())
rs.Get() // ["Product #2"]
```

## Run tests

To run tests with race detector enabled:

```go
go test -race
```

## References

- https://vaughnvernon.co/?p=1012
- https://developers.soundcloud.com/blog/roshi-a-crdt-system-for-timestamped-events
- https://www.youtube.com/watch?list=UU_QIfHvN9auy2CoOdSfMWDw&v=em9zLzM8O7c
- http://blog.plasmaconduit.com/crdts-distributed-semilattices/
