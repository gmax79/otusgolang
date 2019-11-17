package contextscanner

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync/atomic"
	"time"
)

// ContextScanner - main object, processing read from context
type ContextScanner struct {
	data   chan []byte
	err    error
	flag   int32
	closed int32
	//buffer []byte
}

// ScanChunks is a split function for a bufio.Scanner that returns byte chunks as is
/*func scanChunks(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	return len(data), data, nil
}

type readerWithClose struct {
	internal io.Reader
	closed   bool
}

func (r *readerWithClose) Read(p []byte) (n int, err error) {
	if r.closed {
		return 0, io.EOF
	}
	return r.internal.Read(p)
}

func (r *readerWithClose) Close() {
	r.closed = true
}*/

// Create - create helper around read from io.Reader
func Create(ctx context.Context, reader io.Reader, timeout time.Duration) *ContextScanner {
	timeoutchan := make(chan struct{}, 1)
	buffer := make([]byte, 65536)

	s := &ContextScanner{}
	s.data = make(chan []byte, 1)
	go func() {
		<-ctx.Done()
		s.close()
	}()

	if timeout > 0 {
		ticker := time.NewTicker(timeout)
		go func() {
		loop:
			for {
				select {
				case <-ticker.C:
					if atomic.SwapInt32(&s.flag, 0) == 0 {
						timeoutchan <- struct{}{}
						fmt.Println("timeout")
						break loop
					}
				case <-ctx.Done():
					break loop
				}
			}
			ticker.Stop()
			s.close()
		}()
	}

	go func() {
	loop:
		for {
			fmt.Println("read")
			readed, err := reader.Read(buffer)
			if err != nil {
				s.err = err
				break loop
			}
			if readed == 0 {
				select {
				case <-timeoutchan:
					break loop
				case <-ctx.Done():
					break loop
				default:
				}
				continue
			}
			data := make([]byte, readed)
			copy(data, buffer[0:])
			select {
			case s.data <- data:
				atomic.SwapInt32(&s.flag, 1)
			case <-timeoutchan:
				break loop
			case <-ctx.Done():
				break loop
			}
		}
		s.close()
	}()
	return s
}

func (s *ContextScanner) close() {
	if !atomic.CompareAndSwapInt32(&s.closed, 0, 1) {
		return
	}
	close(s.data)
}

// Read - return channel for read from scanner
func (s *ContextScanner) Read() <-chan []byte {
	return s.data
}

// GetLastError - return error if it happend
func (s *ContextScanner) GetLastError() error {
	if s.err == io.EOF {
		return nil
	}
	if _, ok := s.err.(*net.OpError); ok {
		return nil
	}
	return s.err
}
