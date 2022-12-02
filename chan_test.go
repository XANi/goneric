package goneric

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"testing"
	"time"
)

type ctr struct {
	i int
}

func (c *ctr) Counter() int {
	c.i++
	return c.i
}

func TestChanGen(t *testing.T) {
	f := ctr{}
	ch := ChanGen(f.Counter)
	a := <-ch
	b := <-ch
	c := <-ch
	d := <-ch
	assert.Equal(t, []int{1, 2, 3, 4}, []int{a, b, c, d})
}

func TestChanGenCloser(t *testing.T) {
	f := ctr{}
	ch, cl := ChanGenCloser(f.Counter)
	_ = <-ch
	cl()
	// first one after channel drains
	_ = <-ch
	// next one for loop iteration to get to the exit check
	_ = <-ch
	// this one should be empty
	d := <-ch
	assert.NotEqual(t, 4, d)
	assert.Panics(t, func() { close(ch) }, "make sure out channel is closed")
}

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
	chGen := ChanGen(f.Counter)
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

func TestSliceToChan(t *testing.T) {
	data := []int{6, 5, 3, 8}
	ch := SliceToChan(data)
	out := ChanToSliceN(ch, len(data))
	assert.Equal(t, data, out)
	assert.NotPanics(t, func() { close(ch) }, "make sure out channel is open")
}

func TestSliceToChanClose(t *testing.T) {
	data := []int{6, 5, 3, 8}
	ch := SliceToChanClose(data)
	out := ChanToSlice(ch)
	assert.Equal(t, data, out)
	assert.Panics(t, func() { close(ch) }, "make sure out channel is closed")
}

func ExampleSliceToChanClose() {
	// jobs to do
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	output := ChanToSlice( // make slice out of channel
		WorkerPoolBackgroundClose( // that we got out of worker
			SliceToChanClose(input), // that got fed input from slice via channel
			func(v int) float64 {
				// pretend we have some work
				time.Sleep(time.Millisecond*20 + time.Duration(rand.Int31n(20)))
				return float64(v) * 1.5
			},
			16, // in parallel
		))
	fmt.Printf("%+v->(1.5x)->%+v", Sum(input...), Sum(output...))
	// Output: 55->(1.5x)->82.5
}
