package goneric

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

func TestChanToSlice(t *testing.T) {
	ch := make(chan int, 5)
	var sl []int
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		sl = ChanToSlice(ch)
		wg.Done()
	}()
	ch <- 1
	ch <- 2
	ch <- 3
	ch <- 4
	close(ch)
	wg.Wait()
	assert.Equal(t, []int{1, 2, 3, 4}, sl)
}

func TestChanToSliceN(t *testing.T) {
	f := ctr{}
	chGen := GenChan(f.Counter)
	sl := ChanToSliceN(chGen, 10)
	assert.Equal(t, 55, Sum(sl...))
	assert.Len(t, sl, 10)
	ch2 := make(chan int, 5)
	var sl2 []int
	wg2 := sync.WaitGroup{}
	wg2.Add(1)
	go func() {
		sl2 = ChanToSliceN(ch2, 2)
		wg2.Done()
	}()
	ch2 <- 1
	close(ch2)
	wg2.Wait()
	assert.Equal(t, []int{1}, sl2)
}

func TestChanToSliceNTimeout(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// timeout hits first: elements arrive at 100/200/300ms, timeout at 350ms
		// GenChanN (unlike GenChan) terminates, which synctest requires of bubble goroutines
		f1 := ctr{}
		ch1 := GenChanN(func(int) int { return f1.SleepyCounter(time.Millisecond * 100) }, 10, true)
		out1 := ChanToSliceNTimeout(ch1, 10, time.Millisecond*350)
		assert.Equal(t, []int{1, 2, 3}, out1)
		// drain so the generator goroutine can exit before the bubble ends
		for range ch1 {
		}
		// n hits first: all 10 elements arrive within 10ms
		f2 := ctr{}
		ch2 := GenChanN(func(int) int { return f2.SleepyCounter(time.Millisecond) }, 10, true)
		out2 := ChanToSliceNTimeout(ch2, 10, time.Millisecond*350)
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, out2)
		// channel closed before n or timeout: returns what was received, immediately
		f3 := ctr{}
		ch3 := GenChanN(func(int) int { return f3.SleepyCounter(time.Millisecond) }, 2, true)
		out3 := ChanToSliceNTimeout(ch3, 10, time.Hour)
		assert.Equal(t, []int{1, 2}, out3)
	})
}

func TestSliceToChan(t *testing.T) {
	data := []int{6, 5, 3, 8}
	t.Run("chan open", func(t *testing.T) {
		ch := make(chan int, 1)
		SliceToChan(data, ch)
		out := ChanToSliceN(ch, len(data))
		assert.Equal(t, data, out)
		assert.NotPanics(t, func() { close(ch) }, "make sure out channel is open")
	})
	t.Run("chan close", func(t *testing.T) {
		ch := make(chan int, 1)
		SliceToChan(data, ch, true)
		out := ChanToSlice(ch)
		assert.Equal(t, data, out)
		assert.Panics(t, func() { close(ch) }, "make sure out channel is closed")
	})
}
