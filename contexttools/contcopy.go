package contexttools

import (
	"context"
	"io"
	"sync/atomic"
	"time"
)

// Copier - object to copy from reader to writer with context control
type Copier struct {
	ctx     context.Context
	timeout chan struct{}
	flag    int32
	tout    bool
}

// CreateCopier - create it
func CreateCopier(ctx context.Context) *Copier {
	c := &Copier{}
	c.ctx = ctx
	c.timeout = make(chan struct{})
	return c
}

// Copy - function to copy from in to out with context
func (c *Copier) Copy(in io.Reader, out io.Writer) error {
	reader := CreateReader(c.ctx, in)
	for {
		select {
		case data, ok := <-reader.Read():
			if !ok {
				return reader.GetLastError()
			}
			if _, err := out.Write(data); err != nil {
				return err
			}
			atomic.SwapInt32(&c.flag, 1)
		case <-c.timeout:
			c.tout = true
			return nil
		}
	}
}

// AddTimeout - add timeout for read operation
func (c *Copier) AddTimeout(timeout time.Duration) {
	ticker := time.NewTicker(timeout)
	go func() {
	loop:
		for {
			select {
			case <-ticker.C:
				if atomic.SwapInt32(&c.flag, 0) == 0 {
					c.timeout <- struct{}{}
					break loop
				}
			case <-c.ctx.Done():
				break loop
			}
		}
		ticker.Stop()
	}()
}

// IsTimeout - check if copier stopped via timeout
func (c *Copier) IsTimeout() bool {
	return c.tout
}
