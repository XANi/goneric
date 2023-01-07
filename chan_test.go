package goneric

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
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
	f1 := ctr{}
	chGen1 := GenChan(func() int { return f1.SleepyCounter(time.Millisecond * 100) })
	out1 := ChanToSliceNTimeout(chGen1, 10, time.Millisecond*350)
	assert.Equal(t, []int{1, 2, 3}, out1)
	f2 := ctr{}
	chGen2 := GenChan(func() int { return f2.SleepyCounter(time.Millisecond) })
	out2 := ChanToSliceNTimeout(chGen2, 10, time.Millisecond*350)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, out2)
}

func TestSliceToChan(t *testing.T) {
	data := []int{6, 5, 3, 8}
	ch := make(chan int, 1)
	SliceToChan(data, ch)
	out := ChanToSliceN(ch, len(data))
	assert.Equal(t, data, out)
	assert.NotPanics(t, func() { close(ch) }, "make sure out channel is open")
}

func TestSliceToChanClose(t *testing.T) {
	data := []int{6, 5, 3, 8}
	ch := make(chan int, 1)
	SliceToChanClose(data, ch)
	out := ChanToSlice(ch)
	assert.Equal(t, data, out)
	assert.Panics(t, func() { close(ch) }, "make sure out channel is closed")
}
