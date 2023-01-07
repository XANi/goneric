package goneric

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestChanGen(t *testing.T) {
	f := ctr{}
	ch := make(chan int, 1)
	ChanGen(f.Counter, ch)
	a := <-ch
	b := <-ch
	c := <-ch
	d := <-ch
	assert.Equal(t, []int{1, 2, 3, 4}, []int{a, b, c, d})
}

func TestChanGenN(t *testing.T) {
	ch := make(chan int, 1)
	ChanGenN(3, func(i int) int { return i + 1 }, ch)
	a := <-ch
	b := <-ch
	c := <-ch
	var d int
	select {
	case d = <-ch:
	case <-time.After(time.Millisecond * 20):
	}

	assert.Equal(t, []int{1, 2, 3, 0}, []int{a, b, c, d})
}

func TestChanGenNClose(t *testing.T) {
	ch := make(chan int, 1)
	ChanGenNClose(3, func(i int) int { return i + 1 }, ch)
	data := make([]int, 0)
	idx := 0
O:
	for {
		idx++
		select {
		case v := <-ch:
			data = append(data, v)
		case <-time.After(time.Millisecond * 20):
			break O
		}
		if idx > 3 {
			break
		}
	}

	assert.Equal(t, []int{1, 2, 3, 0}, data)
}

func TestChanGenCloser(t *testing.T) {
	t.Run("with closing output channel", func(t *testing.T) {
		f := ctr{}
		ch := make(chan int, 1)
		cl := ChanGenCloser(f.Counter, ch)
		_ = <-ch
		cl(true)
		// first one after channel drains
		_ = <-ch
		// next one for loop iteration to get to the exit check
		_ = <-ch
		// this one should be empty
		d := <-ch
		assert.NotEqual(t, 4, d)
		assert.Panics(t, func() { close(ch) }, "make sure out channel is closed")
	})
	t.Run("without closing output channel", func(t *testing.T) {
		f := ctr{}
		ch := make(chan int, 1)
		cl := ChanGenCloser(f.Counter, ch)
		_ = <-ch
		cl()
		// first one after channel drains
		_ = <-ch
		// next one for loop iteration to get to the exit check
		_ = <-ch
		// this one should be empty
		var d int
		select {
		case d = <-ch:
		case <-time.After(time.Millisecond * 50):
		}
		assert.NotEqual(t, 4, d)
		assert.NotPanics(t, func() { close(ch) }, "make sure out channel is not closed")
	})

}
