package goneric

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenSlice(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3, 4, 5},
		GenSlice(5, func(i int) int { return i + 1 }),
	)
}

func TestGenMap(t *testing.T) {
	assert.Equal(t, map[int]int{0: 1, 1: 2, 2: 3, 3: 4, 4: 5},
		GenMap(5, func(i int) (int, int) { return i, i + 1 }),
	)
}
