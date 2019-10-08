package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// Parse parameters
func main() {

	var from, to string
	var offset, limit int64

	flag.StringVar(&from, "from", "", "Source file path")
	flag.StringVar(&to, "to", "", "Source file path")
	flag.Int64Var(&offset, "offset", 0, "Start position offset [optional]")
	flag.Int64Var(&limit, "limit", 0, "Limit copying data length [optional]")
	flag.Parse()

	if from == "" {
		fmt.Println("Error: from parameter not declared")
	}
	if to == "" {
		fmt.Println("Error: to parameter not declared")
	}
	if offset < 0 {
		fmt.Println("Error: offset parameter in negative")
	}
	if limit < 0 {
		fmt.Println("Error: limit parameter in negative")
	}
	if from == "" || to == "" || offset < 0 || limit < 0 {
		fmt.Println("Program copies part of file to another file.")
		flag.Usage()
		return
	}
	err := copy(from, to, limit, offset)
	if err != nil {
		fmt.Println("Error:", err.Error())
	}
}

// Main job function
func copy(from, to string, offset, limit int64) error {
	src, err := os.Open(from)
	if err != nil {
		return err
	}
	defer src.Close()
	stat, err := src.Stat()
	if err != nil {
		return err
	}
	if offset > 0 {
		if offset >= stat.Size() {
			return fmt.Errorf("Offset exceeds size of source file")
		}
		_, err = src.Seek(offset, io.SeekStart)
		if err != nil {
			return err
		}
	}
	dst, err := os.Create(to)
	if err != nil {
		return err
	}
	defer dst.Close()
	copied, err := copyio(src, dst, limit)
	if err != nil {
		return err
	}
	fmt.Printf("Copied %d bytes, \n", copied)
	return nil
}

// copyio - function to run tests
func copyio(src io.Reader, dst io.Writer, limit int64) (int64, error) {
	var copied int64
	buffer := make([]byte, 0, 16384)
	for {
		readed, err := src.Read(buffer)
		if readed > 0 {
			written, err := dst.Write(buffer)
			if err != nil {
				return copied, err
			}
			copied += int64(written)
		}
		if err == io.EOF {
			return copied, nil
		}
		if err != nil {
			return copied, err
		}
	}
}
