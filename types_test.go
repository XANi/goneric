package goneric

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTypes(t *testing.T) {
	assert.NotEmpty(t, ErrSkip{}.Error())
}
