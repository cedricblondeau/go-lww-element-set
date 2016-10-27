# Last-Writer-Wins (LWW) Element Set [![GoDoc](https://godoc.org/github.com/cedricblondeau/go-lww-element-set?status.svg)](https://godoc.org/github.com/cedricblondeau/go-lww-element-set) [![Build Status](https://api.travis-ci.org/cedricblondeau/go-lww-element-set.svg?branch=develop)](https://travis-ci.org/cedricblondeau/go-lww-element-set)

> A LWW CRDT implementation with Redis support in Go.

[![forthebadge](http://forthebadge.com/images/badges/as-seen-on-tv.svg)](http://forthebadge.com)

## A CR-What?!

A conflict-free replicated data type (CRDT) is a type of data structure that is 
used to achieve **strong eventual consistency** and monotonicity (ie, there are no rollbacks) 
across a set of nodes in a **distributed system**.

## Package

This package focuses on Last-Writer-Wins (LWW) Element Set data structure 
that uses timestamped adds and removes.
It provides 2 implementations (Go maps and Redis).

Unlike [Roshi](https://github.com/soundcloud/roshi), this package **does not** provide replication,
sharding, garbage collection or REST-ish HTTP interface.

## Public API

This package (`lww`) exposes 3 different constructors. Each constructor returns an `ElementSet`.
See [GoDoc reference](https://godoc.org/github.com/cedricblondeau/go-lww-element-set) for details.

## Usage

#### Go maps backed LWW element set

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

#### Redis backed LWW element set

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

You can also use `NewRedisElementSetWithCustomMarshalling()` construtor 
to specify custom marshalling/unmarshalling functions.

## Implementation details

[ElementSet](lww.go) defines the LWW Element Set logic.
An [ElementSet](lww.go) is backed by two underlying 
timestamped sets ([timedSet](timed_set.go)).

[timedSet](timed_set.go) is an interface that defines 
a set where added element are associated with timestamps.
This package provides two implementations of this interface:
- mapTimedSet (using Go maps)
- redisTimedSet (using Redis)

#### Go maps implementation

Go maps implementation uses two [mapTimedSet](timed_set_map.go).
A [mapTimedSet](timed_set_map.go) is a [timedSet](timed_set.go) backed 
with a [Go map](https://blog.golang.org/go-maps-in-action).
Because Go maps are [not safe for concurrent use](https://golang.org/doc/faq#atomic_maps), 
mutual exclusion is used.

#### Redis implementation

Redis implementation uses two [redisTimedSet](timed_set_redis.go).
A [redisTimedSet](timed_set_redis.go) is a [timedSet](timed_set.go) backed 
with a [Redis sorted set](http://redis.io/topics/data-types#sorted-sets).
This implementation uses a Redis script ([which is transactional by definition and by extension atomic](http://redis.io/topics/transactions#redis-scripting-and-transactions)) 
to add elements.

## Run tests

To run tests with race detector enabled:

```go
go test -race
```

## References

- [A comprehensive study of Convergent and Commutative Replicated Data Types by Marc Shapiro](https://hal.inria.fr/file/index/docid/555588/filename/techreport.pdf)
- [Summary of CRDTs by Vaughn Vernon](https://vaughnvernon.co/?p=1012)
- [Roshi: a CRDT system for timestamped events by Peter Bourgon](https://developers.soundcloud.com/blog/roshi-a-crdt-system-for-timestamped-events)
- [Consistency without consensus in production systems by Peter Bourgon](https://www.youtube.com/watch?list=UU_QIfHvN9auy2CoOdSfMWDw&v=em9zLzM8O7c)
- [CRDT notes by Paul Frazee](https://github.com/pfrazee/crdt_notes)
