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
	ch := ParallelMapSliceChannel(func(s string) int {
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
	ch, finish := ParallelMapSliceChannelFinisher(func(s string) int {
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
