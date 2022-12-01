package goneric

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilterMap(t *testing.T) {
	data := map[string]int{"a": 1, "b": 9, "c": 3, "d": 2, "e": 10, "f": 5, "g": 6, "h": 7}
	out := FilterMap(func(k string, v int) bool {
		if v > 8 {
			return false
		}
		return true
	}, data)
	assert.Equal(t, map[string]int{"a": 1, "c": 3, "d": 2, "f": 5, "g": 6, "h": 7}, out)
}

func TestFilterSlice(t *testing.T) {
	data := []int{1, 9, 3, 2, 10, 5, 6, 7}
	out := FilterSlice(func(idx int, v int) bool {
		if v > 8 {
			return false
		}
		return true
	}, data)
	assert.Equal(t, []int{1, 3, 2, 5, 6, 7}, out)
}
