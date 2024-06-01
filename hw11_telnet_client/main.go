package main

import (
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

func main() {
	// Example usage
	var wg sync.WaitGroup
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	// Extract host and port from command-line arguments
	args := flag.Args()
	if len(args) != 2 {
		fmt.Println("Usage: go-telnet [--timeout=<duration>] <host> <port>")
		os.Exit(1)
	}
	host := args[0]
	port := args[1]
	address := net.JoinHostPort(host, port)

	in := io.ReadCloser(os.Stdin) // You need to define your own input stream here
	out := os.Stdout              // You need to define your own output stream here

	client := NewTelnetClient(address, timeout, in, out)
	defer client.Close()

	err := client.Connect()
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}

	wg.Add(2)

	go func() {
		defer wg.Done()
		err := client.Receive()
		if err != nil {
			fmt.Println("Error receiving:", err)
		}
	}()

	go func() {
		defer wg.Done()
		err = client.Send()
		if err != nil {
			fmt.Println("Error sending:", err)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM)

	go func() {
		<-interrupt
		fmt.Println("...EOF")
		os.Exit(0)
	}()

	wg.Wait()
}
