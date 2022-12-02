package goneric

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSum(t *testing.T) {
	assert.Equal(t, 3, Sum(1, 2))
	assert.Equal(t, 6, Sum(1, 2, 3))
	assert.EqualValues(t, 10, Sum(1.0, 2.0, 3.0, 4.0))
}
