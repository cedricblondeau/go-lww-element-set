package lww

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLWWNew(t *testing.T) {
	set := NewElementSet()
	assert.Equal(t, 0, len(set.additions.elements))
	assert.Equal(t, 0, len(set.removals.elements))
}
