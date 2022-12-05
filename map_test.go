package goneric

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sort"
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {

	mappedData := Map(func(v string) int {
		i, _ := strconv.Atoi(v)
		return i
	}, "1", "3", "2")
	assert.Equal(t, []int{1, 3, 2}, mappedData)
}

func TestMapErr(t *testing.T) {
	mappedData, err := MapErr(func(v string) (int, error) {
		i, err := strconv.Atoi(v)
		return i, err
	}, "1", "2", "z", "4")
	assert.Error(t, err)
	assert.Equal(t, []int{1, 2}, mappedData)
	mappedData2, err := MapErr(func(v string) (int, error) {
		i, err := strconv.Atoi(v)
		return i, err
	}, "1", "2", "3", "4")
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4}, mappedData2)
}

func TestMapSlice(t *testing.T) {
	mappedData := MapSlice(func(v string) int {
		i, _ := strconv.Atoi(v)
		return i
	}, []string{"1", "3", "2"})
	assert.Equal(t, []int{1, 3, 2}, mappedData)
}

func TestMapSliceErr(t *testing.T) {
	mappedData, err := MapSliceErr(func(v string) (int, error) {
		i, err := strconv.Atoi(v)
		return i, err
	}, []string{"1", "2", "z", "4"})
	assert.Error(t, err)
	assert.Equal(t, []int{1, 2}, mappedData)
	mappedData2, err := MapSliceErr(func(v string) (int, error) {
		i, err := strconv.Atoi(v)
		return i, err
	}, []string{"1", "2", "3", "4"})
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4}, mappedData2)
}

func TestMapSliceSkip(t *testing.T) {
	mappedData := MapSliceSkip(func(v string) (int, bool) {
		i, err := strconv.Atoi(v)
		return i, err != nil
	}, []string{"1", "2", "z", "4"})
	assert.Equal(t, []int{1, 2, 4}, mappedData)
	mappedData2 := MapSliceSkip(func(v string) (int, bool) {
		i, err := strconv.Atoi(v)
		return i, err != nil
	}, []string{"1", "2", "3", "4"})
	assert.Equal(t, []int{1, 2, 3, 4}, mappedData2)
}

func TestMapSliceErrSkip(t *testing.T) {
	mappedData, err := MapSliceErrSkip(func(v string) (int, error) {
		i, err := strconv.Atoi(v)
		if v == "j" {
			return i, err
		}
		if err != nil {
			return i, ErrSkip{}
		}
		return i, nil
	}, []string{"1", "2", "z", "4", "j", "6"})
	assert.Error(t, err)
	assert.Equal(t, []int{1, 2, 4}, mappedData)
	mappedData2, err := MapSliceErrSkip(func(v string) (int, error) {
		i, err := strconv.Atoi(v)
		return i, err
	}, []string{"1", "2", "3", "4"})
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4}, mappedData2)
}

func TestMapSliceKey(t *testing.T) {
	data := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	dataSlice := MapSliceKey(data)
	sort.Strings(dataSlice)
	assert.EqualValues(t, []string{"a", "b", "c"}, dataSlice)
}

func TestMapSliceValue(t *testing.T) {
	data := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	dataSlice := MapSliceValue(data)
	sort.Ints(dataSlice)
	assert.EqualValues(t, []int{1, 2, 3}, dataSlice)
}

func TestMapToSlice(t *testing.T) {
	assert.True(t,
		CompareSliceSet(
			[]string{"a-1", "b-2"},
			MapToSlice(
				func(k string, v int) string {
					return fmt.Sprintf("%s-%d", k, v)
				}, map[string]int{
					"a": 1,
					"b": 2,
				},
			)))
}
