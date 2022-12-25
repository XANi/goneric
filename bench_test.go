package goneric

import (
	"strconv"
	"testing"
)

func BenchmarkMap(b *testing.B) {
	in := GenSlice(10000, func(i int) int { return i })
	b.Run("Map", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Map(func(i int) string {
				return strconv.Itoa(i)
			}, in...)
		}
	})
	b.Run("manual_map", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			out := make([]string, len(in))
			for idx, i := range in {
				out[idx] = func(i int) string { return strconv.Itoa(i) }(i)
			}

		}
	})
	b.Run("manual_map_inline", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			out := make([]string, len(in))
			for idx, i := range in {
				out[idx] = func(i int) string { return strconv.Itoa(i) }(i)
			}

		}
	})
}
