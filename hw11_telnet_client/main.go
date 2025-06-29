package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout (e.g. 5s, 1m)")
}

func main() {
	// CLI parsing
	flag.Parse()

	if flag.NArg() != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [--timeout=10s] host port\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(2)
	}

	host := flag.Arg(0)
	port := flag.Arg(1)
	address := net.JoinHostPort(host, port)

	// Create / connect client
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		fmt.Fprintln(os.Stderr, "connect error:", err)
		os.Exit(1)
	}
	defer client.Close()
	fmt.Fprintf(os.Stderr, "...Connected to %s\n", address)

	// Context cancelled by SIGINT (Ctrl-C)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Run Send / Receive concurrently
	var wg sync.WaitGroup
	var once sync.Once
	done := make(chan struct{})

	closeDone := func(msg string) {
		once.Do(func() {
			if msg != "" {
				fmt.Fprintln(os.Stderr, msg)
			}
			close(done)
		})
	}

	wg.Add(1)
	go func() { // socket -> stdout
		defer wg.Done()
		if err := client.Receive(); err != nil && !errors.Is(err, io.EOF) {
			closeDone("receive error: " + err.Error())
		} else {
			closeDone("...Connection was closed by peer")
		}
	}()

	wg.Add(1)
	go func() { // stdin -> socket
		defer wg.Done()
		if err := client.Send(); err != nil && !errors.Is(err, io.EOF) {
			closeDone("send error: " + err.Error())
		} else {
			closeDone("...EOF")
		}
	}()

	// Wait for either Ctrl-C or I/O completion
	select {
	case <-done:
	case <-ctx.Done():
		fmt.Fprintln(os.Stderr, "...SIGINT")
	}

	_ = client.Close()
	wg.Wait()
}
