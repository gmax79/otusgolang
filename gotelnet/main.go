package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gmax79/otusgolang/contexttools"
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

	ctx, cancel := context.WithCancel(context.Background())
	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGINT, syscall.SIGTERM)

	var syncflag int32
	go func() {
		<-term
		atomic.SwapInt32(&syncflag, 1)
		cancel()
	}()

	writer := contexttools.CreateCopier(ctx)
	reader := contexttools.CreateCopier(ctx)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		err := writer.Copy(os.Stdin, conn)
		stdErr(err)
		if atomic.SwapInt32(&syncflag, 1) == 0 {
			fmt.Fprintln(os.Stdout, "<<EOF>>, receive last data")
			reader.AddTimeout(time.Second * 3)
		}
		wg.Done()
	}()
	go func() {
		err := reader.Copy(conn, os.Stdout)
		stdErr(err)
		wg.Done()
	}()
	wg.Wait()
}
