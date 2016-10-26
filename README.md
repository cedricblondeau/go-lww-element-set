# Last-Writer-Wins (LWW) Element Set

> A LWW CRDT implementation in Go.

## CR-What?!

A conflict-free replicated data type (CRDT) is a type of data structure that is 
used to achieve strong eventual consistency and monotonicity (ie, there are no rollbacks) 
across a set of nodes in a distributed system.

This package focuses on Last-Writer-Wins (LWW) Element Set 
data structure that uses timestamped adds and removes.

## Implementation

This implementation ([ElementSet](lww.go)) uses 
two separate underlying timestamped sets ([TimedSet](timed_set.go)).

### TimedSetMap (using Go maps)

A [TimedSetMap](timed_set_map.go) is a [TimedSet](timed_set.go) backed 
with a [Go map](https://blog.golang.org/go-maps-in-action).
Because Go maps are [not safe for concurrent use](https://golang.org/doc/faq#atomic_maps), 
mutual exclusion is used.

## Usage

```go
import (
	"time"
	"github.com/cedricblondeau/go-lww-element-set"
)

lww := lww.NewElementSet()
lww.Add("Hello", time.Now())
lww.Add("Hi!", time.Now())
lww.Remove("Hello", time.Now())
elements := lww.Get() // ["Hi!"]
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
