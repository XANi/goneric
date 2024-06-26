package goneric

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
)

type ComplexSlice1 struct {
	Name  string
	Value string
}

var complexSlice1 = []ComplexSlice1{
	{
		Name:  "t1",
		Value: "v1",
	}, {
		Name:  "t2",
		Value: "v2",
	}, {
		Name:  "t3",
		Value: "v3",
	},
}

func TestCompareSliceSet(t *testing.T) {
	equal1a := []string{"a", "ab", "b"}
	equal1b := []string{"a", "ab", "b"}
	assert.True(t, CompareSliceSet(equal1a, equal1b))
	equal2a := []string{"ab", "b", "a", "ab"}
	equal2b := []string{"a", "ab", "b"}
	assert.True(t, CompareSliceSet(equal2a, equal2b))
	equal3a := []int{2, 1, 3}
	equal3b := []int{3, 2, 1}
	assert.True(t, CompareSliceSet(equal3a, equal3b))
	unequal1a := []string{"a", "ab", "c"}
	unequal1b := []string{"a", "ab", "b"}
	assert.False(t, CompareSliceSet(unequal1a, unequal1b))
	unequal2a := []int{2, 1, 4}
	unequal2b := []int{3, 2, 1}
	assert.False(t, CompareSliceSet(unequal2a, unequal2b))
	unequal3a := []int{2, 1}
	unequal3b := []int{3, 2, 1}
	assert.False(t, CompareSliceSet(unequal3a, unequal3b))
	unequal4a := []int{1, 2, 3}
	unequal4b := []int{1, 3, 3}
	assert.False(t, CompareSliceSet(unequal4a, unequal4b))
}

func TestSliceMap(t *testing.T) {
	sliceMap := SliceMap(func(t ComplexSlice1) string {
		return t.Name
	}, complexSlice1)
	assert.Equal(t, complexSlice1[0], sliceMap["t1"])
	assert.Equal(t, complexSlice1[1], sliceMap["t2"])
	assert.Equal(t, complexSlice1[2], sliceMap["t3"])
	assert.Len(t, sliceMap, 3)
}
func TestSliceMapSkip(t *testing.T) {

	sliceMap := SliceMapSkip(func(t ComplexSlice1) (string, bool) {
		if t.Name == "t2" {
			return t.Name, true
		} else {
			return t.Name, false
		}
	}, complexSlice1)
	assert.Equal(t, complexSlice1[0], sliceMap["t1"])
	assert.Equal(t, complexSlice1[2], sliceMap["t3"])
	assert.Len(t, sliceMap, 2)
}

func ExampleSliceMapSkip() {
	type CS struct {
		Name  string
		Value string
	}
	data := []CS{{Name: "t1", Value: "v1"}, {Name: "t2", Value: "v2"}, {Name: "t3", Value: "v3"}, {Name: "t4", Value: "v4"}}
	sliceMap := SliceMapSkip(func(t CS) (string, bool) {
		if t.Name == "t2" {
			return t.Name, true
		} else {
			return t.Name, false
		}
	}, data)
	fmt.Printf("map from slice with skipped t2: [%+v]", sliceMap)
	//Output: map from slice with skipped t2: [map[t1:{Name:t1 Value:v1} t3:{Name:t3 Value:v3} t4:{Name:t4 Value:v4}]]
}

func TestSliceMapSet(t *testing.T) {
	sliceMap := SliceMapSet([]string{"t1", "t2", "t3"})
	assert.Equal(t, sliceMap["t1"], true)
	assert.Equal(t, sliceMap["t2"], true)
	assert.Equal(t, sliceMap["t3"], true)
	assert.Equal(t, sliceMap["t4"], false)
	assert.Len(t, sliceMap, 3)

}
func TestSliceMapSetFunc(t *testing.T) {
	sliceMap := SliceMapSetFunc(func(c ComplexSlice1) string { return c.Name }, complexSlice1)
	assert.Equal(t, sliceMap["t1"], true)
	assert.Equal(t, sliceMap["t2"], true)
	assert.Equal(t, sliceMap["t3"], true)
	assert.Equal(t, sliceMap["t4"], false)
	assert.Len(t, sliceMap, 3)
}

func TestSliceMapFunc(t *testing.T) {
	sliceMap := SliceMapFunc(
		func(s string) (string, string) {
			v := strings.Split(s, ":")
			return v[0], v[1]
		},
		[]string{"a:b", "c:d"},
	)
	assert.Equal(t, map[string]string{"a": "b", "c": "d"}, sliceMap)
}

func TestSliceDiff(t *testing.T) {
	equala := []string{"a", "b", "c"}
	equalb := []string{"a", "b", "c"}
	out_left, out_right := SliceDiff(equala, equalb)
	assert.Len(t, out_left, 0)
	assert.Len(t, out_right, 0)

	lefta := []int{1, 2, 3, 4}
	leftb := []int{1, 2}
	out_lefta, out_leftb := SliceDiff(lefta, leftb)
	assert.Equal(t, []int{3, 4}, out_lefta, "left side")
	assert.Equal(t, []int{}, out_leftb, "right side")

	righta := []float64{3, 4}
	rightb := []float64{1, 2, 3, 4}
	out_righta, out_rightb := SliceDiff(righta, rightb)
	assert.Equal(t, []float64{}, out_righta, "left side")
	assert.Equal(t, []float64{1, 2}, out_rightb, "right side")

}

func TestSliceDiffFunc(t *testing.T) {
	equala := []string{"1", "2", "3"}
	equalb := []int{1, 2, 3}
	out_left, out_right := SliceDiffFunc(
		equala,
		equalb,
		func(s string) int {
			i, err := strconv.Atoi(s)
			if err != nil {
				panic(s)
			}
			return i
		},
		func(i int) int {
			return i
		},
	)
	assert.Len(t, out_left, 0)
	assert.Len(t, out_right, 0)

	lefta := []float32{1, 2, 3, 4}
	leftb := []int{1, 2}
	out_lefta, out_leftb := SliceDiffFunc(
		lefta,
		leftb,
		func(i float32) int {
			return int(i)
		},
		func(i int) int {
			return i
		},
	)
	assert.Equal(t, []float32{3, 4}, out_lefta, "left side")
	assert.Equal(t, []int{}, out_leftb, "right side")

	type Sl1 struct {
		Name string
	}
	type Sl2 struct {
		Login string
	}

	righta := []Sl1{
		{Name: "a"},
		{Name: "b"},
	}
	rightb := []Sl2{
		{Login: "a"},
		{Login: "b"},
		{Login: "c"},
		{Login: "d"},
	}

	out_righta, out_rightb := SliceDiffFunc(
		righta,
		rightb,
		func(v1 Sl1) string {
			return v1.Name
		},
		func(v2 Sl2) string {
			return v2.Login
		},
	)
	assert.Equal(t, []Sl1{}, out_righta, "left side")
	assert.Equal(t, []Sl2{{Login: "c"}, {Login: "d"}}, out_rightb, "right side")
}

func ExampleSliceDiffFunc() {
	data1 := []string{"1", "2", "3", "4", "5"}
	data2 := []float32{1, 7, 3, 4}
	stringToInt := func(s string) int {
		i, _ := strconv.Atoi(s)
		return i
	}
	floatToInt := func(f float32) int {
		return int(f)
	}
	left, right := SliceDiffFunc(data1, data2, stringToInt, floatToInt)
	fmt.Printf("left: %T%+v right: %T%+v", left, left, right, right)
	//Output: left: []string[2 5] right: []float32[7]
}

func TestSliceIn(t *testing.T) {
	type Comparable struct {
		Name string
	}
	assert.True(t, SliceIn([]int{1, 2, 3, 4}, 3))
	assert.True(t, SliceIn(
		[]Comparable{{Name: "t1"}, {Name: "t2"}},
		Comparable{Name: "t2"}))
	assert.False(t, SliceIn([]int{1, 2, 3, 4}, 5))
	assert.False(t, SliceIn(
		[]Comparable{{Name: "t1"}, {Name: "t2"}},
		Comparable{Name: "t0"}))
}

type incomparable struct {
	v []int
}

func TestSliceDedupe(t *testing.T) {

	assert.Equal(t, []int{1, 2, 3, 4}, SliceDedupe([]int{1, 2, 3, 2, 1, 4, 3, 4}))

}

func TestSliceDedupeFunc(t *testing.T) {
	assert.Equal(t,
		[]incomparable{
			{v: []int{1, 2}},
			{v: []int{2, 4}},
			{v: []int{4, 2, 1}},
			{v: []int{1, 3}},
		},
		SliceDedupeFunc(
			[]incomparable{
				{v: []int{1, 2}},
				{v: []int{2, 4}},
				{v: []int{4, 2, 1}},
				{v: []int{1, 3}},
				{v: []int{4, 2, 1}},
				{v: []int{1, 2}},
			},
			func(in incomparable) string {
				return fmt.Sprintf("%+v", in)
			},
		))
}

func TestSliceReverse(t *testing.T) {
	a := []int{1, 9, 4, 3, 7}
	b := SliceReverse(a)
	assert.Equal(t, []int{7, 3, 4, 9, 1}, SliceReverse(a))
	assert.NotEqual(t, a, b)
}

func TestSliceReverseInplace(t *testing.T) {
	a := []int{1, 9, 4, 3, 7}
	SliceReverseInplace(a)
	assert.Equal(t, []int{7, 3, 4, 9, 1}, a)
}

func TestFirstOrEmpty(t *testing.T) {
	assert.Equal(t, 4, FirstOrEmpty([]int{4, 3, 2}))
	assert.Equal(t, 0, FirstOrEmpty([]int{}))
}
func TestLastOrEmpty(t *testing.T) {
	assert.Equal(t, 2, LastOrEmpty([]int{4, 3, 2}))
	assert.Equal(t, 0, LastOrEmpty([]int{}))
}
func TestFirstOrDefault(t *testing.T) {
	assert.Equal(t, 4, FirstOrDefault([]int{4, 3, 2}, 7))
	assert.Equal(t, 7, FirstOrDefault([]int{}, 7))
}
func TestLastOrDefault(t *testing.T) {
	assert.Equal(t, 2, LastOrDefault([]int{4, 3, 2}, 7))
	assert.Equal(t, 7, LastOrDefault([]int{}, 7))
}
