package main

import (
	"fmt"
	"math/rand"
	"time"
)

type worker func() error

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
	counter := make(chan struct{}, maxWorkers)
	results := make(chan error, 1)
	die := make(chan struct{})

	for _, workerfunc := range workers {
		go func(w worker) {
			select {
			case <-die:
				results <- nil
				return
			case counter <- struct{}{}:
				break
			}
			results <- w()
		}(workerfunc)
	}

	errorsCounter := 0
	finished := 0
	diesent := false
	for {
		select {
		case err := <-results:
			finished++
			if err != nil {
				errorsCounter++
			}
			if errorsCounter == maxErrors && !diesent {
				diesent = true
				close(die)
			}
			if !diesent {
				<-counter
			}
		}
		if finished == len(workers) {
			return
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
