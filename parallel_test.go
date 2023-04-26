package goneric

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestParallelMap(t *testing.T) {
	mappedData := ParallelMap(func(v string) int {
		time.Sleep(time.Millisecond * time.Duration(rand.Int31n(10)))
		i, _ := strconv.Atoi(v)
		return i
	},
		3,
		"9", "1", "2", "4", "5", "11")
	assert.Equal(t, []int{9, 1, 2, 4, 5, 11}, mappedData)
}

func TestParallelMapSlice(t *testing.T) {
	mappedData := ParallelMapSlice(func(v string) int {
		time.Sleep(time.Millisecond * time.Duration(rand.Int31n(10)))
		i, _ := strconv.Atoi(v)
		return i
	},
		3,
		[]string{"1", "3", "2", "7", "9", "12"})
	assert.Equal(t, []int{1, 3, 2, 7, 9, 12}, mappedData)
}

func ExampleParallelMapSlice() {
	data := []string{"1", "3", "2", "7", "9", "12"}
	mappedData := ParallelMapSlice(
		// function used to map
		func(v string) int {
			time.Sleep(time.Millisecond * time.Duration(rand.Int31n(10)))
			i, _ := strconv.Atoi(v)
			return i
		},
		3, // run it at least this many times
		data)
	// order stays
	fmt.Printf("%T%+v\n%T%+v", data, data, mappedData, mappedData)
	//Output: []string[1 3 2 7 9 12]
	//[]int[1 3 2 7 9 12]
}

func TestParallelSliceMapChannel(t *testing.T) {
	data := []string{"1", "3", "2", "7", "9", "12"}
	ch := ParallelMapSliceChan(func(s string) int {
		i, _ := strconv.Atoi(s)
		return i
	}, 3, data)
	out := []int{}
	for o := range ch {
		out = append(out, o)
	}
	assert.True(t, CompareSliceSet([]int{1, 2, 3, 7, 9, 12}, out), out)
}

func TestParallelSliceMapChannelFinisher(t *testing.T) {
	data := []string{"1", "3", "2", "7", "9", "12"}
	ch, finish := ParallelMapSliceChanFinisher(func(s string) int {
		i, _ := strconv.Atoi(s)
		return i
	}, 3, data)
	out := []int{}
	for o := range ch {
		out = append(out, o)
	}
	assert.True(t, CompareSliceSet([]int{1, 2, 3, 7, 9, 12}, out), out)
	assert.True(t, <-finish)
}

func TestParallelMapMap(t *testing.T) {
	data := map[string]int{
		"a": 3,
		"b": 1,
		"c": 8,
		"d": 7,
	}
	mappedData := ParallelMapMap(func(k string, v int) (string, string) {
		time.Sleep(time.Millisecond * time.Duration(v))
		return k, strconv.Itoa(v)
	}, 4, data)
	assert.Equal(t, map[string]string{
		"a": "3",
		"b": "1",
		"c": "8",
		"d": "7",
	}, mappedData)
}
