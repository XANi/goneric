package goneric

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	in := make(chan int, 1)
	out := make(chan string, 2)
	go func() {
		for i := 0; i < 32; i++ {
			in <- i
		}
		close(in)
	}()
	go func() {
		WorkerPoolClose(in, out, func(i int) string {
			return strconv.Itoa(i)
		}, 4)
	}()

	outSlice := []string{}
	for o := range out {
		outSlice = append(outSlice, o)
	}
	assert.True(t, CompareSliceSet(
		[]string{"0", "1", "2", "4", "5", "6", "7", "3", "8", "9", "10", "11", "12", "13", "16", "14", "17", "18", "19", "15", "20", "21", "22", "25", "26", "27", "28", "23", "24", "29", "30", "31"},
		outSlice))
	assert.Panics(t, func() {
		WorkerPool(in, out, func(i int) string {
			return strconv.Itoa(i)
		}, 0)
	})

}

func TestWorkerPoolBackground(t *testing.T) {
	in := make(chan int, 1)
	go func() {
		for i := 0; i < 32; i++ {
			in <- i
		}
		close(in)
	}()
	out := WorkerPoolBackground(in, func(i int) string {
		return strconv.Itoa(i)
	}, 4)

	outSlice := ChanToSliceN(out, 32)
	assert.True(t, CompareSliceSet(
		[]string{"0", "1", "2", "4", "5", "6", "7", "3", "8", "9", "10", "11", "12", "13", "16", "14", "17", "18", "19", "15", "20", "21", "22", "25", "26", "27", "28", "23", "24", "29", "30", "31"},
		outSlice))
	assert.Panics(t, func() {
		WorkerPoolBackground(in, func(i int) string {
			return strconv.Itoa(i)
		}, 0)
	})

}

func TestWorkerPoolBackgroundClose(t *testing.T) {
	in := make(chan int, 1)
	go func() {
		for i := 0; i < 32; i++ {
			in <- i
		}
		close(in)
	}()
	out := WorkerPoolBackgroundClose(in, func(i int) string {
		return strconv.Itoa(i)
	}, 4)

	outSlice := ChanToSlice(out)
	assert.True(t, CompareSliceSet(
		[]string{"0", "1", "2", "4", "5", "6", "7", "3", "8", "9", "10", "11", "12", "13", "16", "14", "17", "18", "19", "15", "20", "21", "22", "25", "26", "27", "28", "23", "24", "29", "30", "31"},
		outSlice))
	assert.Panics(t, func() {
		WorkerPoolBackgroundClose(in, func(i int) string {
			return strconv.Itoa(i)
		}, 0)
	})

}

func TestWorkerPoolFinisher(t *testing.T) {
	in := make(chan int, 1)
	out := make(chan string, 2)
	go func() {
		for i := 0; i < 32; i++ {
			in <- i
		}
		close(in)
	}()
	finisher := WorkerPoolFinisher(in, out, func(i int) string {
		return strconv.Itoa(i)
	}, 4)

	outSlice := []string{}
	for o := range out {
		outSlice = append(outSlice, o)
	}
	assert.True(t, CompareSliceSet(
		[]string{"0", "1", "2", "4", "5", "6", "7", "3", "8", "9", "10", "11", "12", "13", "16", "14", "17", "18", "19", "15", "20", "21", "22", "25", "26", "27", "28", "23", "24", "29", "30", "31"},
		outSlice))
	assert.True(t, <-finisher)
	assert.Panics(t, func() {
		WorkerPool(in, out, func(i int) string {
			return strconv.Itoa(i)
		}, 0)
	})
	assert.Panics(t, func() { close(out) })
}

func TestWorkerPoolDrain(t *testing.T) {
	c := make(chan func(), 1)
	end := make(chan bool, 1)
	finish := WorkerPoolDrain(func(f func()) { f() }, 4, c)
	var a1, a2 int
	c <- func() { a1 = 1 }
	c <- func() { a2 = 2 }
	c <- func() { end <- true }
	close(c)
	<-finish
	time.Sleep(time.Millisecond * 10)
	assert.True(t, <-end)
	assert.Equal(t, 1, a1)
	assert.Equal(t, 2, a2)
}

func TestWorkerPoolAsync(t *testing.T) {
	async, finish := WorkerPoolAsync(func(i int) string { return strconv.Itoa(i) }, 2)
	a1 := async(1)
	a2 := async(2)
	a3 := async(3)
	a4 := async(4)
	finish()
	assert.Equal(t, "1", <-a1)
	assert.Equal(t, "2", <-a2)
	assert.Equal(t, "3", <-a3)
	assert.Equal(t, "4", <-a4)

}

func ExampleWorkerPoolAsync() {
	// make our worker with function to mangle data
	async, finish := WorkerPoolAsync(
		func(i int) string { return strconv.Itoa(i) },
		2)
	defer finish() // close the pool once we stop using it
	// queue some jobs
	job1 := async(1)
	job2 := async(2)
	job3 := async(3)
	job4 := async(4)
	// all of those jobs are running in background at this point
	// now get results
	fmt.Printf("%s %s %s %s",
		<-job1,
		<-job3,
		<-job2,
		<-job4,
	)
	//Output: 1 3 2 4
}
