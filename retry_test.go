package goneric

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

type fa struct {
	i int
	sync.Mutex
}

func (f *fa) OkAfter(i int) (int, error) {
	f.Lock()
	defer f.Unlock()
	f.i++
	if f.i >= i {
		return f.i, nil
	} else {
		return f.i, errors.New("fail")
	}
}

func TestTools(t *testing.T) {
	a := fa{}
	r1, e1 := a.OkAfter(3)
	assert.Equal(t, 1, r1)
	assert.Error(t, e1)
	r2, e2 := a.OkAfter(3)
	assert.Equal(t, 2, r2)
	assert.Error(t, e2)
	r3, e3 := a.OkAfter(3)
	assert.Equal(t, 3, r3)
	assert.NoError(t, e3)
}

func TestRetry(t *testing.T) {
	a1 := fa{}
	out1, err1 := Retry(3, func() (int, error) { return a1.OkAfter(3) })
	assert.NoError(t, err1)
	assert.Equal(t, 3, out1)

	a2 := fa{}
	out, err2 := Retry(3, func() (int, error) { return a2.OkAfter(2) })
	assert.NoError(t, err2)
	assert.Equal(t, 2, out)

	a3 := fa{}
	out, err3 := Retry(3, func() (int, error) { return a3.OkAfter(4) })
	assert.Error(t, err3)
	assert.Equal(t, 3, out)
}

func TestRetryExp(t *testing.T) {
	a1 := fa{}
	out1, err1 := RetryAfter(
		time.Millisecond*5,
		time.Millisecond*50,
		time.Millisecond*500,
		func() (int, error) {
			return a1.OkAfter(3)
		})
	assert.NoError(t, err1)
	assert.Equal(t, 3, out1)

	a2 := fa{}
	_, err2 := RetryAfter(
		time.Millisecond*30,
		time.Millisecond*40,
		time.Millisecond*50,
		func() (int, error) { return a2.OkAfter(9) })
	assert.Error(t, err2)
}

func TestTry(t *testing.T) {
	funcList := []func() (int, error){
		func() (int, error) { return 1, errors.New("1") },
		func() (int, error) { return 2, errors.New("2") },
		func() (int, error) { return 3, errors.New("3") },
		func() (int, error) { return 4, nil },
	}
	out, err := Try(funcList...)
	assert.Equal(t, 4, out)
	assert.NoError(t, err)
	out, err = Try(funcList[:2]...)
	assert.Equal(t, 2, out)
	assert.Error(t, err)

}
