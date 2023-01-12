package goneric

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestAsync(t *testing.T) {
	out := Async(func() int { return 1 })
	assert.Equal(t, 1, <-out)
}
func TestAsyncV(t *testing.T) {
	out := AsyncV(
		func() int { return 2 },
		func() int { return 1 },
	)
	assert.True(t, CompareSliceSet([]int{1, 2}, []int{<-out, <-out}))
}

func TestAsyncVUnpanic(t *testing.T) {
	out := AsyncVUnpanic(
		func() int { time.Sleep(time.Millisecond * 10); return 2 },
		func() int { panic("fail") },
		func() int { return 1 },
	)
	assert.True(t, CompareSliceSet([]int{1, 2}, []int{<-out, <-out}))
}

func TestAsyncPipe(t *testing.T) {
	out :=
		AsyncPipe(
			AsyncPipe(
				Async(func() int { return 1 }),
				func(i int) string { return strconv.Itoa(i) },
			),
			func(s string) string {
				return fmt.Sprintf("test: %s", s)
			})
	assert.Equal(t, "test: 1", <-out)
}

func TestAsyncOut(t *testing.T) {
	out := make(chan string, 1)
	AsyncOut(
		AsyncPipe(
			Async(func() int { return 1 }),
			func(i int) string { return strconv.Itoa(i) },
		),
		func(s string) string {
			return fmt.Sprintf("test: %s", s)
		},
		out)
	assert.Equal(t, "test: 1", <-out)
}

func ExampleAsync() {
	// make our conversion output channel
	out := make(chan string, 1)

	// queue first element
	AsyncOut(
		AsyncPipe(
			AsyncIn(time.Unix(123456789, 0)),
			func(in time.Time) (out string) { return in.Format("2006-01-02") },
		),
		func(in string) (out string) {
			// just to make sure that one comes second
			time.Sleep(time.Millisecond * 10)
			return "date: " + in
		},
		out)
	// queue second element
	AsyncOut(
		AsyncPipe(
			AsyncPipe(
				AsyncIn(12345678),
				func(in int) (out time.Time) { return time.Unix(int64(in), 0) },
			),
			func(t time.Time) string { return t.Format("2006-01-02") },
		),
		func(in string) (out string) { return "date: " + in },
		out)
	fmt.Printf("%s %s",
		<-out,
		<-out,
	)
	// Output: date: 1970-05-23 date: 1973-11-29
}
