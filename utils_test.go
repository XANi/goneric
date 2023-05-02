package goneric

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"testing"
)

func TestMust(t *testing.T) {
	v := make([]byte, 4)
	i := Must(rand.Read(v))
	assert.Equal(t, 4, i)

	assert.Panics(t, func() {
		_ = Must(strconv.Atoi("cat"))
	})
}

func TestIgnoreErr(t *testing.T) {
	v := make([]byte, 4)
	i := IgnoreErr(rand.Read(v))
	assert.Equal(t, 4, i)

	assert.NotPanics(t, func() {
		_ = IgnoreErr(strconv.Atoi("cat"))
	})
	assert.Equal(t, 0, IgnoreErr(strconv.Atoi("cat")))
}
