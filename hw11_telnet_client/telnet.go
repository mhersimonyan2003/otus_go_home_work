package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

func (c *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	c.conn = conn
	fmt.Fprintf(os.Stderr, "...Connected to %s\n", c.address)
	return nil
}

func (c *telnetClient) Send() error {
	_, err := io.Copy(c.conn, c.in)
	var netErr *net.OpError

	if err != nil {
		if errors.As(err, &netErr) {
			fmt.Fprintln(os.Stderr, "Connection was closed by peer")
			return nil
		}

		return err
	}

	_, err = c.in.Read([]byte{0})

	if errors.Is(err, io.EOF) {
		fmt.Fprintln(os.Stderr, "...EOF")
		return nil
	}

	return nil
}

func (c *telnetClient) Receive() error {
	_, err := io.Copy(c.out, c.conn)
	if err != nil {
		return err
	}
	return err
}

func (c *telnetClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
