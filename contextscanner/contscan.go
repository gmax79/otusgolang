package contextscanner

import (
	"bufio"
	"context"
	"io"
	"net"
	"sync/atomic"
	"time"
)

// ContextScanner - main object, processing read from context
type ContextScanner struct {
	data chan []byte
	err  error
	flag int32
}

// ScanChunks is a split function for a bufio.Scanner that returns byte chunks as is
func scanChunks(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	return len(data), data, nil
}

// Create - create helper around bufio.Scanner
func Create(ctx context.Context, reader io.ReadCloser, timeout time.Duration) *ContextScanner {
	s := &ContextScanner{}
	s.data = make(chan []byte)

	if timeout > 0 {
		ticker := time.NewTicker(timeout)
		go func() {
			for {
				select {
				case <-ticker.C:
					if atomic.SwapInt32(&s.flag, 0) == 0 {
						ticker.Stop()
						reader.Close()
						return
					}
				case <-ctx.Done():
					ticker.Stop()
					return
				}
			}
		}()
	}

	go func() {
		<-ctx.Done()
		reader.Close()
	}()

	go func() {
		scanner := bufio.NewScanner(reader)
		scanner.Split(scanChunks)
	loop:
		for {
			if !scanner.Scan() {
				s.err = scanner.Err()
				break loop
			}
			data := scanner.Bytes()
			if len(data) == 0 {
				select {
				case <-ctx.Done():
					break loop
				default:
				}
				continue
			}
			select {
			case s.data <- data:
				atomic.SwapInt32(&s.flag, 1)
			case <-ctx.Done():
				break loop
			}
		}
		close(s.data)
	}()
	return s
}

// Read - return channel for read from scanner
func (s *ContextScanner) Read() <-chan []byte {
	return s.data
}

// GetLastError - return error if it happend
func (s *ContextScanner) GetLastError() error {
	if _, ok := s.err.(*net.OpError); ok {
		return nil
	}
	return s.err
}
