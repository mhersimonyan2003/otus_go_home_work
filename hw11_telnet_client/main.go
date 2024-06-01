package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	// Example usage
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

	go func() {
		err := client.Receive()
		if err != nil {
			fmt.Println("Error receiving:", err)
		}
	}()

	err = client.Send()
	if err != nil {
		fmt.Println("Error sending:", err)
	}
}
