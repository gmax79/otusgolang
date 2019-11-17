package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gmax79/otusgolang/contexttools"
)

func main() {
	log.Println("Echo tcp server at port 8080")
	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatalf("Cannot listen: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-term
		cancel()
		listener.Close()
	}()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			conn, err := listener.Accept()
			if err != nil {
				if _, ok := err.(*net.OpError); !ok {
					log.Printf("Cannot accept: %T", err)
				}
				break
			}
			wg.Add(1)
			go handleConnection(ctx, conn, wg)
		}
	}()
	wg.Wait()
	log.Println("Echo server stopped")
}

func handleConnection(ctx context.Context, conn net.Conn, wg *sync.WaitGroup) {
	reader := contexttools.CreateReader(ctx, conn)
	remoteAddr := conn.RemoteAddr()
	defer func() {
		log.Printf("Closing connection with %s\n", remoteAddr)
		conn.Close()
		wg.Done()
	}()
	log.Printf("Connected from %s\n", remoteAddr)
	for {
		select {
		case data, ok := <-reader.Read():
			if !ok {
				err := reader.GetLastError()
				if err != nil {
					log.Println(err.Error())
				}
				return
			}
			log.Printf("Received %d bytes from %s: %s", len(data), remoteAddr, string(data))
			// send echo answer
			conn.Write([]byte("echo:"))
			conn.Write(data)
		}
	}
}
