package goneric

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
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
		WorkerPool(in, out, func(i int) string {
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

}
