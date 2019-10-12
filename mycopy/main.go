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
	flag.StringVar(&to, "to", "", "Destination file path")
	flag.Int64Var(&offset, "offset", 0, "Start position offset [optional, > 0]")
	flag.Int64Var(&limit, "limit", 0, "Limit copying data length [optional, > 0]")
	flag.Parse()
	if from == "" || to == "" || offset < 0 || limit < 0 {
		fmt.Println("Error: Set invalid parameters")
		flag.Usage()
		return
	}
	err := copy(from, to, offset, limit)
	if err != nil {
		fmt.Println("Error:", err.Error())
	} else {
		fmt.Println("")
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
	srcsize := stat.Size()
	if offset > 0 {
		if offset >= srcsize {
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
	_, err = copyio(src, dst, limit, srcsize)
	return err
}

// copyio - function to run tests
func copyio(src io.Reader, dst io.Writer, limit, size int64) (copied int64, err error) {
	buffer := make([]byte, 16384)
	for {
		var readed int
		readed, err = src.Read(buffer)
		if readed > 0 {
			readed64 := int64(readed)
			if limit > 0 && readed64 > limit {
				readed64 = limit
			}
			_, err = dst.Write(buffer[:readed64])
			if err != nil {
				return
			}
			copied += readed64
			var percents float64
			if size > 0 {
				percents = float64(copied) / float64(size)
				percents *= 100
			}
			fmt.Printf("\033[2K\rProgress: %.2f%%", percents) // Inplace print
			//time.Sleep(5 * time.Millisecond)                // Demonstration delay
			if limit > 0 {
				limit -= readed64
				if limit == 0 {
					return copied, nil
				}
			}
		}
		if err == io.EOF {
			return copied, nil
		}
		if err != nil {
			return
		}
	}
}
