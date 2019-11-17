package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gmax79/otusgolang/contextscanner"
)

func stdErrAndExit(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
	os.Exit(1)
}

func stdErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func main() {
	var timeout string
	flag.StringVar(&timeout, "timeout", "10s", "Connection timeout")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "[-timeout value] host port")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "Application read from stdin and send it to remote host.")
		fmt.Fprintln(os.Stderr, "Answers from remote host prints into stdout.")
	}
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		stdErrAndExit("Invalid counts of parameters")
	}
	timeoutDuration, err := time.ParseDuration(timeout)
	if err != nil {
		stdErrAndExit(err.Error())
	}

	remoteHost := args[0] + ":" + args[1]
	fmt.Fprintln(os.Stdout, "Connecting:", remoteHost)

	dialer := &net.Dialer{}
	connctx, conncancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer conncancel()
	conn, err := dialer.DialContext(connctx, "tcp", remoteHost)
	if err != nil {
		stdErrAndExit("Cannot connect:", err.Error())
	}
	defer func() {
		conn.Close()
		fmt.Fprintln(os.Stdout, "Connection closed")
	}()

	writectx, writecancel := context.WithCancel(context.Background())
	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGINT, syscall.SIGTERM)
	readctx, readcancel := context.WithCancel(context.Background())
	go func() {
		<-term
		writecancel()
		readcancel()
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		err := copyWithContext(writectx, os.Stdin, conn, 0)
		fmt.Fprintln(os.Stdout, "Connection write closed")
		stdErr(err)
		timer := time.NewTimer(time.Second * 3)
		<-timer.C
		readcancel()
		wg.Done()
	}()
	go func() {
		err := copyWithContext(readctx, conn, os.Stdout, 0)
		fmt.Fprintln(os.Stdout, "Connection read closed")
		stdErr(err)
		wg.Done()
	}()
	wg.Wait()
}

func copyWithContext(ctx context.Context, in io.ReadCloser, out io.Writer, readtimeout time.Duration) error {
	scanner := contextscanner.Create(ctx, in, readtimeout)
	for {
		data, ok := <-scanner.Read()
		if !ok {
			return scanner.GetLastError()
		}
		if _, err := out.Write(data); err != nil {
			return err
		}
	}
}
