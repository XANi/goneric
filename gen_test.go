package goneric

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
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

func TestGenChan(t *testing.T) {
	f := ctr{}
	ch := GenChan(f.Counter)
	a := <-ch
	b := <-ch
	c := <-ch
	d := <-ch
	assert.Equal(t, []int{1, 2, 3, 4}, []int{a, b, c, d})
}

func TestGenChanN(t *testing.T) {
	t.Run("chan open", func(t *testing.T) {
		ch := GenChanN(func(i int) int { return i + 1 }, 3)
		a := <-ch
		b := <-ch
		c := <-ch
		var d int
		select {
		case d = <-ch:
		case <-time.After(time.Millisecond * 20):
		}

		assert.Equal(t, []int{1, 2, 3, 0}, []int{a, b, c, d})
	})
	t.Run("chan close", func(t *testing.T) {
		ch := GenChanN(func(i int) int { return i + 1 }, 3, true)
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
	})
}

func TestGenChanCloser(t *testing.T) {
	t.Run("with closing output channel", func(t *testing.T) {
		f := ctr{}
		ch, cl := GenChanCloser(f.Counter)
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
		ch, cl := GenChanCloser(f.Counter)
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

func TestGenSliceToChan(t *testing.T) {
	data := []int{6, 5, 3, 8}
	t.Run("chan open", func(t *testing.T) {
		ch := GenSliceToChan(data)
		out := ChanToSliceN(ch, len(data))
		assert.Equal(t, data, out)
		assert.NotPanics(t, func() { close(ch) }, "make sure out channel is open")
	})
	t.Run("chan close", func(t *testing.T) {
		data := []int{6, 5, 3, 8}
		ch := GenSliceToChan(data, true)
		out := ChanToSlice(ch)
		assert.Equal(t, data, out)
		assert.Panics(t, func() { close(ch) }, "make sure out channel is closed")
	})
}
func ExampleGenSliceToChan() {
	// jobs to do
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	output := ChanToSlice( // make slice out of channel
		WorkerPoolBackground( // that we got out of worker
			GenSliceToChan(input, true), // that got fed input from slice via channel
			func(v int) float64 {
				// pretend we have some work
				time.Sleep(time.Millisecond*20 + time.Duration(rand.Int31n(20)))
				return float64(v) * 1.5
			},
			16, true, // in parallel
		))
	fmt.Printf("%+v->(1.5x)->%+v", Sum(input...), Sum(output...))
	// Output: 55->(1.5x)->82.5
}
