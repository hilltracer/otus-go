package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{
		addr:    address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type client struct {
	addr    string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

func (c *client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.addr, c.timeout)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *client) Send() error {
	_, err := io.Copy(c.conn, c.in) // read from stdin, write to socket
	return err                      // nil == graceful EOF (Ctrl-D)
}

func (c *client) Receive() error {
	_, err := io.Copy(c.out, c.conn) // read from socket, write to stdout
	return err                       // nil == peer closed connection
}
