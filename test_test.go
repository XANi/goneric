package goneric

import "time"

// test utils

type ctr struct {
	i int
}

func (c *ctr) Counter() int {
	c.i++
	return c.i
}

func (c *ctr) SleepyCounter(t time.Duration) int {
	c.i++
	time.Sleep(t)
	return c.i
}
