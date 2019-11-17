package contexttools

import (
	"context"
	"io"
	"net"
	"sync/atomic"
)

// ContextReader - main object, processing read from context
type ContextReader struct {
	data   chan []byte
	err    error
	flag   int32
	closed int32
}

// CreateReader - create helper around read from io.Reader
func CreateReader(ctx context.Context, reader io.Reader) *ContextReader {
	buffer := make([]byte, 65536)

	s := &ContextReader{}
	s.data = make(chan []byte, 1)
	go func() {
		<-ctx.Done()
		s.close()
	}()

	go func() {
	loop:
		for {
			readed, err := reader.Read(buffer)
			if err != nil {
				s.err = err
				break loop
			}
			if readed == 0 {
				select {
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
			case <-ctx.Done():
				break loop
			}
		}
		s.close()
	}()
	return s
}

func (s *ContextReader) close() {
	if !atomic.CompareAndSwapInt32(&s.closed, 0, 1) {
		return
	}
	close(s.data)
}

// Read - return channel for read from scanner
func (s *ContextReader) Read() <-chan []byte {
	return s.data
}

// GetLastError - return error if it happend
func (s *ContextReader) GetLastError() error {
	if s.err == io.EOF {
		return nil
	}
	if _, ok := s.err.(*net.OpError); ok {
		return nil
	}
	return s.err
}
