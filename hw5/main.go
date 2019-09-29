package main

import (
	"fmt"
	"math/rand"
	"time"
)

type worker func() error

// test worker, with random working time and chance to error
func createWorker(id int) worker {
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	timeout := random.Intn(100) + 100
	errchance := random.Intn(100)
	return func() error {
		duration := time.Duration(timeout) * time.Millisecond
		fmt.Printf("Started worker #%d\n", id)
		time.Sleep(duration)
		if errchance < 20 {
			fmt.Printf("Stopped worker #%d with error\n", id)
			return fmt.Errorf("Something wrong")
		}
		fmt.Printf("Stopped worker #%d\n", id)
		return nil
	}
}

func workerPool(workers []worker, maxWorkers int, maxErrors int) {

	results := make(chan error)
	pool := make(chan worker)

	// create n workers
	for i := 0; i < maxWorkers; i++ {
		go func() {
			for {
				select {
				case w, ok := <-pool:
					if !ok { // closed
						return
					}
					results <- w()
				}
			}
		}()
	}

	// counting workers and errors
	working := 0
	errorsCounter := 0

stopsent:
	for _, w := range workers {
		select {
		case pool <- w:
			working++
		case err := <-results:
			working--
			if err != nil {
				errorsCounter++
			}
			if errorsCounter == maxErrors {
				close(pool)
				break stopsent
			}
		}
	}
	// wait working
	for working > 0 {
		select {
		case err := <-results:
			working--
			if err != nil {
			}
		}
	}
}

func main() {
	workersCount := 30
	workers := make([]worker, workersCount)
	for i := 0; i < workersCount; i++ {
		workers[i] = createWorker(i + 1)
	}
	workerPool(workers, 5, 3)
}
