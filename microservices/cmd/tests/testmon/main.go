package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gmax79/otusgolang/microservices/internal/grpccon"
)

const host = "localhost:9999"

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	fmt.Println("Testing calendar monitoring. App create some load")
	defer fmt.Println("Load app finished")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	var err error
	cli, err := grpccon.CreateClient(ctx, host)
	defer cli.Close()
	cancel()
	assert(err)

	ctx, cancel = context.WithCancel(context.Background())
	go func() {

	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

}
