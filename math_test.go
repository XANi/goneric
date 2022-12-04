package goneric

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestSum(t *testing.T) {
	assert.Equal(t, 3, Sum(1, 2))
	assert.Equal(t, 6, Sum(1, 2, 3))
	assert.Equal(t, 10.0, Sum(1.0, 2.0, 3.0, 4.0))
}
func TestSumF64(t *testing.T) {
	assert.Equal(t, 3.0, SumF64(1, 2))
	assert.Equal(t, 6.0, SumF64(1, 2, 3))
	assert.Equal(t, 9.2, SumF64(1.0, 2.0, 3.0, 3.2))
	assert.Equal(t, 400.0, SumF64([]uint8{100, 100, 100, 100}...))

}

func TestMax(t *testing.T) {
	assert.Panics(t, func() { Max([]int{}...) })
	assert.Equal(t, 10, Max(10))
	assert.Equal(t, 10, Max(10, 9))
	assert.Equal(t, 10, Max(1, 2, 3, 4, 5, 6, 7, 8, 9, 10))
	assert.Equal(t, 10, Max(1, 2, 3, 4, 10, 6, 7, 8, 9, 5))
	assert.Equal(t, 10, Max(10, 2, 3, 4, 5, 6, 7, 8, 9, 1))
	assert.Equal(t, 10, Max(1, 2, 3, 4, 5, 10, 7, 10, 9, -10))
	assert.Equal(t, -1, Max(-1, -2, -3, -4, -5, -6, -7, -8, -9, -10))
	assert.Equal(t, -1, Max(-1, -2, -3, -4, -10, -6, -7, -8, -9, -5))
	assert.Equal(t, -1, Max(-10, -2, -3, -4, -10, -6, -7, -8, -9, -1))
	assert.Equal(t, -1, Max(-10, -2, -1, -4, -10, -6, -7, -8, -9, -3))
	assert.Equal(t, 0.0, Max(math.Inf(-1), 0))
	assert.Equal(t, 0.0, Max(0.0, math.Inf(-1)))
	assert.Equal(t, 0.0, Max(-0.0, 0.0))
	assert.Equal(t, 0.0, Max(0.0, -0.0))

}

func TestMin(t *testing.T) {
	assert.Panics(t, func() { Min([]int{}...) })
	assert.Equal(t, 10, Min(10))
	assert.Equal(t, 9, Min(10, 9))
	assert.Equal(t, 1, Min(1, 2, 3, 4, 5, 6, 7, 8, 9, 10))
	assert.Equal(t, 1, Min(1, 2, 3, 4, 10, 6, 7, 8, 9, 5))
	assert.Equal(t, 1, Min(10, 2, 3, 4, 5, 6, 7, 8, 9, 1))
	assert.Equal(t, -10, Min(1, 2, 3, 4, 5, 10, 7, 10, 9, -10))
	assert.Equal(t, -10, Min(-1, -2, -3, -4, -5, -6, -7, -8, -9, -10))
	assert.Equal(t, -10, Min(-1, -2, -3, -4, -10, -6, -7, -8, -9, -5))
	assert.Equal(t, -10, Min(-10, -2, -3, -4, -10, -6, -7, -8, -9, -1))
	assert.Equal(t, -10, Min(-10, -2, -1, -4, -10, -6, -7, -8, -9, -3))
	assert.Equal(t, 0.0, Min(math.Inf(1), 0))
	assert.Equal(t, 0.0, Min(0.0, math.Inf(1)))
	assert.Equal(t, -0.0, Min(-0.0, 0.0))
	assert.Equal(t, -0.0, Min(0.0, -0.0))
}
func TestAvg(t *testing.T) {
	assert.Equal(t, 148, Avg(2, 4, 6, 8, 10, 12, 999))
	assert.Equal(t, 1256, Avg(2, 4, 6, 8, 10, 12, 14, 9999))
	assert.Equal(t, 8, Avg(8))
	assert.Equal(t, 9, Avg(8, 10))
	assert.Equal(t, 2.6, Avg(2.0, 3.2))
}
func TestAvgF64(t *testing.T) {
	assert.Equal(t, 149.0, AvgF64(2, 4, 6, 8, 10, 12, 1001))
	assert.Equal(t, 1256.0, AvgF64(2, 4, 6, 8, 10, 12, 14, 9992))
	assert.Equal(t, 8.0, AvgF64(8))
	assert.Equal(t, 9.0, AvgF64(8, 10))
	assert.Equal(t, 2.6, AvgF64(2.0, 3.2))
}
func TestAvgF64F64(t *testing.T) {
	assert.Equal(t, 149.0, AvgF64F64(2, 4, 6, 8, 10, 12, 1001))
	assert.Equal(t, 1256.0, AvgF64F64(2, 4, 6, 8, 10, 12, 14, 9992))
	assert.Equal(t, 8.0, AvgF64F64(8))
	assert.Equal(t, 9.0, AvgF64F64(8, 10))
	assert.Equal(t, 2.6, AvgF64F64(2.0, 3.2))
	assert.Equal(t, 100.0, AvgF64F64([]uint8{100, 100, 100, 100, 100}...))
}
func TestMedian(t *testing.T) {
	assert.Equal(t, 8, Median(2, 4, 6, 8, 10, 12, 9999))
	assert.Equal(t, 9, Median(2, 4, 6, 8, 10, 12, 14, 9999))
	assert.Equal(t, 8, Median(8))
	assert.Equal(t, 9, Median(8, 10))
	assert.Equal(t, 2.5, Median(2.0, 3.0))
}
func TestMedianF64(t *testing.T) {
	assert.Equal(t, 8.0, MedianF64(2, 4, 6, 8, 10, 12, 9999))
	assert.Equal(t, 9.0, MedianF64(2, 4, 6, 8, 10, 12, 14, 9999))
	assert.Equal(t, 8.0, MedianF64(8))
	assert.Equal(t, 8.5, MedianF64(8, 9))
	assert.Equal(t, 2.5, MedianF64(2.0, 3.0))
}
