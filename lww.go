package lww

// ElementSet is a Last-Writer-Wins (LWW) Element Set data structure
// implemented using two TimedSet data structures.
type ElementSet struct {
	additions *TimedSet
	removals  *TimedSet
}

// NewElementSet returns an empty and ready-to-use LWW element set
func NewElementSet() *ElementSet {
	return &ElementSet{
		additions: NewTimedSet(),
		removals:  NewTimedSet(),
	}
}
